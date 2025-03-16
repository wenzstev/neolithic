package agent

import "Neolithic/internal/core"

// agent represents an agent in the world
type Agent struct {
	// name is the name of the agent
	name string
	// behavior holds the agent's decision-making processes
	Behavior  *Behavior
	Inventory core.Inventory
}

var _ core.Agent = (*Agent)(nil)

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) AdjustInventory(resource *core.Resource, i int) core.Agent {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) GetAmount(resource *core.Resource) int {
	//TODO implement me
	panic("implement me")
}

func (a *Agent) Copy() *Agent {
	panic("implement me")
}
