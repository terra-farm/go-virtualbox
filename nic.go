package virtualbox

// NIC represents a virtualized network interface card.
type NIC struct {
	Network         NICNetwork
	Hardware        NICHardware
	HostonlyAdapter string
}

// NICNetwork represents the type of NIC networks.
type NICNetwork string

const (
	// NICNetAbsent when there is no NIC.
	NICNetAbsent = NICNetwork("none")
	// NICNetDisconnected when the NIC is disconnected
	NICNetDisconnected = NICNetwork("null")
	// NICNetNAT when the NIC is NAT-ed to access the external network.
	NICNetNAT = NICNetwork("nat")
	// NICNetBridged when the NIC is the bridge to the external network.
	NICNetBridged = NICNetwork("bridged")
	// NICNetInternal when the NIC does not have access to the external network.
	NICNetInternal = NICNetwork("intnet")
	// NICNetHostonly when the NIC can only access one host-only network.
	NICNetHostonly = NICNetwork("hostonly")
	// NICNetGeneric when the NIC behaves like a standard physical one.
	NICNetGeneric = NICNetwork("generic")
)

// NICHardware represents the type of NIC hardware.
type NICHardware string

const (
	// AMDPCNetPCIII when the NIC emulates a Am79C970A hardware.
	AMDPCNetPCIII = NICHardware("Am79C970A")
	// AMDPCNetFASTIII when the NIC emulates a Am79C973 hardware.
	AMDPCNetFASTIII = NICHardware("Am79C973")
	// IntelPro1000MTDesktop when the NIC emulates a 82540EM hardware.
	IntelPro1000MTDesktop = NICHardware("82540EM")
	// IntelPro1000TServer when the NIC emulates a 82543GC hardware.
	IntelPro1000TServer = NICHardware("82543GC")
	// IntelPro1000MTServer when the NIC emulates a 82545EM hardware.
	IntelPro1000MTServer = NICHardware("82545EM")
	// VirtIO when the NIC emulates a virtio.
	VirtIO = NICHardware("virtio")
)
