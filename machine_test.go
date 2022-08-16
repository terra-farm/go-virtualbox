package virtualbox

import (
	"context"
	"errors"
	"testing"

	"github.com/go-test/deep"
)

var (
	testUbuntuMachine = &Machine{
		Name:     "Ubuntu",
		Firmware: "BIOS",
		UUID:     "37f5d336-bf07-48dd-947c-37e6a56420a7",
		State:    Saved,
		CPUs:     1,
		Memory:   1024, VRAM: 8, CfgFile: "/Users/fix/VirtualBox VMs/go-virtualbox/go-virtualbox.vbox",
		BaseFolder: "/Users/fix/VirtualBox VMs/go-virtualbox", OSType: "", Flag: 0, BootOrder: []string{},
		NICs: []NIC{
			{Network: "nat", Hardware: "82540EM", HostInterface: "", MacAddr: "080027EE1DF7"},
		},
	}
	testGoVirtualboxMachine = &Machine{
		Name:     "go-virtualbox",
		Firmware: "BIOS",
		UUID:     "37f5d336-bf08-48dd-947c-37e6a56420a7",
		State:    Saved,
		CPUs:     1,
		Memory:   1024, VRAM: 8, CfgFile: "/Users/fix/VirtualBox VMs/go-virtualbox/go-virtualbox.vbox",
		BaseFolder: "/Users/fix/VirtualBox VMs/go-virtualbox", OSType: "", Flag: 0, BootOrder: []string{},
		NICs: []NIC{
			{Network: "nat", Hardware: "82540EM", HostInterface: "", MacAddr: "080027EE1DF7"},
		},
	}
)

func TestMachine(t *testing.T) {
	testCases := map[string]struct {
		in   string
		want *Machine
		err  error
	}{
		"by name": {
			in:   "Ubuntu",
			want: testUbuntuMachine,
			err:  nil,
		},
	}
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			m := newTestManager()

			got, err := m.Machine(context.Background(), tc.in)
			if diff := deep.Equal(got, tc.want); !errors.Is(err, tc.err) || diff != nil {
				t.Errorf("Machine(%s) = %+v, %v; want %v, %v; diff = %v",
					tc.in, got, err, tc.want, tc.err, diff)
			}
		})
	}
}

func TestListMachines(t *testing.T) {
	testCases := map[string]struct {
		want []*Machine
		err  error
	}{
		"good": {
			// TODO: If this relies on order we should ensure that it will be
			//       consistent for the tests.
			want: []*Machine{testUbuntuMachine, testGoVirtualboxMachine},
			err:  nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			m := newTestManager()

			got, err := m.ListMachines(context.Background())
			if diff := deep.Equal(got, tc.want); !errors.Is(err, tc.err) || diff != nil {
				t.Errorf("ListMachines() = %v, %v; want %v, %v; diff = %v",
					got, err, tc.want, tc.err, diff)
			}
		})
	}
}

func TestModifyMachine(t *testing.T) {
	// TODO: Figure out how we can do this test, it has pretty extensive flag list
	//       so having a file in the testdata with such a long name doesn't make
	//       sense.
}
