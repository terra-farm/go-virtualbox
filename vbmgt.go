package virtualbox

import (
	"os"
	"os/exec"
	"os/user"
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
	vbprogs := []string{"VBoxManage", "VBoxControl"}
	for _, vbprog := range vbprogs {
		vbprog, err := lookupVBoxProgram(vbprog)
		if err == nil {
			Manage = command{program: vbprog, sudoer: sudoer}
			break
		}
	}
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

func isSudoer() (bool, error) {
	me, err := user.Current()
	if err != nil {
		return false, err
	}
	if groupIDs, err := me.GroupIds(); runtime.GOOS == "linux" {
		if err != nil {
			return false, err
		}
		for _, groupID := range groupIDs {
			if groupID == "sudo" {
				return true, nil
			}
		}
	}
	return false, nil
}
