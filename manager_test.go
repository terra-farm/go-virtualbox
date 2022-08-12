package virtualbox

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
)

func newTestManager() *Manager {
	m := NewManager(Logger(log.Default()))
	m.run = testDataRun

	return m
}

// testDataRun returns the test data for the given args
func testDataRun(_ context.Context, args ...string) (string, string, error) {
	// joined args create the file information
	name := filepath.Join("testdata", strings.Join(args, "_")+".out")
	data, err := os.ReadFile(name)
	if err != nil {
		return "", "", fmt.Errorf("unable to open testdata file %s: %v", name, err)
	}
	return string(data), "", nil
}
