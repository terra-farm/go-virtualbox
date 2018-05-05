package virtualbox

import (
	"os"
	"path/filepath"
	"regexp"
	"runtime"
)

var (
	// Manage holds the command to run VBoxManage.
	Manage Command
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

func init() {
	vbmgt := "VBoxManage"
	p := os.Getenv("VBOX_INSTALL_PATH")

	if p != "" && runtime.GOOS == "windows" {
		vbmgt = filepath.Join(p, "VBoxManage.exe")
	}
	//Trying fallback if nothing works
	if p == "" && runtime.GOOS == "windows" && vbmgt == "VBoxManage" {
		vbmgt = filepath.Join("C:\\", "Program Files", "Oracle", "VirtualBox", "VBoxManage.exe")
	}
	Manage = command{program: vbmgt}
}
