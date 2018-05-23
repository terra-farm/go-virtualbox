package main

import (
	"os/exec"

	"github.com/asnowfix/go-virtualbox"
)

func open(args ...string) *exec.Cmd {
	argv := append([]string{"/c"}, "start")
	argv = append(argv, args...)
	virtualbox.Debug("Executing %v %v", "cmd", argv)
	return exec.Command("cmd", argv...) // #nosec
}
