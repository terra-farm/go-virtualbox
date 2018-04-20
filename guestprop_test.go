package virtualbox

import (
	"os"
	"testing"
)

func TestGuestProperty(t *testing.T) {
	var vm = os.Getenv("TEST_VM")
	if len(vm) <= 0 {
		vm = "go-virtualbox"
		t.Logf("Missing TEST_VM environment variable")
	}
	t.Logf("Using '%s'", vm)

	err := SetGuestProperty(vm, "test_key", "test_val")
	if err != nil {
		t.Fatal(err)
	}

	val, err := GetGuestProperty(vm, "test_key")
	if err != nil {
		t.Fatal(err)
	}
	t.Logf("val='%s'", val)
	if val != "test_val" {
		t.Fatal("Wrong value")
	}
}
