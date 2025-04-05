package agent

import (
	"Neolithic/internal/planner"
	"encoding/gob"
)

func init() {
	gob.Register(plan{})
}

// Plan provides a series of actions for an Agent to complete
type Plan interface {
	// IsComplete returns if the plan has been completed (all actions have been done)
	IsComplete() bool
	// PeekAction looks at the next action on the plan WITHOUT marking it done
	PeekAction() planner.Action
	// PopAction removes the next action on the plan
	PopAction() planner.Action
}

// plan represents the Agent's current plan, as created by the GOAP system
type plan struct {
	// Actions are the actions that make up the plan.
	Actions []planner.Action
	// curLocation is used to determine the current step in the plan.
	curLocation int
}

// Ensure plan implements the Plan interface
var _ Plan = (*plan)(nil)

// IsComplete implements Plan.IsComplete. It indicates if a
// plan has completed all steps.
func (p *plan) IsComplete() bool {
	return p.curLocation >= len(p.Actions)
}

// PeekAction implements Plan.PeekAction. It provides the next action
// in the plan. It does _not_ pop the action.
func (p *plan) PeekAction() planner.Action {
	if p.IsComplete() {
		return nil
	}
	return (p.Actions)[p.curLocation]
}

// PopAction implements Plan.PopAction. It returns the current action
// and increments the counter.
func (p *plan) PopAction() planner.Action {
	if p.IsComplete() {
		return nil
	}
	action := p.PeekAction()
	p.curLocation++
	return action
}
