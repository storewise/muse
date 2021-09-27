module lukechampine.com/muse

go 1.17

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gorilla/handlers v1.5.1
	github.com/pkg/errors v0.9.1
	gitlab.com/NebulousLabs/encoding v0.0.0-20200604091946-456c3dc907fe
	go.sia.tech/siad v1.5.7
	go.uber.org/multierr v1.7.0
	golang.org/x/term v0.0.0-20210421210424-b80969c67360
	lukechampine.com/flagg v1.1.1
	lukechampine.com/frand v1.4.2
	lukechampine.com/shard v0.3.9
	lukechampine.com/us v0.19.4
	lukechampine.com/walrus v0.10.8
)

require (
	filippo.io/edwards25519 v1.0.0-beta.2 // indirect
	github.com/aead/chacha20 v0.0.0-20180709150244-8b13a72661da // indirect
	github.com/dchest/threefish v0.0.0-20120919164726-3ecf4c494abf // indirect
	github.com/felixge/httpsnoop v1.0.1 // indirect
	github.com/gorilla/websocket v1.4.2 // indirect
	github.com/julienschmidt/httprouter v1.3.0 // indirect
	github.com/klauspost/cpuid v1.2.2 // indirect
	github.com/klauspost/reedsolomon v1.9.3 // indirect
	gitlab.com/NebulousLabs/bolt v1.4.4 // indirect
	gitlab.com/NebulousLabs/demotemutex v0.0.0-20151003192217-235395f71c40 // indirect
	gitlab.com/NebulousLabs/entropy-mnemonics v0.0.0-20181018051301-7532f67e3500 // indirect
	gitlab.com/NebulousLabs/errors v0.0.0-20200929122200-06c536cf6975 // indirect
	gitlab.com/NebulousLabs/fastrand v0.0.0-20181126182046-603482d69e40 // indirect
	gitlab.com/NebulousLabs/go-upnp v0.0.0-20210414172302-67b91c9a5c03 // indirect
	gitlab.com/NebulousLabs/log v0.0.0-20200604091839-0ba4a941cdc2 // indirect
	gitlab.com/NebulousLabs/merkletree v0.0.0-20200118113624-07fbf710afc4 // indirect
	gitlab.com/NebulousLabs/monitor v0.0.0-20191205095550-2b0fd3e1012a // indirect
	gitlab.com/NebulousLabs/persist v0.0.0-20200605115618-007e5e23d877 // indirect
	gitlab.com/NebulousLabs/ratelimit v0.0.0-20200811080431-99b8f0768b2e // indirect
	gitlab.com/NebulousLabs/siamux v0.0.0-20210409140711-e667c5f458e4 // indirect
	gitlab.com/NebulousLabs/threadgroup v0.0.0-20200608151952-38921fbef213 // indirect
	go.etcd.io/bbolt v1.3.6 // indirect
	go.uber.org/atomic v1.7.0 // indirect
	golang.org/x/crypto v0.0.0-20210322153248-0c34fe9e7dc2 // indirect
	golang.org/x/net v0.0.0-20210410081132-afb366fc7cd1 // indirect
	golang.org/x/sys v0.0.0-20210330210617-4fbd30eecc44 // indirect
	golang.org/x/text v0.3.6 // indirect
)

replace gitlab.com/NebulousLabs/errors => github.com/storewise/sia-errors v0.0.0-20201017234534-617267505fae
