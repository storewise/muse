module lukechampine.com/muse

go 1.15

require (
	github.com/BurntSushi/toml v0.3.1
	github.com/gorilla/handlers v1.5.1
	github.com/pkg/errors v0.9.1
	gitlab.com/NebulousLabs/Sia v1.4.11
	go.uber.org/multierr v1.6.0
	golang.org/x/crypto v0.0.0-20201016220609-9e8e0b390897
	lukechampine.com/flagg v1.1.1
	lukechampine.com/frand v1.3.0
	lukechampine.com/shard v0.3.4
	lukechampine.com/us v0.18.3
	lukechampine.com/walrus v0.10.0
)

replace (
	gitlab.com/NebulousLabs/errors => github.com/storewise/sia-errors v0.0.0-20201017234534-617267505fae
	lukechampine.com/muse => github.com/storewise/muse v0.0.0-20201007212310-2a7fd6751fc0
	lukechampine.com/us => github.com/storewise/us v0.18.4-0.20201016074118-080dbdd2d6db
)
