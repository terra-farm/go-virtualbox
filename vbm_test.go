package virtualbox

import (
	"context"
	"io"
	"os"
	"testing"
)

func TestVBMOut(t *testing.T) {
	b, err := vbmOut("list", "vms")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}

func setup() {
	Verbose = true

	defaultExecutor = mockExecutor
}

func TestMain(m *testing.M) {
	setup()
	os.Exit(m.Run())
}

func mockExecutor(ctx context.Context, so io.Writer, se io.Writer, args ...string) error {
	// TODO: By returning nil we are causing all the tests to pass because the
	//       current ones do not check the output of the command. Here we would
	//       keep the state of the machines, immitating VBoxManage - thus
	//       eliminating it as a dependency.
	return nil
}
