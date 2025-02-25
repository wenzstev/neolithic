package agent

import "Neolithic/internal/planner"

// WorldState represents the state of the world in the Agent package
type WorldState interface {
	// Copy produces a deep copy of the WorldState
	// Right now this returns the concrete implementation, but we plan to change that on the next
	// PR. This is to avoid cyclic dependencies.
	Copy() *planner.State
}
