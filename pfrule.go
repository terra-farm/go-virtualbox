package virtualbox

import "fmt"

// Port forwarding rule.
type PFRule struct {
	Proto     string // tcp|udp
	HostIP    string
	HostPort  uint16
	GuestIP   string
	GuestPort uint16
}

func (r PFRule) String() string {
	return fmt.Sprintf("%s://%s:%d --> %s:%d",
		r.Proto, r.HostIP, r.HostPort,
		r.GuestIP, r.GuestPort)
}
