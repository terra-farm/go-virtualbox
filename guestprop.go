package virtualbox

import (
	"fmt"
	"log"
	"regexp"
	"strings"
	"sync"
)

// GuestProperty holds key, value and associated flags.
type GuestProperty struct {
	Name  string
	Value string
}

var (
	getRegexp  = regexp.MustCompile("(?m)^Value: ([^,]*)$")
	waitRegexp = regexp.MustCompile("^Name: ([^,]*), value: ([^,]*), flags:.*$")
)

// SetGuestProperty writes a VirtualBox guestproperty to the given value.
func SetGuestProperty(vm string, prop string, val string) error {
	if Manage().isGuest() {
		return Manage().setOpts(sudo(true)).run("guestproperty", "set", prop, val)
	}
	return Manage().run("guestproperty", "set", vm, prop, val)
}

// GetGuestProperty reads a VirtualBox guestproperty.
func GetGuestProperty(vm string, prop string) (string, error) {
	var out string
	var err error
	if Manage().isGuest() {
		out, err = Manage().setOpts(sudo(true)).runOut("guestproperty", "get", prop)
	} else {
		out, err = Manage().runOut("guestproperty", "get", vm, prop)
	}
	if err != nil {
		return "", err
	}
	out = strings.TrimSpace(out)
	Debug("out (trimmed): '%s'", out)
	var match = getRegexp.FindStringSubmatch(out)
	Debug("match:", match)
	if len(match) != 2 {
		return "", fmt.Errorf("No match with get guestproperty output")
	}
	return match[1], nil
}

// WaitGuestProperty blocks until a VirtualBox guestproperty is changed
//
// The key to wait for can be a fully defined key or a key wild-card (glob-pattern).
// The first returned value is the property name that was changed.
// The second returned value is the new property value,
// Deletion of the guestproperty causes WaitGuestProperty to return the
// string.
func WaitGuestProperty(vm string, prop string) (string, string, error) {
	var out string
	var err error
	Debug("WaitGuestProperty(): wait on '%s'", prop)
	if Manage().isGuest() {
		out, err = Manage().setOpts(sudo(true)).runOut("guestproperty", "wait", prop)
	}
	out, err = Manage().runOut("guestproperty", "wait", vm, prop)
	if err != nil {
		log.Print(err)
		return "", "", err
	}
	out = strings.TrimSpace(out)
	Debug("WaitGuestProperty(): out (trimmed): '%s'", out)
	var match = waitRegexp.FindStringSubmatch(out)
	Debug("WaitGuestProperty(): match:", match)
	if len(match) != 3 {
		return "", "", fmt.Errorf("No match with VBoxManage wait guestproperty output")
	}
	return match[1], match[2], nil
}

// WaitGuestProperties wait for changes in GuestProperties
//
// WaitGetProperties wait for changes in the VirtualBox GuestProperties matching
// the given propsPattern, for the given VM.  The given bool channel indicates
// caller-required closure.  The optional sync.WaitGroup enabke the caller program
// to wait for Go routine completion.
//
// It returns a channel of GuestProperty objects (name-values pairs) populated
// as they change.
//
// If the bool channel is never closed, the Waiter Go routine never ends,
// but on VBoxManage error.
//
// Each GuestProperty change must be read from thwe channel before the waiter Go
// routine resumes waiting for the next matching change.
//
func WaitGuestProperties(vm string, propPattern string, done chan bool, wg *sync.WaitGroup) chan GuestProperty {

	props := make(chan GuestProperty)
	wg.Add(1)

	go func() {
		defer close(props)
		defer wg.Done()

		for {
			Debug("WaitGetProperties(): waiting for: '%s' changes", propPattern)
			name, value, err := WaitGuestProperty(vm, propPattern)
			if err != nil {
				log.Printf("WaitGetProperties(): err=%v", err)
				return
			}
			prop := GuestProperty{name, value}
			select {
			case props <- prop:
				Debug("WaitGetProperties(): stacked: %+v", prop)
			case <-done:
				Debug("WaitGetProperties(): done channel closed")
				return
			}
		}
	}()

	return props
}

// DeleteGuestProperty deletes a VirtualBox guestproperty.
func DeleteGuestProperty(vm string, prop string) error {
	if Manage().isGuest() {
		return Manage().setOpts(sudo(true)).run("guestproperty", "delete", prop)
	}
	return Manage().run("guestproperty", "delete", vm, prop)
}
