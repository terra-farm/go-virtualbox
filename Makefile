PKGS := $(filter-out /vendor%,$(shell go list ./...))

SHELL = bash

.PHONY: test
test:
	go test $(PKGS)
#	$(MAKE) lint

BIN_DIR := $(GOPATH)/bin
GOMETALINTER := $(BIN_DIR)/gometalinter

$(GOMETALINTER):
	go get -u github.com/alecthomas/gometalinter
	gometalinter --install 1>/dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	gometalinter ./... --vendor

BINARY := mytool
VERSION ?= vlatest
PLATFORMS := windows linux darwin
os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: windows linux darwin