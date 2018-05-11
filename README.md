# go-virtualbox

This is a wrapper package for Golang to interact with VirtualBox. The API is
experimental at the moment and you should expect frequent changes.

API doc at http://godoc.org/github.com/riobard/go-virtualbox

**Table of Contents**

<!-- TOC depthFrom:2 depthTo:4 -->

- [Status](#status)
- [Building](#building)
- [Testing](#testing)
    - [Preparation](#preparation)
    - [Run tests](#run-tests)
    - [Re-generate mock](#re-generate-mock)

<!-- /TOC -->

## Status

| Project | Status | Notes |
|---------|--------|-------|
| [CircleCI](https://circleci.com/gh/asnowfix/go-virtualbox) | [![CircleCI](https://circleci.com/gh/asnowfix/go-virtualbox.svg?style=svg)](https://circleci.com/gh/asnowfix/go-virtualbox) | |
| [TravisCI](https://travis-ci.org/asnowfix/go-virtualbox) | [![Build Status](https://travis-ci.org/asnowfix/go-virtualbox.svg?branch=master)](https://travis-ci.org/asnowfix/go-virtualbox) | |
| [Go Report Card](https://goreportcard.com/) | [![Go Report Card](https://goreportcard.com/badge/github.com/asnowfix/go-virtualbox?style=flat-square)](https://goreportcard.com/report/github.com/asnowfix/go-virtualbox) | scan  code with `gofmt`, `go vet`, `gocyclo`, `golint`, `ineffassign`, `license` and `misspell`. |
| [GoDoc](http://godoc.org) | [![Go Doc](https://img.shields.io/badge/godoc-reference-blue.svg?style=flat-square)](http://godoc.org/github.com/asnowfix/go-virtualbox) | |
| Release | [![Release](https://img.shields.io/github/release/asnowfix/go-virtualbox.svg?style=flat-square)](https://github.com/asnowfix/go-virtualbox/releases/latest) | |

## Usage

The part of the library that manages guest properties can run both from the Host and the Guest.

```go
    err := virtualbox.SetGuestProperty("MyVM", "test_key", "test_val")
    val, err := GetGuestProperty(VM, "test_key")
```

See [GoDoc](https://godoc.org/github.com/asnowfix/go-virtualbox) for full details.

## Building

First install dependencies

```bash
$ go get -v github.com/golang/dep/cmd/dep
$ dep ensure -v
```

Then build:

```bash
$ go build -v
```

## Documentation

Run `godoc -http=:6060 &` and then `open http://localhost:6060/pkg/github.com/asnowfix/go-virtualbox/`.

## Testing 

### Preparation

Tests run using mocked interfaces, unless the `TEST_VM` environment variable is set, in order to test against real VirtualBox. You either need to  have a pre-provisioned VirtualBox VM and  to set its name using the `TEST_VM` environment variable, or use [Vagrant](https://www.vagrantup.com/intro/getting-started/).

```bash
$ vagrant box add bento/ubuntu-16.04
# select virtualbox as provider

$ vagrant up
```

Then run the tests 

```bash
$ export TEST_VM=go-virtualbox
$ go test
```

...or (on Windows):

```shell
> set TEST_VM=go-virtualbox
> go test
```

Once you are done with testing, run `vagrant halt` to same resources.

### Run tests

As usual, run `go test`, or `go test -v`.  To run one test in particular,
run `go test --run TestGuestProperty`.



### Re-generate mock

```bash
mockgen -source=vbcmd.go -destination=mockvbcmd_test.go -package=virtualbox -mock_names=Command=MockCommand
```

## Caveats

> I _know_ this code is _awful_: it is my GoLang ramp-up playground.  I hope to
> make progress in GoLang & improve the code accordingly.
