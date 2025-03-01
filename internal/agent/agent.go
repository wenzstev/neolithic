package agent

// Agent is the public interface for interacting with the agent.
type Agent interface {
	// Name provides the name of the agent
	Name() string
	// SetCurState allows for setting the Agent's current state
	SetCurState(state State)
	// Plan provides the Agent's current plan
	Plan() Plan
}

// agent represents an agent in the world
type agent struct {
	// name is the name of the agent
	name string
	// behavior holds the agent's decision-making processes
	behavior *Behavior
}

// Plan provides the agent's current Plan
func (a *agent) Plan() Plan {
	return a.behavior.curPlan
}

// Name provides the name of the agent.
func (a *agent) Name() string {
	return a.name
}

// SetCurState sets the state of the agent.
func (a *agent) SetCurState(state State) {
	a.behavior.curState = state
}
