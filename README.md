# go-virtualbox

This is a wrapper package for Golang to interact with VirtualBox. The API is
experimental at the moment and you should expect frequent changes.

API doc at http://godoc.org/github.com/riobard/go-virtualbox

## Testing 

### Preparation

You either need to  have a pre-provisioned VirtualBox VM and  to set its name
using the `TEST_VM` environment variable, or use [Vagrant](https://www.vagrantup.com/intro/getting-started/).

```bash
$ vagrant box add bento/ubuntu-16.04
$ vagrant up
```

### Run tests

As usual, run `go test`, or `go test -v`.  To run one test in particular,
run `go test --run TestGuestProperty`.