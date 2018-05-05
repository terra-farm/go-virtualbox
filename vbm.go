package virtualbox

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var (
	// VBM holds the inferred path to the VBoxManage utility.
	VBM string
	// Verbose when set toggle the library in verbose execution mode.
	Verbose bool
	// Manage holds the command to run VBoxManage
	Manage Command
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
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
	Manage = command{program: VBM}
}
