package virtualbox

import (
	"fmt"
	"testing"
)

func init() {
	Verbose = true
	VBM = "VBoxManage"
}

func TestExtra(t *testing.T) {
	e := SetExtra("ihaoyue-1.1", "vbox_graph_mode", "360x640-16")
	if e != nil {
		t.Fatal(e)
	}

	b, err := GetExtra("ihaoyue-1.1", "vbox_graph_mode")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
	fmt.Printf("Value: %s\n", b)
}
