module lukechampine.com/muse

go 1.16

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

replace gitlab.com/NebulousLabs/errors => github.com/storewise/sia-errors v0.0.0-20201017234534-617267505fae
