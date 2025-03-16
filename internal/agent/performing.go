package agent

import (
	"Neolithic/internal/core"
	"Neolithic/internal/planner"
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
func (p *Performing) Execute(world *core.WorldState, deltaTime float64) (*core.WorldState, error) {
	behavior := p.agent.Behavior
	curPlan := behavior.CurPlan

	// get the next action and determine if time is needed
	if p.action == nil {
		p.action = curPlan.PeekAction()
		if p.action == nil { // still nil, plan complete
			behavior.curState = &Idle{}
			return (*core.WorldState)(nil), nil
		}

		actionDuration, ok := p.action.(RequiresTime)
		if ok && p.timeLeft == 0 {
			p.timeLeft = actionDuration.TimeNeeded()
		}
	}

	// if there's time on the clock, increment by delta time and return
	if p.timeLeft > 0 {
		p.timeLeft -= deltaTime // called every tick, update is called 60 times p second
		return (*core.WorldState)(nil), nil
	}

	// perform the action
	newWorldState := p.action.Perform(world, p.agent)
	if newWorldState == nil { // action failed
		behavior.curState = &Idle{}
		return (*core.WorldState)(nil), nil
	}

	curPlan.PopAction()
	if curPlan.IsComplete() {
		behavior.curState = &Idle{}
	} else {
		behavior.curState = &Moving{}
	}

	return newWorldState, nil
}
