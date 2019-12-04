package virtualbox

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestMachine(t *testing.T) {
	Setup(t)

	if ManageMock != nil {
		listVmsOut := ReadTestData("vboxmanage-list-vms-1.out")
		vmInfoOut := ReadTestData("vboxmanage-showvminfo-1.out")
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
