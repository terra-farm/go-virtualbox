package virtualbox

// Network interface card.
type NIC struct {
	Network         string // none|null|nat|bridged|intnet|hostonly|generic
	Hardware        string // Am79C970A|Am79C973|82540EM|82543GC|82545EM|virtio
	HostonlyAdapter string
}
