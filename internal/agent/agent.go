package agent

import (
	"Neolithic/internal/core"
	"fmt"
)

// Agent struct represents an agent in the simulation world that can interact with its environment.
// It contains the agent's name, behavior patterns, inventory and position information.
type Agent struct {
	// name is the name of the agent
	name string
	// Behavior holds the agent's decision-making processes and planning capabilities
	Behavior *Behavior
	// inventory stores the items and resources the agent currently possesses
	inventory core.Inventory
	// Position represents the agent's current location in the world using coordinates
	Position core.Coord
}

// Ensure Agent implements core.Agent interface
var _ core.Agent = (*Agent)(nil)

// Name returns the name of the agent
func (a *Agent) Name() string {
	return a.name
}

// Inventory returns the agent's current inventory
func (a *Agent) Inventory() core.Inventory {
	return a.inventory
}

// DeepCopy creates a deep copy of the agent and returns it
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

// String returns a string representation of the agent including name, inventory and position
func (a *Agent) String() string {
	return fmt.Sprintf("Agent: %s \nInventory %s\n Position %v\n", a.name, a.inventory, a.Position)
}
