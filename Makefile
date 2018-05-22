## Setup -- do not modify

ME := $(shell id -un)
$(info ME=$(ME))
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

ifeq ($(GOROOT),)
$(error GOROOT undefined)
endif

ifeq ($(GOPATH),)
$(error GOPATH undefined)
endif

## Package -- 

ifeq ($(ME),vagrant)
# FIXME: any better way?
ROOTPKG := github.com/asnowfix/go-virtualbox
PKGS := $(filter-out /vendor%,$(shell cd $(GOPATH)/src/$(ROOTPKG) && go list ./...))
else
ROOTPKG := $(shell go list .)
PKGS := $(filter-out /vendor%,$(shell go list ./...))
endif
$(info PKGS=$(PKGS))

default: deps test lint build-pkgs

DEP_NAME := dep
DEP := $(GOPATH)/bin/$(DEP_NAME)$(EXE)

.PHONY: deps
deps: $(DEP)
	go get -t -d -v ./...
ifeq ($(ME),vagrant)
	cd $(GOPATH)/src/$(ROOTPKG) && $(DEP) ensure -v
else
	$(DEP) ensure -v
endif

$(DEP):
	go get -v github.com/golang/dep/cmd/dep

.PHONY: build test
build test:
	go $(@) -v ./...

.PHONY: build-pkgs
build-pkgs: $(foreach pkg,$(PKGS),build-pkg-$(basename $(pkg)))
#	go build -v ./cmd/vbhostd

define build-pkg
build-pkg-$(basename $(1)):
	go build -v $(1)
endef

$(foreach pkg,$(PKGS),$(eval $(call build-pkg,$(pkg))))

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
	GOOS=$(os) GOARCH=amd64 go build -o release/$(BINARY)-$(VERSION)-$(os)-amd64

.PHONY: release
release: $(PLATFORMS)