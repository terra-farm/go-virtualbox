PKGS := $(filter-out /vendor%,$(shell go list ./...))

SHELL = bash

INTERACTIVE:=$(shell [ -t 0 ] && echo 1)

ifdef INTERACTIVE
# is a terminal
else
# cron job / other
endif

default: test lint

.PHONY: test
test:
	go test $(PKGS)

# go get asks for credentials when needed
ifdef INTERACTIVE
GIT_TERMINAL_PROMPT := 1
export GIT_TERMINAL_PROMPT
endif

#GOMETALINTER := gometalinter.v2
GOMETALINTER := gometalinter

$(GOMETALINTER):
ifeq ($(GOMETALINTER),gometalinter)
	go get -u github.com/alecthomas/$(@)
else
	go get -u gopkg.in/alecthomas/$(@)
endif
	$(@) --install 1>/dev/null

.PHONY: lint
lint: $(GOMETALINTER)
	$(GOMETALINTER) ./... --vendor

BINARY := mytool

VERSION ?= $(shell git describe --tags)

PLATFORMS := windows linux darwin

os = $(word 1, $@)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	mkdir -p release
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: windows linux darwin