package core

import "fmt"

// Agent is the core interface for all packages that require an agent
type Agent interface {
	fmt.Stringer
	// Name returns the name of the agent
	Name() string
	// DeepCopy returns a deep copy of the agent
	DeepCopy() Agent
	// Inventory returns the agent's current inventory
	Inventory() Inventory
}
