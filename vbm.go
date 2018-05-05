package virtualbox

import (
	"bytes"
	"errors"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	// VBM holds the inferred path to the VBoxManage utility.
	VBM string
	// Verbose when set toggle the library in verbose execution mode.
	Verbose bool
)

func init() {
	VBM = "VBoxManage"
	p := os.Getenv("VBOX_INSTALL_PATH")

	if p != "" && runtime.GOOS == "windows" {
		VBM = filepath.Join(p, "VBoxManage.exe")
	}
	//Trying fallback if nothing works
	if p == "" && runtime.GOOS == "windows" && VBM == "VBoxManage" {
		VBM = filepath.Join("C:\\", "Program Files", "Oracle", "VirtualBox", "VBoxManage.exe")
	}
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	// ErrMachineExist holds the error message when the machine already exists.
	ErrMachineExist = errors.New("machine already exists")
	// ErrMachineNotExist holds the error message when the machine does not exist.
	ErrMachineNotExist = errors.New("machine does not exist")
	// ErrVBMNotFound holds the error message when the VBoxManage commands was not found.
	ErrVBMNotFound = errors.New("VBoxManage not found")
)

var (
	// Manage holds the command to run VBoxManage
	Manage Command = manage{}
)

