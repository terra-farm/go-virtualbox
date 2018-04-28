package virtualbox

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

// GuestProperty holds key, value and associated flags.
type GuestProperty struct {
	Name  string
	Value string
}

var (
	getRegexp  = regexp.MustCompile("^Value: ([^,]*)$")
	waitRegexp = regexp.MustCompile("^Name: ([^,]*), value: ([^,]*), flags:.*$")
)

// SetGuestProperty writes a VirtualBox guestproperty to the given value.
func SetGuestProperty(vm string, prop string, val string) error {
	return Manage.run("guestproperty", "set", vm, prop, val)
}

// GetGuestProperty reads a VirtualBox guestproperty.
func GetGuestProperty(vm string, prop string) (string, error) {
	var out string
	var err error
	out, err = Manage.runOut("guestproperty", "get", vm, prop)
	if err != nil {
		log.Print(err)
		return "", err
	}
	out = strings.TrimSpace(out)
	if Verbose {
		log.Printf("out (trimmed): '%s'", out)
	}
	var match = getRegexp.FindStringSubmatch(out)
	if Verbose {
		log.Print("match:", match)
	}
	if len(match) != 2 {
		return "", fmt.Errorf("No match with VBoxManage get guestproperty output")
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
	if Verbose {
		log.Printf("WaitGuestProperty(): wait on '%s'", prop)
	}
	out, err = Manage.runOut("guestproperty", "wait", vm, prop)
	if err != nil {
		log.Print(err)
		return "", "", err
	}
	out = strings.TrimSpace(out)
	if Verbose {
		log.Printf("WaitGuestProperty(): out (trimmed): '%s'", out)
	}
	var match = waitRegexp.FindStringSubmatch(out)
	if Verbose {
		log.Print("WaitGuestProperty(): match:", match)
	}
	if len(match) != 3 {
		return "", "", fmt.Errorf("No match with VBoxManage wait guestproperty output")
	}
	return match[1], match[2], nil
}

func WaitGetProperties(vm string, propPattern string, propsChan *chan GuestProperty, doneC chan bool) {
	go func() {
		defer close(*propsChan)
		for {
			if Verbose {
				log.Printf("WaitGetProperties(): waiting for: '%s' changes", propPattern)
			}
			name, value, err := WaitGuestProperty(vm, propPattern)
			if err != nil {
				log.Printf("WaitGetProperties(): err=%v", err)
				return
			}
			prop := GuestProperty{name, value}
			select {
			case *propsChan <- prop:
				if Verbose {
					log.Printf("WaitGetProperties(): stacked: %+v", prop)
				}
			case done := <-doneC:
				if Verbose {
					log.Printf("WaitGetProperties(): done=%v", done)
				}
				if done {
					return
				}
			}
		}
	}()
}

// DeleteGuestProperty deletes a VirtualBox guestproperty.
func DeleteGuestProperty(vm string, prop string) error {
	return Manage.run("guestproperty", "delete", vm, prop)
}
