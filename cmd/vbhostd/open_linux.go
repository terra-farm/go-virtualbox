package main

import (
	"os/exec"
)

func open(args ...string) {
	exec.Command("xdg_open", args...)
}
