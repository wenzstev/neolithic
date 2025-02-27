package agent

import (
	"Neolithic/internal/planner"
	"errors"
)

type Performing struct {
	action   planner.Action
	timeLeft float64
	agent    Agent
}

// Execute implements State.Execute and simulates the performance of an action. If the Action takes a period of time,
// the Agent will stay in the Execute state until the necessary amount of time has elapsed. Afterward (of if the Action
// is instant), the agent will call planner.Action.Perform, and return the result. It will also change the Agent's State
// to either Idle or Moving, depending on if the plan is complete.
func (p *Performing) Execute(world WorldState, deltaTime float64) (WorldState, error) {
	curPlan := p.agent.Plan()

	if p.action == nil {
		p.action = curPlan.PeekAction()
		if p.action == nil { // still nil, plan complete
			p.agent.SetCurState(&Idle{})
			return (*planner.State)(nil), nil
		}
	}

	actionDuration, ok := p.action.(RequiresTime)
	if ok && p.timeLeft == 0 {
		p.timeLeft = actionDuration.TimeNeeded()
	}

	if p.timeLeft > 0 {
		p.timeLeft -= deltaTime // called every tick, update is called 60 times p second
		return (*planner.State)(nil), nil
	}

	// for now, cast the state into planner state
	worldState, ok := world.(*planner.State)
	if !ok {
		return nil, errors.New("world state was not a planner.State")
	}

	newState := p.action.Perform(worldState, p.agent)
	if newState == nil { // action failed
		p.agent.SetCurState(&Idle{})
		return (*planner.State)(nil), nil
	}
	curPlan.PopAction()
	if curPlan.IsComplete() {
		p.agent.SetCurState(&Idle{})
	} else {
		p.agent.SetCurState(&Moving{})
	}

	return newState, nil
}
