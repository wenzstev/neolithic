package agent

import (
	"Neolithic/internal/goalengine"
	"encoding/gob"

	"Neolithic/internal/core"
	"Neolithic/internal/planner"
)

func init() {
	gob.Register(Idle{})
	gob.Register(Moving{})
	gob.Register(Performing{})
}

// State represents an Agent's behavioral state.
type State interface {
	// Execute runs the state for a single unit of discrete time. It may produce a new world state, which indicates that
	// the Agent has changed the world in some way. It may make changes to the Agent, such as changing the Agent's State,
	// goal, or plan.
	Execute(world *core.WorldState, deltaTime float64) (*core.WorldState, error)
}

// Behavior encapsulates the parts of the Agent that are not in the physical WorldState
type Behavior struct {
	// PossibleActions represents all possible actions the Agent can do. This is NOT the same as the actions in the
	// current plan
	PossibleActions []planner.Action
	// CurPlan is the current plan the Agent is attempting to execute
	//nolint:unused
	CurPlan Plan
	// Goal is the Agent's desired WorldState
	Goal *core.WorldState
	// CurState is the current State the Agent is in.
	CurState State
	// GoalEngine is used to determine the agent's current and future goals
	GoalEngine *goalengine.GoalEngine
}
