package muse

import (
	"context"
	"crypto/ed25519"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"sync"
	"time"

	"go.sia.tech/siad/modules"
	"lukechampine.com/frand"
	"lukechampine.com/shard"
	"lukechampine.com/us/hostdb"
	"lukechampine.com/us/renter"
	"lukechampine.com/us/renter/proto"
)

func writeJSON(w http.ResponseWriter, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	// encode nil slices as [] instead of null
	if val := reflect.ValueOf(v); val.Kind() == reflect.Slice && val.Len() == 0 {
		w.Write([]byte("[]\n"))
		return
	}
	enc := json.NewEncoder(w)
	enc.SetIndent("", "\t")
	enc.Encode(v)
}

type server struct {
	hostSets map[string][]hostdb.HostPublicKey
	dir      string

	wallet proto.Wallet
	tpool  proto.TransactionPool
	shard  *shard.Client
	mu     sync.Mutex
	utxoMu sync.Mutex // separate mutex for utxos, preventing reuse
}

func (s *server) handleForm(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var rf RequestForm
	if err := json.NewDecoder(req.Body).Decode(&rf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	start := time.Now()
	log.Println("resolving a host key:", rf.HostKey)
	hostAddr, err := s.shard.ResolveHostKey(rf.HostKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rf.Settings.NetAddress = hostAddr
	host := hostdb.ScannedHost{
		HostSettings: rf.Settings,
		PublicKey:    rf.HostKey,
	}
	key := ed25519.NewKeyFromSeed(frand.Bytes(32))
	log.Println("trying to lock utxoMu:", rf.HostKey, time.Since(start))
	s.utxoMu.Lock()
	log.Println("forming a contract:", rf.HostKey, time.Since(start))
	rev, txnSet, err := proto.FormContract(s.wallet, s.tpool, key, host, rf.Funds, rf.StartHeight, rf.EndHeight)
	if err != nil {
		log.Println("release utxoMu:", rf.HostKey, time.Since(start))
		s.utxoMu.Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// submit txnSet to tpool
	//
	// NOTE: if the tpool rejects the transaction, we log the error, but
	// still treat the contract as valid. The host has signed the
	// transaction, so presumably they already submitted it to their own
	// tpool without error, and intend to honor the contract. Our tpool
	// *shouldn't* reject the transaction, but it might if we desync from
	// the network somehow.
	log.Println("submitting transaction set:", rf.HostKey, time.Since(start))
	submitErr := s.tpool.AcceptTransactionSet(txnSet)
	log.Println("release utxoMu:", rf.HostKey, time.Since(start))
	s.utxoMu.Unlock()
	if submitErr != nil && submitErr != modules.ErrDuplicateTransactionSet {
		log.Println("WARN: contract transaction was not accepted", submitErr)
	}

	c := Contract{
		Contract: renter.Contract{
			HostKey:   rev.HostKey(),
			ID:        rev.ID(),
			RenterKey: key,
		},
		HostAddress: host.NetAddress,
		EndHeight:   rf.EndHeight,
	}
	writeJSON(w, c)
	log.Println("forming a contract finishes:", rf.HostKey, time.Since(start))
}

func (s *server) handleRenew(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var rf RequestRenew
	if err := json.NewDecoder(req.Body).Decode(&rf); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	start := time.Now()
	log.Println("resolving a host key:", rf.HostKey)
	hostAddr, err := s.shard.ResolveHostKey(rf.HostKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	rf.Settings.NetAddress = hostAddr
	host := hostdb.ScannedHost{
		PublicKey:    rf.HostKey,
		HostSettings: rf.Settings,
	}
	log.Println("trying to lock utxoMu:", rf.HostKey, time.Since(start))
	s.utxoMu.Lock()
	log.Println("renewing a contract:", rf.HostKey, time.Since(start))
	rev, txnSet, err := proto.RenewContract(s.wallet, s.tpool, rf.ID, rf.RenterKey, host, rf.Funds, rf.StartHeight, rf.EndHeight)
	if err != nil {
		log.Println("release utxoMu:", rf.HostKey, time.Since(start))
		s.utxoMu.Unlock()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// submit txnSet to tpool (see handleForm)
	log.Println("submitting transaction set:", rf.HostKey, time.Since(start))
	submitErr := s.tpool.AcceptTransactionSet(txnSet)
	s.utxoMu.Unlock()
	if submitErr != nil && submitErr != modules.ErrDuplicateTransactionSet {
		log.Println("WARN: contract transaction was not accepted", submitErr)
	}

	c := Contract{
		Contract: renter.Contract{
			HostKey:   rev.HostKey(),
			ID:        rev.ID(),
			RenterKey: rf.RenterKey,
		},
		HostAddress: rf.Settings.NetAddress,
		EndHeight:   rf.EndHeight,
	}
	writeJSON(w, c)
	log.Println("renewing a contract finishes:", rf.HostKey, time.Since(start))
}

func (s *server) handleHostSets(w http.ResponseWriter, req *http.Request) {
	setName := strings.TrimPrefix(req.URL.Path, "/hostsets/")
	if strings.Contains(setName, "/") {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	switch req.Method {
	case http.MethodGet:
		if setName == "" {
			s.mu.Lock()
			setNames := make([]string, 0, len(s.hostSets))
			for name := range s.hostSets {
				setNames = append(setNames, name)
			}
			s.mu.Unlock()
			sort.Strings(setNames)
			writeJSON(w, setNames)
		} else {
			if setName == "" {
				http.Error(w, "No host set name provided", http.StatusBadRequest)
				return
			}
			s.mu.Lock()
			set, ok := s.hostSets[setName]
			hostKeys := append([]hostdb.HostPublicKey(nil), set...)
			s.mu.Unlock()
			if !ok {
				http.Error(w, "No record of that host set", http.StatusBadRequest)
				return
			}
			writeJSON(w, hostKeys)
		}

	case http.MethodPut:
		if setName == "" {
			http.Error(w, "No host set name provided", http.StatusBadRequest)
			return
		}

		var hostKeys []hostdb.HostPublicKey
		if err := json.NewDecoder(req.Body).Decode(&hostKeys); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		sort.Slice(hostKeys, func(i, j int) bool {
			return hostKeys[i] < hostKeys[j]
		})
		s.mu.Lock()
		s.hostSets[setName] = hostKeys
		if len(hostKeys) == 0 {
			delete(s.hostSets, setName)
		}
		hostSetsJSON, _ := json.MarshalIndent(s.hostSets, "", "  ")
		s.mu.Unlock()
		err := ioutil.WriteFile(filepath.Join(s.dir, "hostSets.json"), hostSetsJSON, 0660)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

	default:
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
	}
}

func (s *server) handleScan(w http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}
	var rs RequestScan
	if err := json.NewDecoder(req.Body).Decode(&rs); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	hostAddr, err := s.shard.ResolveHostKey(rs.HostKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	host, err := hostdb.Scan(ctx, hostAddr, rs.HostKey)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, host.HostSettings)
}

// NewServer returns an HTTP handler that serves the muse API.
func NewServer(dir string, wallet proto.Wallet, tpool proto.TransactionPool, shardAddr string) (http.Handler, error) {
	srv := &server{
		wallet: wallet,
		tpool:  tpool,
		shard:  shard.NewClient(shardAddr),
		dir:    dir,
	}

	// load host sets
	hostSetsJSON, err := ioutil.ReadFile(filepath.Join(dir, "hostSets.json"))
	if os.IsNotExist(err) {
		srv.hostSets = make(map[string][]hostdb.HostPublicKey)
	} else {
		if err != nil {
			return nil, err
		}
		if err := json.Unmarshal(hostSetsJSON, &srv.hostSets); err != nil {
			return nil, err
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/form", srv.handleForm)
	mux.HandleFunc("/renew", srv.handleRenew)
	mux.HandleFunc("/hostsets/", srv.handleHostSets)
	mux.HandleFunc("/scan", srv.handleScan)

	// shard proxy
	shardURL, err := url.Parse(shardAddr)
	if err != nil {
		return nil, err
	}
	mux.Handle("/shard/", &httputil.ReverseProxy{Director: func(req *http.Request) {
		req.URL.Scheme = shardURL.Scheme
		req.URL.Host = shardURL.Host
		req.URL.Path = strings.TrimPrefix(req.URL.Path, "/shard")
	}})
	return mux, nil
}
