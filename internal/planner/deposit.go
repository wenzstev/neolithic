package planner

import (
	"encoding/gob"
	"fmt"

	"Neolithic/internal/core"
)

func init() {
	gob.Register(Deposit{})
}

// Deposit implements Action, and represents the act of depositing a Resource at a location
type Deposit struct {
	// Resource is the Resource being deposited
	Resource *core.Resource
	// Amount is the Amount of Resource being deposited
	Amount int
	// ActionLocation is the Location the Resource is being deposited
	ActionLocation *core.Location
	// ActionCost is the ActionCost of taking the Action
	ActionCost float64
}

// Force Deposit to implement Action
var _ Action = (*Deposit)(nil)

// Perform implements Action.Perform, and simulates the act of depositing a Resource in a location
func (d *Deposit) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()

	endLoc, ok := end.Locations[d.ActionLocation.Name]
	if !ok {
		return nil // error, no location of that type in State
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // error, no agent of that type in State
	}

	endAgentInv := endAgent.Inventory()
	amountToDeposit := minInt(endAgentInv.GetAmount(d.Resource), d.Amount)

	if amountToDeposit <= 0 {
		return nil // fail, no Resource to deposit
	}

	endLoc.Inventory.AdjustAmount(d.Resource, amountToDeposit)
	endAgentInv.AdjustAmount(d.Resource, -amountToDeposit)

	return end
}

// Cost implements Action.Cost, and returns the energy ActionCost of depositing the Resource.
func (d *Deposit) Cost(_ core.Agent) float64 {
	return d.ActionCost // TODO: more dynamic ActionCost
}

// Description implements Action.Description, and returns a string representation of the Action.
func (d *Deposit) Description() string {
	return fmt.Sprintf("deposit %d %s at %s", d.Amount, d.Resource.Name, d.ActionLocation)
}

func (d *Deposit) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{
		{
			Entity:     agent.Name(),
			EntityType: AgentEntity,
			Resource:   d.Resource,
			Amount:     -d.Amount,
		},
		{
			Entity:     d.ActionLocation.Name,
			EntityType: LocationEntity,
			Resource:   d.Resource,
			Amount:     d.Amount,
		},
	}
}

func (d *Deposit) Location() *core.Location {
	return d.ActionLocation
}
