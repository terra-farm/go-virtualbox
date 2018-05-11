package virtualbox

import (
	"os"
	"os/exec"
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
	sudoer, _ := isSudoer()

	if vbprog, err := lookupVBoxProgram("VBoxManage"); err == nil {
		Manage = command{program: vbprog, sudoer: sudoer, guest: false}
	} else if vbprog, err := lookupVBoxProgram("VBoxControl"); err == nil {
		Manage = command{program: vbprog, sudoer: sudoer, guest: true}
	} else {
		// Did not find a VirtualBox management command
		Manage = command{program: "false", sudoer: false, guest: false}
	}
	Debug("Manage: '%+v'", Manage)
}

func lookupVBoxProgram(vbprog string) (string, error) {

	if runtime.GOOS == "windows" {
		if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" {
			vbprog = filepath.Join(p, vbprog+".exe")
		} else {
			vbprog = filepath.Join("C:\\", "Program Files", "Oracle", "VirtualBox", vbprog+".exe")
		}
	}

	return exec.LookPath(vbprog)
}
