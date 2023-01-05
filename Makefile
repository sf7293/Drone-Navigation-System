# Caution: Read https://stackoverflow.com/questions/16931770/makefile4-missing-separator-stop article about very common Makefile mistakes
# Versioning
CURRENT_BRANCH = $(shell git rev-parse --abbrev-ref HEAD)
COMMIT_SHORT_HASH = $(shell git rev-parse --short HEAD)
COMMIT_HASH = $(shell git rev-parse HEAD)
DATE = $(shell date -u +%Y.%m.%d-%H%M%S)
VERSION = v$(DATE)-$(COMMIT_SHORT_HASH)

# Build the binary statically to avoid alpine libmusl's incompatibility issues
LDFLAGS = -ldflags '-extldflags "-static" -X app.Version=$(VERSION) -X app.CommitHash=$(COMMIT_HASH)'
BUILDFLAGS = -a -installsuffix cgo
PKG_BASE = github.com/sf7293/Drone-Navigation-System

.NOTPARALLEL:
.PHONY: all build test clean docker_test docker_base docker_lite

all: build test

clean:
	go clean -cache
	rm -rf release

build_server:
#	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build $(BUILDFLAGS) $(LDFLAGS) -o ./release/server $(PKG_BASE)
	go build -o ./release/server $(PKG_BASE)

build: build_server

push: tag
	git push origin $(VERSION)

tag:
	git tag $(VERSION)
	echo "Tagged $(VERSION)"

test:
	go test -v ./...

docker:
	docker build --no-cache -f Dockerfile.server -t sf7293/Drone-Navigation-System:$(VERSION) .
