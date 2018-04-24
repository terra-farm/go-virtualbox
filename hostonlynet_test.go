package virtualbox

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestHostonlyNets(t *testing.T) {
	Setup(t)

	if ManageMock != nil {
		listHostOnlyIfsOut := ReadTestData("vboxmanage-list-hostonlyifs-1.out")
		gomock.InOrder(
			ManageMock.EXPECT().runOut("list", "hostonlyifs").Return(listHostOnlyIfsOut, nil).Times(1),
		)
	}
	m, err := HostonlyNets()
	if err != nil {
		t.Fatal(err)
	}
	for _, n := range m {
		t.Logf("%+v", n)
	}

	Teardown()
}
