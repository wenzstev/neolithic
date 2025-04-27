package agent

import (
	"errors"
	"log/slog"

	"Neolithic/internal/core"
)

// Performing is the State where the Agent tries to perform its provided Action
type Performing struct {
	// action is the action that will be performed
	action core.Action
	// timeLeft is the amount of time before the action is completed, if necessary
	timeLeft float64
	// agent is the agent that is performing the action
	agent *Agent
	// logger is the logger
	logger *slog.Logger
}

// Execute implements State.Execute and simulates the performance of an action. If the Action takes a period of time,
// the Agent will stay in the Execute state until the necessary amount of time has elapsed. Afterward (of if the Action
// is instant), the Agent will call planner.Action.Perform, and return the result. It will also change the Agent's State
// to either Idle or Moving, depending on if the plan is complete.
func (p *Performing) Execute(world *core.WorldState, deltaTime float64) (*core.WorldState, error) {
	p.logger.Debug("performing state execute", "agent", p.agent.Name(), "deltaTime", deltaTime)

	behavior := p.agent.Behavior
	curPlan := behavior.CurPlan

	// get the next action and determine if time is needed
	if p.action == nil {
		p.action = curPlan.PeekAction()
		if p.action == nil { // still nil, plan complete
			p.logger.Info("plan complete, transitioning to idle", "agent", p.agent.Name())
			behavior.CurState = &Idle{agent: p.agent, logger: p.logger}
			return (*core.WorldState)(nil), nil
		}

		p.logger.Debug("starting new action", "agent", p.agent.Name(), "action", p.action)
		actionDuration, ok := p.action.(core.RequiresTime)
		if ok && p.timeLeft == 0 {
			p.timeLeft = actionDuration.TimeNeeded()
			p.logger.Debug("action requires time", "agent", p.agent.Name(), "timeNeeded", p.timeLeft)
		}
	}

	// if there's time on the clock, increment by delta time and return
	if p.timeLeft > 0 {
		p.timeLeft -= deltaTime // called every tick, update is called 60 times p second
		p.logger.Debug("action in progress", "agent", p.agent.Name(), "timeLeft", p.timeLeft)
		return (*core.WorldState)(nil), nil
	}

	p.logger.Info("performing action", "agent", p.agent.Name(), "action", p.action)
	newWorldState := p.action.Perform(world, p.agent)
	if newWorldState == nil { // action failed
		p.logger.Error("action failed", "agent", p.agent.Name(), "action", p.action)
		behavior.CurState = &Idle{agent: p.agent, logger: p.logger}
		return (*core.WorldState)(nil), nil
	}

	newAgentInterface, exists := newWorldState.GetAgent(p.agent.Name())
	if !exists {
		p.logger.Error("agent does not exist in deep copied world", "agent", p.agent.Name())
		return nil, errors.New("agent does not exist in deep copied world")
	}
	newAgent := newAgentInterface.(*Agent)
	behavior = newAgent.Behavior
	curPlan = behavior.CurPlan

	curPlan.PopAction()
	if curPlan.IsComplete() {
		p.logger.Info("plan complete after action, transitioning to idle", "agent", p.agent.Name())
		behavior.CurState = &Idle{agent: newAgent, logger: p.logger}
	} else {
		p.logger.Info("action complete, transitioning to moving", "agent", p.agent.Name())
		behavior.CurState = &Moving{agent: newAgent, logger: p.logger}
	}

	return newWorldState, nil
}

// NewPerforming creates a new Performing state
func NewPerforming(agent *Agent, logger *slog.Logger) *Performing {
	return &Performing{
		agent:  agent,
		logger: logger,
	}
}
