package virtualbox

import (
	"bytes"
	"context"
	"errors"
	"io"
	"io/ioutil"
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
	if p := os.Getenv("VBOX_INSTALL_PATH"); p != "" && runtime.GOOS == "windows" {
		VBM = filepath.Join(p, "VBoxManage.exe")
	}
}

var (
	reVMNameUUID      = regexp.MustCompile(`"(.+)" {([0-9a-f-]+)}`)
	reVMInfoLine      = regexp.MustCompile(`(?:"(.+)"|(.+))=(?:"(.*)"|(.*))`)
	reColonLine       = regexp.MustCompile(`(.+):\s+(.*)`)
	reMachineNotFound = regexp.MustCompile(`Could not find a registered machine named '(.+)'`)
)

var (
	ErrMachineExist    = errors.New("machine already exists")
	ErrMachineNotExist = errors.New("machine does not exist")
	ErrVBMNotFound     = errors.New("VBoxManage not found")
)

// executor abstracts the execution method that is being used to run the
// command.
type executor func(context.Context, io.Writer, io.Writer, ...string) error

var defaultExecutor executor = cmdExecutor

func cmdExecutor(ctx context.Context, so io.Writer, se io.Writer, args ...string) error {
	cmd := exec.CommandContext(ctx, VBM, args...)
	if Verbose {
		log.Printf("executing: %v %v", VBM, strings.Join(args, " "))
	}
	cmd.Stdout = so
	cmd.Stderr = se
	err := cmd.Run()
	if err != nil {
		if ee, ok := err.(*exec.Error); ok && ee == exec.ErrNotFound {
			err = ErrVBMNotFound
		}
	}
	return err
}

func vbm(args ...string) error {
	so, se := ioutil.Discard, ioutil.Discard
	if Verbose {
		so = os.Stdout
		se = os.Stderr
	}
	return defaultExecutor(context.Background(), so, se, args...)
}

func vbmOut(args ...string) (string, error) {
	so, se := new(bytes.Buffer), ioutil.Discard
	if Verbose {
		se = os.Stderr
	}
	err := defaultExecutor(context.Background(), so, se, args...)
	return so.String(), err
}

func vbmOutErr(args ...string) (string, string, error) {
	so, se := new(bytes.Buffer), new(bytes.Buffer)
	err := defaultExecutor(context.Background(), so, se, args...)
	return so.String(), se.String(), err
}
