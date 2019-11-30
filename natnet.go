package virtualbox

import (
	"bufio"
	"net"
	"strings"
)

// A NATNet defines a NAT network.
type NATNet struct {
	Name    string
	IPv4    net.IPNet
	IPv6    net.IPNet
	DHCP    bool
	Enabled bool
}

// NATNets gets all NAT networks in a  map keyed by NATNet.Name.
func NATNets() (map[string]NATNet, error) {
	out, err := Manage().runOut("list", "natnets")
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
		case "IP":
			n.IPv4.IP = net.ParseIP(val)
		case "Network":
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv4.Mask = ipnet.Mask
		case "IPv6 Prefix":
			// TODO: IPv6 CIDR parsing works fine on macOS, check on Windows
			// if val == "" {
			// 	continue
			// }
			// l, err := strconv.ParseUint(val, 10, 7)
			// if err != nil {
			// 	return nil, err
			// }
			// n.IPv6.Mask = net.CIDRMask(int(l), net.IPv6len*8)
			_, ipnet, err := net.ParseCIDR(val)
			if err != nil {
				return nil, err
			}
			n.IPv6.Mask = ipnet.Mask
		case "DHCP Enabled":
			n.DHCP = (val == stringYes)
		case "Enabled":
			n.Enabled = (val == stringYes)
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return m, nil
}
