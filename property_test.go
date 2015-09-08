package virtualbox

import (
	_ "fmt"
	"testing"
)

func init() {
	Verbose = true
	VBM = "VBoxManage"
}

func TestProperty(t *testing.T) {
	e := guestPropertySet("ihaoyue-1.1", "vbox_graph_mode", "360x640-16")
	if e != nil {
		t.Fatal(e)
	}

	b, err := guestPropertyGet("ihaoyue-1.1", "vbox_graph_mode")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("%s", b)
	//fmt.Printf("Value: %s\n", b)

	ba, eerr := guestPropertyEnumerate("ihaoyue-1.1")
	if eerr != nil {
		t.Fatal(eerr)
	}
	t.Logf("%s", ba)
	//fmt.Println("Value: ", ba)
}
