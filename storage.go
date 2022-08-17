package virtualbox

// StorageController represents a virtualized storage controller.
type StorageController struct {
	SysBus      SystemBus
	Ports       uint // SATA port count 1--30
	Chipset     StorageControllerChipset
	HostIOCache bool
	Bootable    bool
}

// SystemBus represents the system bus of a storage controller.
type SystemBus string

const (
	// SysBusIDE when the storage controller provides an IDE bus.
	SysBusIDE = SystemBus("ide")
	// SysBusSATA when the storage controller provides a SATA bus.
	SysBusSATA = SystemBus("sata")
	// SysBusSCSI when the storage controller provides an SCSI bus.
	SysBusSCSI = SystemBus("scsi")
	// SysBusFloppy when the storage controller provides access to Floppy drives.
	SysBusFloppy = SystemBus("floppy")
	// SysBusSAS storage controller provides a SAS bus.
	SysBusSAS = SystemBus("sas")
	// SysBusUSB storage controller proveds an USB bus.
	SysBusUSB = SystemBus("usb")
	// SysBusPCIE storage controller proveds a PCIe bus.
	SysBusPCIE = SystemBus("pcie")
	// SysBusVirtio storage controller proveds a Virtio bus.
	SysBusVirtio = SystemBus("virtio")
)

// StorageControllerChipset represents the hardware of a storage controller.
type StorageControllerChipset string

const (
	// CtrlLSILogic when the storage controller emulates LSILogic hardware.
	CtrlLSILogic = StorageControllerChipset("LSILogic")
	// CtrlLSILogicSAS when the storage controller emulates LSILogicSAS hardware.
	CtrlLSILogicSAS = StorageControllerChipset("LSILogicSAS")
	// CtrlBusLogic when the storage controller emulates BusLogic hardware.
	CtrlBusLogic = StorageControllerChipset("BusLogic")
	// CtrlIntelAHCI when the storage controller emulates IntelAHCI hardware.
	CtrlIntelAHCI = StorageControllerChipset("IntelAHCI")
	// CtrlPIIX3 when the storage controller emulates PIIX3 hardware.
	CtrlPIIX3 = StorageControllerChipset("PIIX3")
	// CtrlPIIX4 when the storage controller emulates PIIX4 hardware.
	CtrlPIIX4 = StorageControllerChipset("PIIX4")
	// CtrlICH6 when the storage controller emulates ICH6 hardware.
	CtrlICH6 = StorageControllerChipset("ICH6")
	// CtrlI82078 when the storage controller emulates I82078 hardware.
	CtrlI82078 = StorageControllerChipset("I82078")
	// CtrlUSB storage controller emulates USB hardware.
	CtrlUSB = StorageControllerChipset("USB")
	// CtrlNVME storage controller emulates NVME hardware.
	CtrlNVME = StorageControllerChipset("NVMe")
	// CtrlVirtIO storage controller emulates VirtIO hardware.
	CtrlVirtIO = StorageControllerChipset("VirtIO")
)

// StorageMedium represents the storage medium attached to a storage controller.
type StorageMedium struct {
	Port      uint
	Device    uint
	DriveType DriveType
	Medium    string // none|emptydrive|<uuid>|<filename|host:<drive>|iscsi
}

// DriveType represents the hardware type of a drive.
type DriveType string

const (
	// DriveDVD when the drive is a DVD reader/writer.
	DriveDVD = DriveType("dvddrive")
	// DriveHDD when the drive is a hard disk or SSD.
	DriveHDD = DriveType("hdd")
	// DriveFDD when the drive is a floppy.
	DriveFDD = DriveType("fdd")
)

// CloneHD virtual harddrive
func CloneHD(input, output string) error {
	_, _, err := Manage().run("clonehd", input, output)
	return err
}
