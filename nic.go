package virtualbox

// Network interface card.
type NIC struct {
	Type            string // none|null|nat|bridged|intnet|hostonly|generic
	HwType          string // Am79C970A|Am79C973|82540EM|82543GC|82545EM|virtio
	HostonlyAdapter string
}
