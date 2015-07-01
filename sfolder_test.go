package virtualbox

import (
	"testing"
)

func TestSharedFolderParse(t *testing.T) {
	lines := []struct{ key, value string }{
		{"SharedFolderNameMachineMapping2", "GGG"},
		{"nic2", "hostonly"},
		{"nictype2", "virtio"},
		{"nicspeed2", "0"},
		{"SharedFolderNameMachineMapping1", "Users"},
		{"SharedFolderPathMachineMapping1", "/Users"},
		{"SharedFolderPathMachineMapping2", "/Users/Guest"},
	}

	var sf SharedFolders
	for _, l := range lines {
		if err := sf.parseProperty(l.key, l.value); err != nil {
			t.Error("Error parsing line: ", err)
		}
	}

	expected := []SharedFolder{
		{"GGG", "/Users/Guest"},
		{"Users", "/Users"},
	}

	list := sf.List()
	t.Log("Shared folders:", list)

	for i := range expected {
		if expected[i] != list[i] {
			t.Errorf("Error, different at %d, %v!=%v", i, expected[i], list[i])
		}
	}
}
