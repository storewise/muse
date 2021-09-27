LD_FLAGS := -X 'main.githash=`git rev-parse --short HEAD`' \
           -X 'main.builddate=`date`'
IMAGE := ghcr.io/storewise/muse/muse
TAG := $(shell date '+%Y%m%d%H%M')
bold := $(shell tput bold)
sgr0 := $(shell tput sgr0)

# all builds a binary with the current commit hash
all: bin/muse bin/musec
	@

# static is like all, but for static binaries
static:
	go install -ldflags "$(ldflags) -s -w -extldflags='-static'" -tags='timetzdata' ./cmd/...

# dev builds a binary with dev constants
dev:
	go install -ldflags "$(LD_FLAGS)" -tags='dev' ./cmd/...

test:
	go test -short ./...

test-long:
	go test -v -race ./...

bench:
	go test -v -run=XXX -bench=. ./...

lint:
	@golint ./...
	@golangci-lint run ./...

bin/muse bin/musec:
	@echo "$(bold)Building binaries$(sgr0)"
	@mkdir -p bin
	GOARCH=amd64 GOOS=linux go build -ldflags "$(LD_FLAGS)" -o bin ./cmd/...

build: bin/muse bin/musec
	@echo "$(bold)Building a Docker image$(sgr0)"
	docker build -t $(IMAGE):$(TAG) -f build/Dockerfile .

clean:
	@rm -rf bin/*

.PHONY: all static dev test test-long bench lint build clean
