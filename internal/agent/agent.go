package agent

type Agent interface {
	Name() string
	SetCurState(state State)
	Plan() Plan
}

// agent represents an agent in the world
type agent struct {
	// name is the name of the agent
	name string
	// behavior holds the agent's decision-making processes
	behavior *Behavior
}

func (a *agent) Plan() Plan {
	return a.behavior.curPlan
}

// Name provides the name of the agent
func (a *agent) Name() string {
	return a.name
}

func (a *agent) SetCurState(state State) {
	a.behavior.curState = state
}
