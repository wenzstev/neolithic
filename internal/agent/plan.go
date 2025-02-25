package agent

import "Neolithic/internal/planner"

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

// Ensure plan implements the Plan interface
var _ Plan = (*plan)(nil)

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
	if p.IsComplete() {
		return nil
	}
	action := p.PeekAction()
	p.curLocation++
	return action
}
