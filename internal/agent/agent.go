package agent

import (
	"encoding/gob"
	"log/slog"
	"strings"

	"Neolithic/internal/core"
)

func init() {
	gob.Register(Agent{})
}

// Agent struct represents an Agent in the simulation world that can interact with its environment.
// It contains the Agent's name, behavior patterns, inventory and position information.
type Agent struct {
	// name is the name of the Agent
	name string
	// Behavior holds the Agent's decision-making processes and planning capabilities
	Behavior *Behavior
	// inventory stores the items and resources the Agent currently possesses
	inventory core.Inventory
	// Position represents the Agent's current location in the world using coordinates
	Position core.Coord
}

// Ensure Agent implements core.Agent interface
var _ core.Agent = (*Agent)(nil)

// Name returns the name of the Agent
func (a *Agent) Name() string {
	return a.name
}

// Inventory returns the Agent's current inventory
func (a *Agent) Inventory() core.Inventory {
	return a.inventory
}

// DeepCopy creates a deep copy of the Agent and returns it
func (a *Agent) DeepCopy() core.Agent {
	newAgent := &Agent{}
	newAgent.name = a.name
	if a.Behavior != nil {
		newAgent.Behavior = &Behavior{
			PossibleActions: a.Behavior.PossibleActions,
			CurPlan:         a.Behavior.CurPlan,
			CurState:        a.Behavior.CurState,
			GoalEngine:      a.Behavior.GoalEngine,
		}
	}
	if a.inventory != nil {
		newAgent.inventory = a.inventory.DeepCopy()
	}
	newAgent.Position = a.Position
	return newAgent
}

// String returns a string representation of the Agent including name, inventory and position
func (a *Agent) String() string {
	var sb strings.Builder
	sb.WriteString("Agent: ")
	sb.WriteString(a.name)
	sb.WriteString("\nInventory ")
	sb.WriteString(a.inventory.String())
	sb.WriteString("\nPosition ")
	sb.WriteString(a.Position.String())
	sb.WriteString("\n")
	return sb.String()
}

func (a *Agent) Tick(worldState *core.WorldState, deltaTime float64) (*core.WorldState, error) {
	return a.Behavior.CurState.Execute(worldState, deltaTime)
}

func NewAgent(name string, logger *slog.Logger) *Agent {
	newAgent := &Agent{
		Behavior:  &Behavior{},
		name:      name,
		inventory: core.NewInventory(),
	}
	newAgent.Behavior.CurState = &Idle{agent: newAgent, logger: logger}
	return newAgent
}
