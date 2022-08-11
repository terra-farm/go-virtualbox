package virtualbox

import (
	"context"
	"sync"
)

// runFn is the function which is used to actually run the commands. This is
// abstracted into a function so it can be easily replaced for testing purposes.
type runFn func(context.Context, ...string) (string, string, error)

// Manager of the virtualbox instance.
type Manager struct {
	// lock the whole manager to only allow one action at a time
	// TODO: Decide if this is a good idea, maybe one mutex per type of operation?
	lock sync.Mutex

	run runFn
}

// NewManager returns a manager capable of managing everything in virtualbox.
func NewManager() *Manager {
	return &Manager{
		run: vboxManageRun,
	}
}

// vboxManageRun is a function which actually runs the VboxManage
func vboxManageRun(_ context.Context, args ...string) (string, string, error) {
	// TODO: reimplement and do not use the old function
	return Manage().runOutErr(args...)
}

// defaultManager is used for backwards compatibility so that the older
// functions can use it.
var defaultManager = NewManager()
