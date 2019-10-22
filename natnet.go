package virtualbox

import (
	"bufio"
	"net"
	"strings"
	"errors"
)
var (
		ErrNoSuchNet = errors.New("No such NAT network")
)
// A NATNet defines a NAT network.
type NATNet struct {
	Name    string
	IPv4    net.IPNet
	IPv6    net.IPNet
	DHCP    bool
	Enabled bool
}

func (n *NATNet) Delete() error {
	err := vbm("natnetwork", "remove", "--netname", n.Name)
	if err != nil {
		return err
	}
	return nil
}

func (n *NATNet) Config() error {
	if n.IPv4.IP != nil && n.IPv4.Mask != nil {
		if err := vbm("natnetwork", "modify", "--netname", n.Name, "--network", n.IPv4.String()); err != nil {
			return err
		}
	}

	if err := vbm("natnetwork", "modify", "--netname", n.Name, "--dhcp", bool2string(n.DHCP)); err != nil {
		return err
	}

	if n.Enabled {
		if err := vbm("natnetwork", "modify", "--netname", n.Name, "--enable"); err != nil {
			return err
		}
	} else {
		if err := vbm("natnetwork", "modify", "--netname", n.Name, "--disable"); err != nil {
			return err
		}
	}

	return nil
}

func CreateNATNet(name string, network string, dhcp bool) (*NATNet, error) {
	err := vbm("natnetwork", "add", "--netname", name, "--network", network, "--dhcp", bool2string(dhcp))
	if err != nil {
		return nil, err
	}
	_, ipnet, err := net.ParseCIDR(network)
	return &NATNet{Name: name, IPv4: *ipnet, DHCP: dhcp, Enabled: true}, nil
}

func GetNATNetwork(name string) (*NATNet, error) {
	natnets, err := NATNets()
	if err != nil {
		return nil, err
	}
	natnet, ok := natnets[name]
	if !ok {
		return nil, ErrNoSuchNet
	}
	return &natnet, nil
}

// NATNets gets all NAT networks in a  map keyed by NATNet.Name.
func NATNets() (map[string]NATNet, error) {
	out, err := vbmOut("list", "natnets")
	if err != nil {
		return nil, err
	}
	s := bufio.NewScanner(strings.NewReader(out))
	m := map[string]NATNet{}
	n := NATNet{}
	for s.Scan() {
		line := s.Text()
		if line == "" {
			m[n.Name] = n
			n = NATNet{}
			continue
		}
		res := reColonLine.FindStringSubmatch(line)
		if res == nil {
			continue
		}
		switch key, val := res[1], res[2]; key {
		case "NetworkName":
			n.Name = val
		case "Network":
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv4.IP = ipnet.IP
			n.IPv4.Mask = ipnet.Mask
		case "IPv6 Prefix":
			if val == "" {
				continue
			}
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv6.IP = ipnet.IP
			n.IPv6.Mask = ipnet.Mask
		case "DHCP Enabled":
			n.DHCP = (val == "Yes")
		case "Enabled":
			n.Enabled = (val == "Yes")
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
