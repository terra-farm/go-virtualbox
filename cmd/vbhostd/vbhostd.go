package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/terra-farm/go-virtualbox"
)

var (
	openRegexp = regexp.MustCompile("^(http|https|mailto):")
)

func main() {
	vm := flag.String("vm", "all", "VM to wait events from (all)")
	verbose := flag.Bool("v", false, "run in verbose mode")
	help := flag.Bool("h", false, "this message")
	flag.Parse()
	if *help {
		flag.PrintDefaults()
		os.Exit(0)
	}

	logger := log.New(os.Stderr, "", 0)
	virtualbox.Verbose = *verbose
	virtualbox.Debug = func(format string, args ...interface{}) {
		if !*verbose {
			return
		}
		msg := fmt.Sprintf(format, args...)
		logger.SetPrefix("\t  ")
		logger.Print(msg + "\n")
	}

	var vms []string
	if *vm == "all" {
		machines, err := virtualbox.ListMachines()
		if err != nil {
			panic(err)
		}
		if *verbose {
			virtualbox.Debug("machines: %+v\n", machines)
		}
		for _, machine := range machines {
			vms = append(vms, machine.Name)
		}
		if *verbose {
			virtualbox.Debug("vms: %+v\n", vms)
		}
	} else {
		vms = append(vms, *vm)
	}

	wg := new(sync.WaitGroup)
	agg := make(chan virtualbox.GuestProperty)
	done := make(map[string]chan bool)

	for _, vm := range vms {
		done[vm] = make(chan bool)
		props := virtualbox.WaitGuestProperties(vm, "vbhostd/*", done[vm], wg)
		go func(c chan virtualbox.GuestProperty) {
			for prop := range c {
				agg <- prop
			}
		}(props)
	}

	func() {
		for prop := range agg {
			virtualbox.Debug("Got prop: %+v.\n", prop)
			switch prop.Name {
			case "vbhostd/open":
				fmt.Printf("opening: %v\n", prop.Value)
				virtualbox.Debug("opening: %v", prop.Value)
				if openRegexp.MatchString(prop.Value) {
					args := strings.Split(prop.Value, " ")
					cmd := open(args...)
					err := cmd.Run()
					if err != nil {
						fmt.Printf("Error: %v\n", err)
					}
				} else {
					fmt.Printf("Error: not a supported URL=%v\n", prop.Value)
					virtualbox.Debug("Error: not a supported URL=%v", prop.Value)
				}
			case "vbhostd/error":
				fmt.Printf("Error: %v\n", prop.Value)
				virtualbox.Debug("Error: %v", prop.Value)
				return
			case "":
				fmt.Printf("Unexpected error: %v\n", prop.Value)
				virtualbox.Debug("Unexpected error: %v", prop.Value)
				return
			}
		}
	}()

	for vm, d := range done {
		virtualbox.Debug("Closing WaitGuestProperties(%s)...\n", vm)
		close(d)
		virtualbox.Debug("Closing WaitGuestProperties(%s)... Ok\n", vm)
	}

	virtualbox.Debug("Waiting completion or timeout...\n")
	wait := make(chan struct{})
	go func() {
		wg.Wait()
		close(wait)
	}()

	select {
	case <-wait:
		virtualbox.Debug("Every WaitGuestProperties() have completed.\n")
	case <-time.After(2000 * time.Millisecond):
		virtualbox.Debug("Timeout.\n")
	}

	fmt.Printf("Exiting....\n")
	virtualbox.Debug("Exiting....\n")
}
