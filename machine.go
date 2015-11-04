package virtualbox

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type MachineState string

const (
	Poweroff = MachineState("poweroff")
	Running  = MachineState("running")
	Paused   = MachineState("paused")
	Saved    = MachineState("saved")
	Aborted  = MachineState("aborted")
)

type Flag int

// Flag names in lowercases to be consistent with VBoxManage options.
const (
	F_acpi Flag = 1 << iota
	F_ioapic
	F_rtcuseutc
	F_cpuhotplug
	F_pae
	F_longmode
	F_synthcpu
	F_hpet
	F_hwvirtex
	F_triplefaultreset
	F_nestedpaging
	F_largepages
	F_vtxvpid
	F_vtxux
	F_accelerate3d
)

// Convert bool to "on"/"off"
func bool2string(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

func string2bool(s string) bool {
	if s == "on" {
		return true
	}
	return false
}

// Test if flag is set. Return "on" or "off".
func (f Flag) Get(o Flag) string {
	return bool2string(f&o == o)
}

// Machine information.
type Machine struct {
	Name       string
	UUID       string
	State      MachineState
	CPUs       uint
	Memory     uint // main memory (in MB)
	VRAM       uint // video memory (in MB)
	CfgFile    string
	BaseFolder string
	OSType     string
	Flag       Flag
	BootOrder  []string // max 4 slots, each in {none|floppy|dvd|disk|net}
	NICs       map[int]*NIC
}

// Refresh reloads the machine information.
func (m *Machine) Refresh() error {
	id := m.Name
	if id == "" {
		id = m.UUID
	}
	mm, err := GetMachine(id)
	if err != nil {
		return err
	}
	*m = *mm
	return nil
}

// Start starts the machine.
func (m *Machine) Start() error {
	switch m.State {
	case Paused:
		return vbm("controlvm", m.Name, "resume")
	case Poweroff, Saved, Aborted:
		return vbm("startvm", m.Name, "--type", "headless")
	}
	return nil
}

// Suspend suspends the machine and saves its state to disk.
func (m *Machine) Save() error {
	switch m.State {
	case Paused:
		if err := m.Start(); err != nil {
			return err
		}
	case Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.Name, "savestate")
}

// Pause pauses the execution of the machine.
func (m *Machine) Pause() error {
	switch m.State {
	case Paused, Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.Name, "pause")
}

// Stop gracefully stops the machine.
func (m *Machine) Stop() error {
	switch m.State {
	case Poweroff, Aborted, Saved:
		return nil
	case Paused:
		if err := m.Start(); err != nil {
			return err
		}
	}

	for m.State != Poweroff { // busy wait until the machine is stopped
		//if err := vbm("controlvm", m.Name, "acpipowerbutton"); err != nil {
		if err := vbm("controlvm", m.Name, "poweroff"); err != nil {
			return err
		}
		time.Sleep(1 * time.Second)
		if err := m.Refresh(); err != nil {
			return err
		}
	}
	return nil
}

// Poweroff forcefully stops the machine. State is lost and might corrupt the disk image.
func (m *Machine) Poweroff() error {
	switch m.State {
	case Poweroff, Aborted, Saved:
		return nil
	}
	return vbm("controlvm", m.Name, "poweroff")
}

// Restart gracefully restarts the machine.
func (m *Machine) Restart() error {
	switch m.State {
	case Paused, Saved:
		if err := m.Start(); err != nil {
			return err
		}
	}
	if err := m.Stop(); err != nil {
		return err
	}
	return m.Start()
}

// Reset forcefully restarts the machine. State is lost and might corrupt the disk image.
func (m *Machine) Reset() error {
	switch m.State {
	case Paused, Saved:
		if err := m.Start(); err != nil {
			return err
		}
	}
	return vbm("controlvm", m.Name, "reset")
}

// Delete deletes the machine and associated disk images.
func (m *Machine) Delete() error {
	if err := m.Poweroff(); err != nil {
		return err
	}
	return vbm("unregistervm", m.Name, "--delete")
}

// GetMachine finds a machine by its name or UUID.
func GetMachine(id string) (*Machine, error) {
	stdout, stderr, err := vbmOutErr("showvminfo", id, "--machinereadable")
	if err != nil {
		if reMachineNotFound.FindString(stderr) != "" {
			return nil, ErrMachineNotExist
		}
		return nil, err
	}
	s := bufio.NewScanner(strings.NewReader(stdout))
	m := &Machine{}
	m.NICs = make(map[int]*NIC)

	for s.Scan() {
		res := reVMInfoLine.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		key := res[1]
		if key == "" {
			key = res[2]
		}
		val := res[3]
		if val == "" {
			val = res[4]
		}

		/*for k, v := range res {
			fmt.Printf("%d %s\n", k, v)
		}*/

		switch key {
		case "name":
			m.Name = val
		case "UUID":
			m.UUID = val
		case "VMState":
			m.State = MachineState(val)
		case "memory":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.Memory = uint(n)
		case "cpus":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.CPUs = uint(n)
		case "vram":
			n, err := strconv.ParseUint(val, 10, 32)
			if err != nil {
				return nil, err
			}
			m.VRAM = uint(n)
		case "CfgFile":
			m.CfgFile = val
			m.BaseFolder = filepath.Dir(val)

		case "nic1":
			if n, ok := m.NICs[1]; ok {
				n.Network = NICNetwork(val)
			} else {
				m.NICs[1] = &NIC{}
				m.NICs[1].Network = NICNetwork(val)
			}
		case "nictype1":
			if n, ok := m.NICs[1]; ok {
				n.Hardware = NICHardware(val)
			} else {
				m.NICs[1] = &NIC{}
				m.NICs[1].Hardware = NICHardware(val)
			}
		case "hostonlyadapter1", "nat-network1", "bridgeadapter1", "intnet1":
			if n, ok := m.NICs[1]; ok {
				n.InterfaceName = val
			} else {
				m.NICs[1] = &NIC{}
				m.NICs[1].InterfaceName = val
			}
		case "cableconnected1":
			if n, ok := m.NICs[1]; ok {
				n.CableConnected = string2bool(val)
			} else {
				m.NICs[1] = &NIC{}
				m.NICs[1].CableConnected = string2bool(val)
			}
		case "macaddress1":
			if n, ok := m.NICs[1]; ok {
				n.MACAddress = val
			} else {
				m.NICs[1] = &NIC{}
				m.NICs[1].MACAddress = val
			}

		case "nic2":
			if n, ok := m.NICs[2]; ok {
				n.Network = NICNetwork(val)
			} else {
				m.NICs[2] = &NIC{}
				m.NICs[2].Network = NICNetwork(val)
			}
		case "nictype2":
			if n, ok := m.NICs[2]; ok {
				n.Hardware = NICHardware(val)
			} else {
				m.NICs[2] = &NIC{}
				m.NICs[2].Hardware = NICHardware(val)
			}
		case "hostonlyadapter2", "nat-network2", "bridgeadapter2", "intnet2":
			if n, ok := m.NICs[2]; ok {
				n.InterfaceName = val
			} else {
				m.NICs[2] = &NIC{}
				m.NICs[2].InterfaceName = val
			}
		case "cableconnected2":
			if n, ok := m.NICs[2]; ok {
				n.CableConnected = string2bool(val)
			} else {
				m.NICs[2] = &NIC{}
				m.NICs[2].CableConnected = string2bool(val)
			}
		case "macaddress2":
			if n, ok := m.NICs[2]; ok {
				n.MACAddress = val
			} else {
				m.NICs[2] = &NIC{}
				m.NICs[2].MACAddress = val
			}
		}
	}
	if err := s.Err(); err != nil {
		return nil, err
	}

	return m, nil
}

func ImportMachine(filename string) error {
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		return errors.New("no such file or directory")
	}
	return vbm("import", filename)
}

// ListMachines lists all registered machines.
func ListMachines() ([]*Machine, error) {
	out, err := vbmOut("list", "vms")
	if err != nil {
		return nil, err
	}
	ms := []*Machine{}
	s := bufio.NewScanner(strings.NewReader(out))
	for s.Scan() {
		res := reVMNameUUID.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		if res[1] == "<inaccessible>" {
			continue
		}
		if len(res) != 3 {
			continue
		}

		m, err := GetMachine(res[2]) // res[1] is name, res[2] is uuid, now use uuid
		if err != nil {
			return nil, err
		}
		ms = append(ms, m)
	}
	if err := s.Err(); err != nil {
		return nil, err
	}
	return ms, nil
}

// CreateMachine creates a new machine. If basefolder is empty, use default.
func CreateMachine(name, basefolder string) (*Machine, error) {
	if name == "" {
		return nil, fmt.Errorf("machine name is empty")
	}

	// Check if a machine with the given name already exists.
	ms, err := ListMachines()
	if err != nil {
		return nil, err
	}
	for _, m := range ms {
		if m.Name == name {
			return nil, ErrMachineExist
		}
	}

	// Create and register the machine.
	args := []string{"createvm", "--name", name, "--register"}
	if basefolder != "" {
		args = append(args, "--basefolder", basefolder)
	}
	if err := vbm(args...); err != nil {
		return nil, err
	}

	m, err := GetMachine(name)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Modify changes the settings of the machine.
func (m *Machine) Modify() error {
	args := []string{"modifyvm", m.Name,
		"--firmware", "bios",
		"--bioslogofadein", "off",
		"--bioslogofadeout", "off",
		"--bioslogodisplaytime", "0",
		"--biosbootmenu", "disabled",

		"--ostype", m.OSType,
		"--cpus", fmt.Sprintf("%d", m.CPUs),
		"--memory", fmt.Sprintf("%d", m.Memory),
		"--vram", fmt.Sprintf("%d", m.VRAM),

		"--acpi", m.Flag.Get(F_acpi),
		"--ioapic", m.Flag.Get(F_ioapic),
		"--rtcuseutc", m.Flag.Get(F_rtcuseutc),
		"--cpuhotplug", m.Flag.Get(F_cpuhotplug),
		"--pae", m.Flag.Get(F_pae),
		"--longmode", m.Flag.Get(F_longmode),
		"--synthcpu", m.Flag.Get(F_synthcpu),
		"--hpet", m.Flag.Get(F_hpet),
		"--hwvirtex", m.Flag.Get(F_hwvirtex),
		"--triplefaultreset", m.Flag.Get(F_triplefaultreset),
		"--nestedpaging", m.Flag.Get(F_nestedpaging),
		"--largepages", m.Flag.Get(F_largepages),
		"--vtxvpid", m.Flag.Get(F_vtxvpid),
		"--vtxux", m.Flag.Get(F_vtxux),
		"--accelerate3d", m.Flag.Get(F_accelerate3d),
	}

	for i, dev := range m.BootOrder {
		if i > 3 {
			break // Only four slots `--boot{1,2,3,4}`. Ignore the rest.
		}
		args = append(args, fmt.Sprintf("--boot%d", i+1), dev)
	}
	if err := vbm(args...); err != nil {
		return err
	}
	return m.Refresh()
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *Machine) AddNATPF(n int, name string, rule PFRule) error {
	return vbm("controlvm", m.Name, fmt.Sprintf("natpf%d", n),
		fmt.Sprintf("%s,%s", name, rule.Format()))
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *Machine) DelNATPF(n int, name string) error {
	return vbm("controlvm", m.Name, fmt.Sprintf("natpf%d", n), "delete", name)
}

func (m *Machine) SetFrontend(front string) error {
	return vbm("modifyvm", m.UUID, "--defaultfrontend", front)
}

// SetNIC set the n-th NIC.
func (m *Machine) SetNIC(n int, nic NIC) error {
	if m.State == Running || m.State == Saved || m.State == Paused {
		return errors.New("Session has been locked")
	}

	args := []string{"modifyvm", m.Name,
		fmt.Sprintf("--nic%d", n), string(nic.Network),
		fmt.Sprintf("--nictype%d", n), string(nic.Hardware),
		fmt.Sprintf("--cableconnected%d", n), "on",
	}

	if nic.Network == "hostonly" {
		args = append(args, fmt.Sprintf("--hostonlyadapter%d", n), nic.InterfaceName)
	}
	return vbm(args...)
}

// SetNIC set the n-th NIC.
func (m *Machine) SetVRdpPort(port int) error {
	if m.State == Running || m.State == Saved || m.State == Paused {
		return errors.New("Session has been locked")
	}

	if port > 0 {
		cmd := []string{"modifyvm", m.Name,
			"--vrde", "on"}
		err := vbm(cmd...)
		if err != nil {
			return err
		}

		cmd = []string{"modifyvm", m.Name,
			"--vrdeport", strconv.Itoa(port),
		}
		return vbm(cmd...)
	}

	return nil
}

// AddStorageCtl adds a storage controller with the given name.
func (m *Machine) AddStorageCtl(name string, ctl StorageController) error {
	args := []string{"storagectl", m.Name, "--name", name}
	if ctl.SysBus != "" {
		args = append(args, "--add", string(ctl.SysBus))
	}
	if ctl.Ports > 0 {
		args = append(args, "--portcount", fmt.Sprintf("%d", ctl.Ports))
	}
	if ctl.Chipset != "" {
		args = append(args, "--controller", string(ctl.Chipset))
	}
	args = append(args, "--hostiocache", bool2string(ctl.HostIOCache))
	args = append(args, "--bootable", bool2string(ctl.Bootable))
	return vbm(args...)
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *Machine) DelStorageCtl(name string) error {
	return vbm("storagectl", m.Name, "--name", name, "--remove")
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *Machine) AttachStorage(ctlName string, medium StorageMedium) error {
	return vbm("storageattach", m.Name, "--storagectl", ctlName,
		"--port", fmt.Sprintf("%d", medium.Port),
		"--device", fmt.Sprintf("%d", medium.Device),
		"--type", string(medium.DriveType),
		"--medium", medium.Medium,
	)
}

// Property Set
func (m *Machine) PropertySet(key, val string) error {
	return guestPropertySet(m.UUID, key, val)
}

// Property Get
func (m *Machine) PropertyGet(key string) (string, error) {
	return guestPropertyGet(m.UUID, key)
}

// Property Get
func (m *Machine) PropertyDel(key string) error {
	return guestPropertyDel(m.UUID, key)
}

// Property Wait
func (m *Machine) PropertyWait(key string, timeout int) (string, error) {
	return guestPropertyWait(m.UUID, key, timeout)
}

// Property Enumerate
func (m *Machine) PropertyEnumerate() (map[string]string, error) {
	return guestPropertyEnumerate(m.UUID)
}

// MAC Address Set
func (m *Machine) MACAddressSet(solt int, mac string) error {
	return modifyMacAddress(m.UUID, solt, mac)
}

func (m *Machine) Clone(name string) error {
	if m.State != Poweroff {
		return errors.New("Machine is not poweroff")
	}

	_, errstr, err := vbmOutErr("clonevm", m.UUID,
		"--name", name,
		"--mode", "all",
		"--register",
	)
	if err != nil {
		return errors.New(errstr)
	}
	return nil
}
