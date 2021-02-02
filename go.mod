module lukechampine.com/muse

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gorilla/handlers v1.5.1
	github.com/pkg/errors v0.9.1
	gitlab.com/NebulousLabs/Sia v1.5.4
	gitlab.com/NebulousLabs/encoding v0.0.0-20200604091946-456c3dc907fe
	go.uber.org/multierr v1.6.0
	golang.org/x/crypto v0.0.0-20201221181555-eec23a3978ad
	lukechampine.com/flagg v1.1.1
	lukechampine.com/frand v1.3.0
	lukechampine.com/shard v0.3.5
	lukechampine.com/us v0.19.0
	lukechampine.com/walrus v0.10.2
)

replace (
	gitlab.com/NebulousLabs/errors => github.com/storewise/sia-errors v0.0.0-20201017234534-617267505fae
	lukechampine.com/us => github.com/storewise/us v0.18.4-0.20210201204520-65324afd1bcf
)
