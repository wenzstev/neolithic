package agent

// Agent represents an agent in the world
type Agent struct {
	// name is the name of the agent
	name string
	// todo add behavior and other necessary values
}

// Name provides the name of the agent
func (a *Agent) Name() string {
	return a.name
}
