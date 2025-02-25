package agent

import (
	"Neolithic/internal/planner"
)

// State represents an Agent's behavioral state.
type State interface {
	// Execute runs the state for a single unit of discrete time. It may produce a new world state, which indicates that
	// the agent has changed the world in some way. It may make changes to the agent, such as changing the agent's State,
	// goal, or plan.
	Execute(world WorldState) (WorldState, error)
}

// Behavior encapsulates the parts of the agent that are not in the physical WorldState
type Behavior struct {
	// PossibleActions represents all possible actions the agent can do. This is NOT the same as the actions in the
	// current plan
	PossibleActions *[]planner.Action
	// curPlan is the current plan the Agent is attempting to execute
	//nolint:unused
	curPlan Plan
	// Goal is the agent's desired WorldState
	Goal WorldState
	// curState is the current State the agent is in.
	//nolint:unused
	curState State
}

// RequiresTime is an interface that provides a required amount of time.
type RequiresTime interface { // TODO find better location
	TimeNeeded() float64
}
