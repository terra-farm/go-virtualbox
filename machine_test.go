package virtualbox

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMachine(t *testing.T) {
	Setup(t)

	if ManageMock != nil {
		listVmsOut := "\"Ubuntu\" {2e16b1fc-675d-4a7a-a9a1-e89a8bde7874}\n" +
			"\"go-virtualbox\" {def44546-e3da-4902-8d15-b91c99c80cbc}"
		vmInfoOut := ReadTestData("vboxmanage-showvminfo-1.properties")
		gomock.InOrder(
			ManageMock.EXPECT().runOut("list", "vms").Return(listVmsOut, nil).Times(1),
			ManageMock.EXPECT().runOutErr("showvminfo", "Ubuntu", "--machinereadable").Return(vmInfoOut, "", nil).Times(1),
			ManageMock.EXPECT().runOutErr("showvminfo", "go-virtualbox", "--machinereadable").Return(vmInfoOut, "", nil).Times(1),
		)
	}
	ms, err := ListMachines()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range ms {
		t.Logf("%+v", m)
	}

	Teardown()
}
