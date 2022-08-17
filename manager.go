package virtualbox

import (
	"context"
	"io"
	"log"
	"os"
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

	log *log.Logger
}

// NewManager returns a manager capable of managing everything in virtualbox.
func NewManager(opts ...Option) *Manager {
	m := &Manager{
		run: vboxManageRun,
		log: log.New(io.Discard, "", 0),
	}

	// if the debug env var for the virtualbox is set to true, we want to set the
	// logger to be bit more useful, and the default logger will suffice.
	if os.Getenv("DEBUG") == "virtualbox" {
		m.log = log.Default()
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}

// vboxManageRun is a function which actually runs the VboxManage
func vboxManageRun(_ context.Context, args ...string) (string, string, error) {
	// TODO: reimplement and do not use the old function
	return Manage().run(args...)
}

// defaultManager is used for backwards compatibility so that the older
// functions can use it.
var defaultManager = NewManager()

// Option modifies the manager options
type Option func(*Manager)

// Logger allows to override the logger used by the manager.
func Logger(l *log.Logger) Option {
	return func(m *Manager) {
		m.log = l
	}
}
