package muse

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strings"

	"gitlab.com/NebulousLabs/Sia/types"
	"go.uber.org/multierr"
	"lukechampine.com/shard"
	"lukechampine.com/us/hostdb"
	"lukechampine.com/us/renter"
)

// Error is an error wrapper that provides Is function.
type Error struct {
	error
}

// NewError returns an error that formats as the given text.
func NewError(str string) Error {
	return Error{error: errors.New(str)}
}

// Is reports whether this error matches target.
func (e Error) Is(err error) bool {
	return strings.Contains(e.Error(), err.Error())
}

// A Client communicates with a muse server.
type Client struct {
	addr string
	ctx  context.Context
}

func (c *Client) req(method string, route string, data, resp interface{}) (err error) {
	var body io.Reader
	if data != nil {
		js, _ := json.Marshal(data)
		body = bytes.NewReader(js)
	}
	req, err := http.NewRequestWithContext(c.ctx, method, fmt.Sprintf("%v%v", c.addr, route), body)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	r, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer func() {
		if e := r.Body.Close(); e != nil {
			err = multierr.Append(err, e)
		}
	}()
	if r.StatusCode != 200 {
		err, _ := ioutil.ReadAll(r.Body)
		return NewError(strings.TrimSpace(string(err)))
	}
	if resp == nil {
		return nil
	}
	return json.NewDecoder(r.Body).Decode(resp)
}

func (c *Client) get(route string, r interface{}) error     { return c.req("GET", route, nil, r) }
func (c *Client) post(route string, d, r interface{}) error { return c.req("POST", route, d, r) }
func (c *Client) put(route string, d, r interface{}) error  { return c.req("PUT", route, d, r) }

// WithContext returns a new Client whose requests are subject to the supplied
// context.
func (c *Client) WithContext(ctx context.Context) *Client {
	return &Client{
		addr: c.addr,
		ctx:  ctx,
	}
}

// AllContracts returns all contracts formed by the server.
func (c *Client) AllContracts() (cs []Contract, err error) {
	err = c.get("/contracts", &cs)
	return
}

// Contracts returns the contracts in the specified set.
func (c *Client) Contracts(set string) (cs []Contract, err error) {
	if set == "" {
		return nil, errors.New("no host set provided; to retrieve all contracts, use AllContracts")
	}
	err = c.get("/contracts?hostset="+set, &cs)
	return
}

// Scan queries the specified host for its current settings.
//
// Note that the host may also be scanned via the hostdb.Scan function.
func (c *Client) Scan(host hostdb.HostPublicKey) (settings hostdb.HostSettings, err error) {
	err = c.post("/scan", RequestScan{
		HostKey: host,
	}, &settings)
	return
}

// Form forms a contract with a host. The settings should be obtained from a
// recent call to Scan. If the settings have changed in the interim, the host
// may reject the contract.
func (c *Client) Form(host *hostdb.ScannedHost, funds types.Currency, start, end types.BlockHeight) (contract Contract, err error) {
	err = c.post("/form", RequestForm{
		HostKey:     host.PublicKey,
		Funds:       funds,
		StartHeight: start,
		EndHeight:   end,
		Settings:    host.HostSettings,
	}, &contract)
	return
}

// Renew renews the contract with the specified ID, which must refer to a
// contract previously formed by the server. The settings should be obtained
// from a recent call to Scan. If the settings have changed in the interim, the
// host may reject the contract.
func (c *Client) Renew(host *hostdb.ScannedHost, old *renter.Contract, funds types.Currency, start, end types.BlockHeight) (contract Contract, err error) {
	err = c.post("/renew", RequestRenew{
		ID:          old.ID,
		Funds:       funds,
		StartHeight: start,
		EndHeight:   end,
		Settings:    host.HostSettings,
		HostKey:     host.PublicKey,
		RenterKey:   old.RenterKey,
	}, &contract)
	return
}

// HostSets returns the current list of host sets.
func (c *Client) HostSets() (hs []string, err error) {
	err = c.get("/hostsets/", &hs)
	return
}

// HostSet returns the contents of the named host set.
func (c *Client) HostSet(name string) (hosts []hostdb.HostPublicKey, err error) {
	err = c.get("/hostsets/"+name, &hosts)
	return
}

// SetHostSet sets the contents of a host set, creating it if it does not exist.
// If an empty slice is passed, the host set is deleted.
func (c *Client) SetHostSet(name string, hosts []hostdb.HostPublicKey) (err error) {
	err = c.put("/hostsets/"+name, hosts, nil)
	return
}

// SHARD returns a client for the muse server's shard endpoints.
func (c *Client) SHARD() *shard.Client {
	u, err := url.Parse(c.addr)
	if err != nil {
		panic(err)
	}
	u.Path = path.Join(u.Path, "shard")
	return shard.NewClient(u.String())
}

// NewClient returns a client that communicates with a muse server listening
// on the specified address.
func NewClient(addr string) *Client {
	return &Client{addr, context.Background()}
}

func modifyURL(str string, fn func(*url.URL)) string {
	u, err := url.Parse(str)
	if err != nil {
		panic(err)
	}
	fn(u)
	return u.String()
}
