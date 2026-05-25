PKG := github.com/334456777/weflow-api/internal/version
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT  := $(shell git rev-parse --short HEAD 2>/dev/null || echo none)
DATE    := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -s -w -X $(PKG).Version=$(VERSION) -X $(PKG).Commit=$(COMMIT) -X $(PKG).Date=$(DATE)

.PHONY: build test tidy fmt vet snapshot release clean install

build:
	go build -ldflags "$(LDFLAGS)" -o weflow .

test:
	go test ./...

tidy:
	go mod tidy

fmt:
	gofmt -s -w .

vet:
	go vet ./...

snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean

install: build
	install -m 0755 weflow $${GOBIN:-$$HOME/go/bin}/weflow

clean:
	rm -rf dist weflow weflow.exe
