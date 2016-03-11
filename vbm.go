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
	VBM     string // Path to VBoxManage utility.
	Verbose bool   // Verbose mode.
)

func init() {
	VBM = "VBoxManage"
	if p := os.Getenv("VBOX_MSI_INSTALL_PATH" /*"VBOX_INSTALL_PATH"*/); p != "" && runtime.GOOS == "windows" {
		//VBM = filepath.Join(p, "VBoxManage.exe")
		VBM = "hbx\\VBoxManage.exe"
	} else if runtime.GOOS == "windows" {
		VBM = "hbx\\VBoxManage.exe"
	}
	Verbose = false
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)

	reVMProperty          = regexp.MustCompile(`Value:\s+(.+)`)
	reVMPropertyEnumerate = regexp.MustCompile(`Name:\s+(.+), value:\s+(.+), timestamp:\s+(.+)`)

	reHdNotFound = regexp.MustCompile(`Could not find file for the medium '(.+)'`)
	reHdUUID     = regexp.MustCompile(`UUID:\s+([0-9a-f-]+)`)
	//reHdInfoLine = regexp.MustCompile(`((.+)):\s((.*))`)
	reHdInfoLine = regexp.MustCompile(`(.+):\s+(.*)`)
	reHdParent   = regexp.MustCompile(`Parent UUID:\s+(.+)`)
	reHdState    = regexp.MustCompile(`State:\s+(.+)`)
	reHdType     = regexp.MustCompile(`Type:\s+(.+)`)
	reHdLocation = regexp.MustCompile(`Location:\s+(.+)`)
	reHdFormat   = regexp.MustCompile(`Storage format:\s+(.+)`)
	reHdCap      = regexp.MustCompile(`\s+(.+) MBytes`)
	reHdSize     = regexp.MustCompile(`Size on disk:\s+([0-9]+) MBytes`)
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrMediumNotExist  = errors.New("medium does not exist")
	ErrVBMNotFound     = errors.New("VBoxManage not found")
)

func vbm(args ...string) error {
	cmd := exec.Command(VBM, args...)
	if Verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", VBM, strings.Join(args, " "))
	}
	if err := cmd.Run(); err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			return ErrVBMNotFound
		}
		return err
	}
	return nil
}

func vbmOut(args ...string) (string, error) {
	cmd := exec.Command(VBM, args...)
	if Verbose {
		cmd.Stderr = os.Stderr
		log.Printf("executing: %v %v", VBM, strings.Join(args, " "))
	}

	b, err := cmd.Output()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return string(b), err
}

func vbmOutErr(args ...string) (string, string, error) {
	cmd := exec.Command(VBM, args...)
	if Verbose {
		log.Printf("executing: %v %v", VBM, strings.Join(args, " "))
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}

func VirtualBoxCmd(arg string) (string, string, error) {
	var cmd *exec.Cmd
	if runtime.GOOS == "windows" {
		cmd = exec.Command("cmd", "/c", "\""+VBM+"\" "+arg)
	} else {
		cmd = exec.Command("sh", "-c", VBM+" "+arg)
	}
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return stdout.String(), stderr.String(), err
}
