package agent

import (
	"Neolithic/internal/core"
	"fmt"
)

// Agent represents an agent in the world
type Agent struct {
	// name is the name of the agent
	name string
	// behavior holds the agent's decision-making processes
	Behavior  *Behavior
	inventory core.Inventory
	Position  core.Coord
}

var _ core.Agent = (*Agent)(nil)

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) Inventory() core.Inventory {
	return a.inventory
}

func (a *Agent) DeepCopy() core.Agent {
	newAgent := &Agent{}
	newAgent.name = a.name
	if a.Behavior != nil {
		newAgent.Behavior = &Behavior{
			PossibleActions: a.Behavior.PossibleActions,
			CurPlan:         a.Behavior.CurPlan,
			Goal:            a.Behavior.Goal,
			curState:        a.Behavior.curState,
		}
	}
	if a.inventory != nil {
		newAgent.inventory = a.inventory.DeepCopy()
	}
	newAgent.Position = a.Position
	return newAgent
}

func (a *Agent) String() string {
	return fmt.Sprintf("Agent: %s \nInventory %s\n Position %v\n", a.name, a.inventory, a.Position)
}
