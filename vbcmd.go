package virtualbox

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"strings"
)

// Command is the mock-able interface to run VirtualBox commands
// such as VBoxManage (host side) or VBoxControl (guest side)
type Command interface {
	path() string
	run(args ...string) error
	runOut(args ...string) (string, error)
	runOutErr(args ...string) (string, string, error)
}

var (
	// Verbose toggles the library in verbose execution mode.
	Verbose bool
	// ErrMachineExist holds the error message when the machine already exists.
	ErrMachineExist = errors.New("machine already exists")
	// ErrMachineNotExist holds the error message when the machine does not exist.
	ErrMachineNotExist = errors.New("machine does not exist")
	// ErrCommandNotFound holds the error message when the VBoxManage commands was not found.
	ErrCommandNotFound = errors.New("command not found")
)

type command struct {
	program string
}

func (vbcmd command) path() string {
	return vbcmd.program
}

func (vbcmd command) run(args ...string) error {
	cmd := exec.Command(vbcmd.program, args...)
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", vbcmd.program, strings.Join(args, " "))
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrCommandNotFound
		}
		return err
	}
	return nil
}

func (vbcmd command) runOut(args ...string) (string, error) {
	cmd := exec.Command(vbcmd.program, args...)
	if Verbose {
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", vbcmd.program, strings.Join(args, " "))
	}

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrCommandNotFound
		}
	}
	return string(b), err
}

func (vbcmd command) runOutErr(args ...string) (string, string, error) {
	cmd := exec.Command(vbcmd.program, args...)
	if Verbose {
		log.Printf("executing: %v %v", vbcmd.program, strings.Join(args, " "))
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrCommandNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}
