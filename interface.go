package virtualbox

import "context"

// Virtualbox interface defines all the actions which can be performed by the
// Manager. This is mostly a utility interface designed for the customers of the
// package.
type Virtualbox interface {
	MachineManager
}

// MachineManager defines the actions that can be performed to manage machines
type MachineManager interface {
	// Machine gets a machine name based on its name or UUID
	Machine(context.Context, string) (*Machine, error)

	// ListMachines returns a list of all machines
	ListMachines(context.Context) ([]*Machine, error)

	// UpdateMachine allows to update the machine. The returned is the machine
	// in the current state after the update.
	UpdateMachine(context.Context, *Machine) (*Machine, error)

	// DeleteMachine deletes a machine by its name or UUID
	DeleteMachine(context.Context, string) error
}
