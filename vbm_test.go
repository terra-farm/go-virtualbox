package virtualbox

import (
	"testing"
)

func init() {
	Verbose = true
	VBM = "VBoxManage"
}

func TestVBMOut(t *testing.T) {
	b, err := vbmOut("list", "vms")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
}
