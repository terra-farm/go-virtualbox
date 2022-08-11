package virtualbox

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

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
