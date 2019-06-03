.PHONY: all

BIN := openair-apiserver
PKG := github.com/openairtech/apiserver

BINDIR = bin

VERSION_VAR := cmd.Version
TIMESTAMP_VAR := cmd.Timestamp

VERSION ?= $(shell git describe --always --dirty --tags)
TIMESTAMP := $(shell date -u '+%Y-%m-%d_%I:%M:%S%p')

GOBUILD_LDFLAGS := -ldflags "-s -w -X $(VERSION_VAR)=$(VERSION) -X $(TIMESTAMP_VAR)=$(TIMESTAMP)"

default: all

all: build

build:
	go build -x $(GOBUILD_LDFLAGS) -v -o $(BINDIR)/$(BIN)

build-static:
	env CGO_ENABLED=0 GOOS=linux go build -a -installsuffix "static" $(GOBUILD_LDFLAGS) -o $(BINDIR)/$(BIN)

clean:
	rm -dRf $(BINDIR)

deploy: build-static
	scp $(BINDIR)/$(BIN) root@openair.city:/usr/bin && ssh root@openair.city 'systemctl restart openair-apiserver'

fmt:
	go fmt ./...

# https://golangci.com/
# curl -sfL https://install.goreleaser.com/github.com/golangci/golangci-lint.sh | sh -s -- -b $GOPATH/bin v1.10.2
lint:
	golangci-lint run

test:
	go test ./...
