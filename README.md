# go-virtualbox

This is a wrapper package for Golang to interact with VirtualBox. The API is
experimental at the moment and you should expect frequent changes.

API doc at http://godoc.org/github.com/riobard/go-virtualbox

## Testing 

### Preparation

Not every tests are curently using mocked interfaces, so
you either need to  have a pre-provisioned VirtualBox VM and  to set its name
using the `TEST_VM` environment variable, or use [Vagrant](https://www.vagrantup.com/intro/getting-started/).

```bash
$ vagrant box add bento/ubuntu-16.04
# select virtualbox as provider

$ vagrant up
```

Once you are done with testing, run `vagrant halt` to same resources.

### Run tests

As usual, run `go test`, or `go test -v`.  To run one test in particular,
run `go test --run TestGuestProperty`.

In order to activate the `VBoxManage` mocked stubs, set the `TEST_MOCK_VBM` environment variable to a non-empty version.
