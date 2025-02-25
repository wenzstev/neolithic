package agent

import (
	"Neolithic/internal/planner"
)

type Plan interface {
	IsComplete() bool
	PeekAction() planner.Action
	PopAction() planner.Action
}

// plan represents the agent's current plan, as created by the GOAP system
type plan struct {
	// Actions are the actions that make up the plan.
	Actions *[]planner.Action
	// curLocation is used to determine the current step in the plan.
	curLocation int
}

// IsComplete indicates if a plan has completed all steps
func (p *plan) IsComplete() bool {
	return p.curLocation >= len(*p.Actions)
}

// PeekAction provides the next action in the plan. It does _not_ pop the action.
func (p *plan) PeekAction() planner.Action {
	if p.IsComplete() {
		return nil
	}
	return (*p.Actions)[p.curLocation]
}

// PopAction returns the current action and increments the counter.
func (p *plan) PopAction() planner.Action {
	if !p.IsComplete() {
		return nil
	}
	action := p.PeekAction()
	p.curLocation++
	return action
}

// Behavior encapsulates the parts of the agent that are not in the physical WorldState
type Behavior struct {
	// PossibleActions represents all possible actions the agent can do. This is NOT the same as the actions in the
	// current plan
	PossibleActions *[]planner.Action
	// curPlan is the current plan the Agent is attempting to execute
	curPlan Plan
	// Goal is the agent's desired WorldState
	Goal WorldState
	// curState is the current State the agent is in.
	curState State
}

// RequiresTime is an interface that provides a required amount of time.
type RequiresTime interface { // TODO find better location
	TimeNeeded() float64
}

// State represents an Agent's behavioral state.
type State interface {
	// Execute runs the state for a single unit of discrete time. It may produce a new world state, which indicates that
	// the agent has changed the world in some way. It may make changes to the agent, such as changing the agent's State,
	// goal, or plan.
	Execute(world WorldState) (WorldState, error)
}
