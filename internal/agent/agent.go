package agent

// Agent is the public interface for interacting with the agent.
type Agent interface {
	// Name provides the name of the agent
	Name() string
	// Behavior provides the behavior struct for the agent.
	Behavior() *Behavior
}

// agent represents an agent in the world
type agent struct {
	// name is the name of the agent
	name string
	// behavior holds the agent's decision-making processes
	behavior *Behavior
}

// Name implements Agent.Name and returns the name of the agent
func (a *agent) Name() string {
	return a.name
}

// Behavior implements Agent.Behavior and returns the Behavior of the agent
func (a *agent) Behavior() *Behavior {
	return a.behavior
}
