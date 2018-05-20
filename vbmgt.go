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
	manage Command
)

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

// Manage returns the Command to run VBoxManage/VBoxControl.
func Manage() Command {
	if manage != nil {
		return manage
	}

	sudoer, _ := isSudoer()

	if vbprog, err := lookupVBoxProgram("VBoxManage"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: false}
	} else if vbprog, err := lookupVBoxProgram("VBoxControl"); err == nil {
		manage = command{program: vbprog, sudoer: sudoer, guest: true}
	} else {
		// Did not find a VirtualBox management command
		manage = command{program: "false", sudoer: false, guest: false}
	}
	Debug("manage: '%+v'", manage)
	return manage
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
	Debug("User: '%+v'", me)
	if groupIDs, err := me.GroupIds(); runtime.GOOS == "linux" {
		if err != nil {
			return false, err
		}
		Debug("groupIDs: '%+v'", groupIDs)
		for _, groupID := range groupIDs {
			group, err := user.LookupGroupId(groupID)
			if err != nil {
				return false, err
			}
			Debug("group: '%+v'", group)
			if group.Name == "sudo" {
				return true, nil
			}
		}
	}
	return false, nil
}
