package virtualbox

import "testing"

var machineList = "\"Ubuntu\" {2e16b1fc-675d-4a7a-a9a1-e89a8bde7874}\n" +
	"\"go-virtualbox\" {def44546-e3da-4902-8d15-b91c99c80cbc}"

func TestMachine(t *testing.T) {
	ms, err := ListMachines()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range ms {
		t.Logf("%+v", m)
	}
}
