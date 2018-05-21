PKGS := $(filter-out /vendor%,$(shell go list ./...))
$(info PKGS=$(PKGS))

INTERACTIVE:=$(shell [ -t 0 ] && echo 1)
$(info INTERACTIVE=$(INTERACTIVE))
ifdef INTERACTIVE
# is a terminal
else
# cron job / other
endif

ifeq ($(OS),Windows_NT)
EXE := .exe
endif

default: deps test lint

DEP_NAME := dep
DEP := $(GOPATH)/bin/$(DEP_NAME)$(EXE)

.PHONY: deps
deps: $(DEP)
	go get -t -d -v ./...
	$(DEP) ensure -v

$(DEP):
	go get -v github.com/golang/dep/cmd/dep

.PHONY: build test
build test:
	go $(@) -v ./...

# go get asks for credentials when needed
ifdef INTERACTIVE
GIT_TERMINAL_PROMPT := 1
export GIT_TERMINAL_PROMPT
endif

#GOMETALINTER_NAME := gometalinter.v2
GOMETALINTER_NAME := gometalinter
GOMETALINTER := $(GOPATH)/bin/$(GOMETALINTER_NAME)$(EXE)

$(GOMETALINTER):
#	@echo PATH=$(PATH)
ifeq ($(GOMETALINTER_NAME),gometalinter)
	go get -u github.com/alecthomas/$(GOMETALINTER_NAME)
else
	go get -u gopkg.in/alecthomas/$(GOMETALINTER_NAME)
endif
	$(@) --install

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