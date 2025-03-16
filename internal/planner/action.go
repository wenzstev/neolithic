package planner

import "Neolithic/internal/core"

// Action represents a thing that can be done
type Action interface {
	// Perform takes a start State and manipulates it into a new State, based on what the Action does.
	Perform(start *core.WorldState, agent core.Agent) *core.WorldState
	// Cost returns the cost of the Action, representing how much work it takes to do the Action.
	Cost(agent core.Agent) float64
	// Description returns a string description of the Action, used to make the Action more legible.
	Description() string
	// GetStateChange returns the difference in State before and after applying the Action
	GetStateChange(agent core.Agent) *core.WorldState
}
