package agent

import (
	"Neolithic/internal/planner"
	"errors"
)

type Performing struct {
	action   planner.Action
	timeLeft float64
	agent    *Agent
}

// Execute implements State.Execute and simulates the performance of an action. If the Action takes a period of time,
// the Agent will stay in the Execute state until the necessary amount of time has elapsed. Afterward (of if the Action
// is instant), the agent will call planner.Action.Perform, and return the result. It will also change the Agent's State
// to either Idle or Moving, depending on if the plan is complete.
func (p *Performing) Execute(world WorldState) (WorldState, error) {
	if p.action == nil {
		p.action = p.agent.behavior.curPlan.PeekAction()
		if p.action == nil { // still nil, plan complete
			p.agent.behavior.curState = &Idle{}
			return (*planner.State)(nil), nil
		}
	}

	actionDuration, ok := p.action.(RequiresTime)
	if ok {
		p.timeLeft = actionDuration.TimeNeeded()
	}

	if p.timeLeft > 0 {
		p.timeLeft -= 1.0 / 60.0 // called every tick, update is called 60 times p second
		return (*planner.State)(nil), nil
	}

	// for now, cast the state into planner state
	worldState, ok := world.(*planner.State)
	if !ok {
		return nil, errors.New("world state was not a planner.State")
	}

	newState := p.action.Perform(worldState, p.agent)
	if newState == nil { // action failed
		p.agent.behavior.curState = &Idle{}
		return (*planner.State)(nil), nil
	}
	p.agent.behavior.curPlan.PopAction()
	if p.agent.behavior.curPlan.IsComplete() {
		p.agent.behavior.curState = &Idle{}
	} else {
		p.agent.behavior.curState = &Moving{}
	}

	return newState, nil
}
