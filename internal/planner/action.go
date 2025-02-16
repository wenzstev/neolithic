package planner

// Action represents a thing that can be done
type Action interface {
	// Perform takes a start state and manipulates it into a new state, based on what the Action does.
	Perform(start *State, agent *Agent) *State
	// Cost returns the cost of the Action, representing how much work it takes to do the Action.
	Cost(agent *Agent) float64
	// Description returns a string description of the action, used to make the action more legible.
	Description() string
	// GetStateChange returns the difference in state before and after applying the Action
	GetStateChange(agent *Agent) *State
}
