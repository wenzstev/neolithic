package planner

import "Neolithic/internal/core"

type EntityType string

var AgentEntity EntityType = "agent"
var LocationEntity EntityType = "location"

type StateChange struct {
	Entity     string
	EntityType EntityType
	Resource   *core.Resource
	Amount     int
}

// Action represents a thing that can be done
type Action interface {
	// Perform takes a start State and manipulates it into a new State, based on what the Action does.
	Perform(start *core.WorldState, agent core.Agent) *core.WorldState
	// Cost returns the cost of the Action, representing how much work it takes to do the Action.
	Cost(agent core.Agent) float64
	// Description returns a string description of the Action, used to make the Action more legible.
	Description() string
	// GetChanges returns the difference in state before and after applying the Action
	GetChanges(agent core.Agent) []StateChange
}
