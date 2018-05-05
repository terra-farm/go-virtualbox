package virtualbox

import (
	"os"
	"path/filepath"
	"runtime"
)

var (
	// Control holds the command to run VBoxControl.
	Control Command
)

func init() {
	vbctl := "VBoxControl"
	p := os.Getenv("VBOX_INSTALL_PATH")

	if p != "" && runtime.GOOS == "windows" {
		vbctl = filepath.Join(p, "VBoxControl.exe")
	}
	//Trying fallback if nothing works
	if p == "" && runtime.GOOS == "windows" && vbctl == "VBoxControl" {
		vbctl = filepath.Join("C:\\", "Program Files", "Oracle", "VirtualBox", "VBoxControl.exe")
	}
	Control = command{program: vbctl}
}
