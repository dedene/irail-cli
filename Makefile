.PHONY: build test lint fmt ci install clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -ldflags "-s -w -X github.com/dedene/irail-cli/internal/cmd.version=$(VERSION) -X github.com/dedene/irail-cli/internal/cmd.commit=$(COMMIT) -X github.com/dedene/irail-cli/internal/cmd.date=$(DATE)"

build:
	@mkdir -p bin
	go build $(LDFLAGS) -o bin/irail ./cmd/irail

test:
	go test -race -v ./...

lint:
	golangci-lint run

fmt:
	go fmt ./...
	goimports -w .

ci: fmt lint test build

install:
	go install $(LDFLAGS) ./cmd/irail

clean:
	rm -rf bin/
	rm -rf dist/
