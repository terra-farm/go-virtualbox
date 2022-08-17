package virtualbox

import (
	"testing"

	"github.com/golang/mock/gomock"
)

func TestDHCPs(t *testing.T) {
	Setup(t)

	if ManageMock != nil {
		listDhcpServersOut := ReadTestData("vboxmanage-list-dhcpservers-1.out")
		gomock.InOrder(
			ManageMock.EXPECT().run("list", "dhcpservers").Return(listDhcpServersOut, "", nil).Times(1),
		)
	}
	m, err := DHCPs()
	if err != nil {
		t.Fatal(err)
	}

	for _, dhcp := range m {
		t.Logf("%+v", dhcp)
	}

	Teardown()
}
