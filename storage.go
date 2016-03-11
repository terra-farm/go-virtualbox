package virtualbox

import (
	"errors"
	"fmt"
	gouuid "github.com/satori/go.uuid"
)

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
	SysBusIDE    = SystemBus("ide")
	SysBusSATA   = SystemBus("sata")
	SysBusSCSI   = SystemBus("scsi")
	SysBusFloppy = SystemBus("floppy")
)

// StorageControllerChipset represents the hardware of a storage controller.
type StorageControllerChipset string

const (
	CtrlLSILogic    = StorageControllerChipset("LSILogic")
	CtrlLSILogicSAS = StorageControllerChipset("LSILogicSAS")
	CtrlBusLogic    = StorageControllerChipset("BusLogic")
	CtrlIntelAHCI   = StorageControllerChipset("IntelAHCI")
	CtrlPIIX3       = StorageControllerChipset("PIIX3")
	CtrlPIIX4       = StorageControllerChipset("PIIX4")
	CtrlICH6        = StorageControllerChipset("ICH6")
	CtrlI82078      = StorageControllerChipset("I82078")
)

// StorageMedium represents the storage medium attached to a storage controller.
type StorageMedium struct {
	Port      uint
	Device    uint
	DriveType DriveType
	Medium    string // none|emptydrive|<uuid>|<filename|host:<drive>|iscsi
	Hdd       *MediumInfo
	Ctl       string
}

func (m *StorageMedium) Attach(hdd string) error {
	if m.Hdd != nil {
		_, eo, e := vbmOutErr("storageattach", m.Hdd.VMId, "--storagectl", m.Ctl,
			"--port", fmt.Sprintf("%d", m.Port),
			"--device", fmt.Sprintf("%d", m.Device),
			"--type", string(m.DriveType),
			"--medium", hdd,
		)
		if e != nil {
			return errors.New(e.Error() + ":" + eo)
		}
	}
	return nil
}

func (m *StorageMedium) AttachByVM(vm, hdd string) error {
	huuid := gouuid.NewV4().String()
	_, eo, e := vbmOutErr("internalcommands", "sethduuid", hdd, huuid)
	if e != nil {
		return errors.New(e.Error() + ":" + eo)
	}
	_, eo, e = vbmOutErr("storageattach", vm, "--storagectl", m.Ctl,
		"--port", fmt.Sprintf("%d", m.Port),
		"--device", fmt.Sprintf("%d", m.Device),
		"--type", string(m.DriveType),
		"--medium", hdd,
	)
	if e != nil {
		return errors.New(e.Error() + ":" + eo)
	}
	return nil
}

func (m *StorageMedium) Remove() error {
	if m.Hdd != nil {
		_, eo, e := vbmOutErr("storageattach", m.Hdd.VMId, "--storagectl", m.Ctl,
			"--port", fmt.Sprintf("%d", m.Port),
			"--device", fmt.Sprintf("%d", m.Device),
			"--type", string(m.DriveType),
			"--medium", "none",
		)
		if e != nil {
			return errors.New(e.Error() + ":" + eo)
		}
	}
	return nil
}

func (m *StorageMedium) Delete() error {
	e := m.Remove()
	if e != nil {
		fmt.Println(e.Error())
	}
	if m.Hdd != nil {
		return m.Hdd.Delete()
	} else {
		return errors.New("Hdd is nil")
	}
	return nil
}

// DriveType represents the hardware type of a drive.
type DriveType string

const (
	DriveDVD = DriveType("dvddrive")
	DriveHDD = DriveType("hdd")
	DriveFDD = DriveType("fdd")
)
