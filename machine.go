package virtualbox

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// Machine returns the information about existing virtualbox machine identified
// by either its UUID or name.
func (m *Manager) Machine(ctx context.Context, id string) (*Machine, error) {
	m.log.Printf("getting information for %q", id)
	// There is a strage behavior where running multiple instances of
	// 'VBoxManage showvminfo' on same VM simultaneously can return an error of
	// 'object is not ready (E_ACCESSDENIED)', so we sequential the operation with a mutex.
	// Note if you are running multiple process of go-virtualbox or 'showvminfo'
	// in the command line side by side, this not gonna work.
	// TODO: Verify the above is still true.
	m.lock.Lock()
	stdout, stderr, err := m.run(ctx, "showvminfo", id, "--machinereadable")
	m.lock.Unlock()
	if err != nil {
		if reMachineNotFound.FindString(stderr) != "" {
			return nil, ErrMachineNotExist
		}
		return nil, err
	}

	/* Read all VM info into a map */
	props := make(map[string]string)
	s := bufio.NewScanner(strings.NewReader(stdout))
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
		props[key] = val
	}

	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("unable to scan all fields: %w", err)
	}

	// error that occured during parsing
	var perr error

	sp := func(field string, def ...string) string {
		if v, exists := props[field]; exists {
			return v
		}
		if len(def) < 1 {
			return ""
		}
		return def[0]
	}

	up := func(field string, def ...uint) uint {
		if v, exists := props[field]; exists {
			n, err := strconv.ParseUint(v, 10, 32)
			if err != nil {
				perr = err
				return 0
			}
			return uint(n)
		}
		if len(def) < 1 {
			return 0
		}
		return def[0]
	}

	/* Extract basic info */
	vm := &Machine{
		// TODO: This was in New, verify is this still correct.
		BootOrder:  make([]string, 0, 4),
		NICs:       make([]NIC, 0, 4),
		Name:       sp("name"),
		Firmware:   sp("firmware"),
		UUID:       sp("UUID"),
		State:      MachineState(sp("VMState")),
		Memory:     up("memory"),
		CPUs:       up("cpus"),
		VRAM:       up("vram"),
		CfgFile:    sp("CfgFile"),
		BaseFolder: filepath.Dir(sp("CfgFile")),
	}

	/* Extract NIC info */
	for i := 1; i <= 4; i++ {
		var nic NIC
		nicType, ok := props[fmt.Sprintf("nic%d", i)]
		if !ok || nicType == "none" {
			break
		}
		nic.Network = NICNetwork(nicType)
		nic.Hardware = NICHardware(props[fmt.Sprintf("nictype%d", i)])
		if nic.Hardware == "" {
			return nil, fmt.Errorf("Could not find corresponding 'nictype%d'", i)
		}
		nic.MacAddr = props[fmt.Sprintf("macaddress%d", i)]
		if nic.MacAddr == "" {
			return nil, fmt.Errorf("Could not find corresponding 'macaddress%d'", i)
		}
		if nic.Network == NICNetHostonly {
			nic.HostInterface = props[fmt.Sprintf("hostonlyadapter%d", i)]
		} else if nic.Network == NICNetBridged {
			nic.HostInterface = props[fmt.Sprintf("bridgeadapter%d", i)]
		}
		vm.NICs = append(vm.NICs, nic)
	}

	if perr != nil {
		return nil, fmt.Errorf("parsing machine props failed: %w", perr)
	}

	return vm, nil
}

// ListMachines returns the list of the machines
func (m *Manager) ListMachines(ctx context.Context) ([]*Machine, error) {
	m.log.Println("listing vms")
	stdout, _, err := m.run(ctx, "list", "vms")
	if err != nil {
		return nil, fmt.Errorf("unable to list vms: %w", err)
	}
	vms := []*Machine{}
	s := bufio.NewScanner(strings.NewReader(stdout))
	for s.Scan() {
		res := reVMNameUUID.FindStringSubmatch(s.Text())
		if res == nil {
			continue
		}
		m, err := m.Machine(ctx, res[1])
		if err != nil {
			// Sometimes a VM is listed but not available, so we need to handle this.
			if errors.Is(err, ErrMachineNotExist) {
				continue
			} else {
				return nil, fmt.Errorf("unable to get machine info: %w", err)
			}
		}
		vms = append(vms, m)
	}
	if err := s.Err(); err != nil {
		return nil, fmt.Errorf("error reading machine list: %w", err)
	}
	return vms, nil
}

// ModifyMachine modifies the data of the machine
func (m *Manager) ModifyMachine(ctx context.Context, vm *Machine) error {
	args := []string{"modifyvm", vm.Name,
		"--firmware", vm.Firmware,
		"--bioslogofadein", "off",
		"--bioslogofadeout", "off",
		"--bioslogodisplaytime", "0",
		"--biosbootmenu", "disabled",

		"--ostype", vm.OSType,
		"--cpus", fmt.Sprintf("%d", vm.CPUs),
		"--memory", fmt.Sprintf("%d", vm.Memory),
		"--vram", fmt.Sprintf("%d", vm.VRAM),

		"--acpi", vm.Flag.Get(ACPI),
		"--ioapic", vm.Flag.Get(IOAPIC),
		"--rtcuseutc", vm.Flag.Get(RTCUSEUTC),
		"--cpuhotplug", vm.Flag.Get(CPUHOTPLUG),
		"--pae", vm.Flag.Get(PAE),
		"--longmode", vm.Flag.Get(LONGMODE),
		"--hpet", vm.Flag.Get(HPET),
		"--hwvirtex", vm.Flag.Get(HWVIRTEX),
		"--triplefaultreset", vm.Flag.Get(TRIPLEFAULTRESET),
		"--nestedpaging", vm.Flag.Get(NESTEDPAGING),
		"--largepages", vm.Flag.Get(LARGEPAGES),
		"--vtxvpid", vm.Flag.Get(VTXVPID),
		"--vtxux", vm.Flag.Get(VTXUX),
		"--accelerate3d", vm.Flag.Get(ACCELERATE3D),
	}

	for i, dev := range vm.BootOrder {
		if i > 3 {
			break // Only four slots `--boot{1,2,3,4}`. Ignore the rest.
		}
		args = append(args, fmt.Sprintf("--boot%d", i+1), dev)
	}

	for i, nic := range vm.NICs {
		n := i + 1
		args = append(args,
			fmt.Sprintf("--nic%d", n), string(nic.Network),
			fmt.Sprintf("--nictype%d", n), string(nic.Hardware),
			fmt.Sprintf("--cableconnected%d", n), "on")
		if nic.Network == NICNetHostonly {
			args = append(args, fmt.Sprintf("--hostonlyadapter%d", n), nic.HostInterface)
		} else if nic.Network == NICNetBridged {
			args = append(args, fmt.Sprintf("--bridgeadapter%d", n), nic.HostInterface)
		}
		if nic.MacAddr != "" {
			args = append(args, fmt.Sprintf("--macaddress%d", n), nic.MacAddr)
		}
	}

	if _, _, err := m.run(ctx, args...); err != nil {
		return err
	}
	return vm.Refresh()
}

// StartMachine will start the machine based on its current state.
func (m *Manager) StartMachine(ctx context.Context, id string) error {
	var args []string

	vm, err := m.Machine(ctx, id)
	if err != nil {
		return fmt.Errorf("unable to get machine to check its status: %w", err)
	}

	switch vm.State {
	case Paused:
		args = []string{"controlvm", id, "resume"}
	case Poweroff, Saved, Aborted:
		args = []string{"startvm", id, "--type", "headless"}
	}

	_, msg, err := m.run(ctx, args...)
	if err != nil {
		return errors.New(msg)
	}

	return nil
}

// MachineState stores the last retrieved VM state.
type MachineState string

const (
	// Poweroff is a MachineState value.
	Poweroff = MachineState("poweroff")
	// Running is a MachineState value.
	Running = MachineState("running")
	// Paused is a MachineState value.
	Paused = MachineState("paused")
	// Saved is a MachineState value.
	Saved = MachineState("saved")
	// Aborted is a MachineState value.
	Aborted = MachineState("aborted")
)

// Flag is an active VM configuration toggle
type Flag int

// Flag names in lowercases to be consistent with VBoxManage options.
const (
	ACPI Flag = 1 << iota
	IOAPIC
	RTCUSEUTC
	CPUHOTPLUG
	PAE
	LONGMODE
	HPET
	HWVIRTEX
	TRIPLEFAULTRESET
	NESTEDPAGING
	LARGEPAGES
	VTXVPID
	VTXUX
	ACCELERATE3D
)

// Convert bool to "on"/"off"
func bool2string(b bool) string {
	if b {
		return "on"
	}
	return "off"
}

// Get tests if flag is set. Return "on" or "off".
func (f Flag) Get(o Flag) string {
	return bool2string(f&o == o)
}

// Machine information.
type Machine struct {
	Name       string
	Firmware   string
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
	NICs       []NIC
}

// New creates a new machine.
func New() *Machine {
	return &Machine{
		BootOrder: make([]string, 0, 4),
		NICs:      make([]NIC, 0, 4),
	}
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

// Start the machine, and return the underlying error when unable to do so.
func (m *Machine) Start() error {
	return defaultManager.StartMachine(context.Background(), m.UUID)
}

// DisconnectSerialPort sets given serial port to disconnected.
func (m *Machine) DisconnectSerialPort(portNumber int) error {
	_, _, err := Manage().run("modifyvm", m.Name, fmt.Sprintf("--uartmode%d", portNumber), "disconnected")
	return err
}

// Save suspends the machine and saves its state to disk.
func (m *Machine) Save() error {
	switch m.State {
	case Paused:
		if err := m.Start(); err != nil {
			return err
		}
	case Poweroff, Aborted, Saved:
		return nil
	}
	_, _, err := Manage().run("controlvm", m.Name, "savestate")
	return err
}

// Pause pauses the execution of the machine.
func (m *Machine) Pause() error {
	switch m.State {
	case Paused, Poweroff, Aborted, Saved:
		return nil
	}
	_, _, err := Manage().run("controlvm", m.Name, "pause")
	return err
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
		if _, _, err := Manage().run("controlvm", m.Name, "acpipowerbutton"); err != nil {
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
	_, _, err := Manage().run("controlvm", m.Name, "poweroff")
	return err
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
	_, _, err := Manage().run("controlvm", m.Name, "reset")
	return err
}

// Delete deletes the machine and associated disk images.
func (m *Machine) Delete() error {
	if err := m.Poweroff(); err != nil {
		return err
	}
	_, _, err := Manage().run("unregistervm", m.Name, "--delete")
	return err
}

// GetMachine finds a machine by its name or UUID.
// DEPRECATED: Use (*Manager).Machine
func GetMachine(id string) (*Machine, error) {
	return defaultManager.Machine(context.Background(), id)
}

// ListMachines lists all registered machines.
// DEPRECATED: Use (*Manager).ListMachines
func ListMachines() ([]*Machine, error) {
	return defaultManager.ListMachines(context.Background())
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
	if _, _, err = Manage().run(args...); err != nil {
		return nil, err
	}

	m, err := GetMachine(name)
	if err != nil {
		return nil, err
	}

	return m, nil
}

// Modify changes the settings of the machine.
// DEPRECATED: Use (*Manager).ModifyMachine
func (m *Machine) Modify() error {
	return defaultManager.ModifyMachine(context.Background(), m)
}

// AddNATPF adds a NAT port forarding rule to the n-th NIC with the given name.
func (m *Machine) AddNATPF(n int, name string, rule PFRule) error {
	_, _, err := Manage().run("controlvm", m.Name, fmt.Sprintf("natpf%d", n),
		fmt.Sprintf("%s,%s", name, rule.Format()))
	return err
}

// DelNATPF deletes the NAT port forwarding rule with the given name from the n-th NIC.
func (m *Machine) DelNATPF(n int, name string) error {
	_, _, err := Manage().run("controlvm", m.Name, fmt.Sprintf("natpf%d", n), "delete", name)
	return err
}

// SetNIC set the n-th NIC.
func (m *Machine) SetNIC(n int, nic NIC) error {
	args := []string{"modifyvm", m.Name,
		fmt.Sprintf("--nic%d", n), string(nic.Network),
		fmt.Sprintf("--nictype%d", n), string(nic.Hardware),
		fmt.Sprintf("--cableconnected%d", n), "on",
	}

	if nic.Network == NICNetHostonly {
		args = append(args, fmt.Sprintf("--hostonlyadapter%d", n), nic.HostInterface)
	} else if nic.Network == NICNetBridged {
		args = append(args, fmt.Sprintf("--bridgeadapter%d", n), nic.HostInterface)
	}
	_, _, err := Manage().run(args...)
	return err
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

	_, _, err := Manage().run(args...)
	return err
}

// DelStorageCtl deletes the storage controller with the given name.
func (m *Machine) DelStorageCtl(name string) error {
	_, _, err := Manage().run("storagectl", m.Name, "--name", name, "--remove")
	return err
}

// AttachStorage attaches a storage medium to the named storage controller.
func (m *Machine) AttachStorage(ctlName string, medium StorageMedium) error {
	_, _, err := Manage().run("storageattach", m.Name, "--storagectl", ctlName,
		"--port", fmt.Sprintf("%d", medium.Port),
		"--device", fmt.Sprintf("%d", medium.Device),
		"--type", string(medium.DriveType),
		"--medium", medium.Medium,
	)
	return err
}

// SetExtraData attaches custom string to the VM.
func (m *Machine) SetExtraData(key, val string) error {
	_, _, err := Manage().run("setextradata", m.Name, key, val)
	return err
}

// GetExtraData retrieves custom string from the VM.
func (m *Machine) GetExtraData(key string) (*string, error) {
	value, _, err := Manage().run("getextradata", m.Name, key)
	if err != nil {
		return nil, err
	}
	value = strings.TrimSpace(value)
	/* 'getextradata get' returns 0 even when the key is not found,
	so we need to check stdout for this case */
	if strings.HasPrefix(value, "No value set") {
		return nil, nil
	}
	trimmed := strings.TrimPrefix(value, "Value: ")
	return &trimmed, nil
}

// DeleteExtraData removes custom string from the VM.
func (m *Machine) DeleteExtraData(key string) error {
	_, _, err := Manage().run("setextradata", m.Name, key)
	return err
}

// CloneMachine clones the given machine name into a new one.
func CloneMachine(baseImageName string, newImageName string, register bool) error {
	if register {
		_, _, err := Manage().run("clonevm", baseImageName, "--name", newImageName, "--register")
		return err
	}
	_, _, err := Manage().run("clonevm", baseImageName, "--name", newImageName)
	return err
}
