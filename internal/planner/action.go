package planner

import "Neolithic/internal/core"

// EntityType is the type of Entity that is being changed
type EntityType string

// AgentEntity is used to represent an agent
var AgentEntity EntityType = "agent"

// LocationEntity is used to represent a location
var LocationEntity EntityType = "location"

// StateChange represents a change to the state of the world
type StateChange struct {
	// Entity is the name of the entity that is being changed
	Entity string
	// EntityType is the type of entity that is being changed
	EntityType EntityType
	// Resource is the Resource that is being changed
	Resource *core.Resource
	// Amount is the Amount of the Resource that is being changed
	Amount int
}

// Action represents a thing that can be done
type Action interface {
	// Perform takes a start State and manipulates it into a new State, based on what the Action does.
	Perform(start *core.WorldState, agent core.Agent) *core.WorldState
	// Cost returns the ActionCost of the Action, representing how much work it takes to do the Action.
	Cost(agent core.Agent) float64
	// Description returns a string description of the Action, used to make the Action more legible.
	Description() string
	// GetChanges returns the difference in state before and after applying the Action
	GetChanges(agent core.Agent) []StateChange
}

// RequiresTime is an interface that provides a required Amount of time.
type RequiresTime interface {
	TimeNeeded() float64
}
