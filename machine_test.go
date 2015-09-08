package virtualbox

import (
	"fmt"
	"log"
	"testing"
)

func TestMachine(t *testing.T) {
	ms, err := ListMachines()
	if err != nil {
		t.Fatal(err)
	}
	for _, m := range ms {
		t.Logf("%+v", m)

		if m.Name == "ihaoyue" {
			if n, ok := m.NICs[1]; ok {
				fmt.Println("MAC is : ", n.MACAddress)
			} else {
				fmt.Println("MAC is nof found, map len ", len(m.NICs))
			}

			fmt.Println("Machine: ", m)
			log.Println("NIC1: ", m.NICs[1])
			log.Println("NICNetwork: ", m.NICs[1].Network)
			log.Println("NICHardware: ", m.NICs[1].Hardware)
			log.Println("InterfaceName: ", m.NICs[1].InterfaceName)
			log.Println("InterfaceName: ", m.NICs[2].InterfaceName)
		}

	}

}
