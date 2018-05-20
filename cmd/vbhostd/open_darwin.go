package main

import (
	"os/exec"
)

func open(args ...string) *exec.Cmd {
	return exec.Command("open", args...)
}
