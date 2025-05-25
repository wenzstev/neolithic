package attributes

import (
	"fmt"

	"Neolithic/internal/core"
)

// Deposit implements Action, and represents the act of depositing a Resource at a location
type Deposit struct {
	// DepResource is the Resource being deposited
	DepResource *core.Resource
	// Amount is the Amount of Resource being deposited
	Amount int
	// ActionLocation is the Location the Resource is being deposited
	ActionLocation *core.Location
	// ActionCost is the ActionCost of taking the Action
	ActionCost float64
}

// Force Deposit to implement Action
var _ core.Action = (*Deposit)(nil)

// Perform implements Action.Perform, and simulates the act of depositing a Resource in a location
func (d *Deposit) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()

	endLoc, ok := end.GetLocation(d.ActionLocation.Name)
	if !ok {
		return nil // error, no location of that type in State
	}

	endAgent, ok := end.GetAgent(agent.Name())
	if !ok {
		return nil // error, no agent of that type in State
	}

	endAgentInv := endAgent.Inventory()
	amountToDeposit := minInt(endAgentInv.GetAmount(d.DepResource), d.Amount)

	if amountToDeposit <= 0 {
		return nil // fail, no DepResource to deposit
	}

	endLoc.Inventory.AdjustAmount(d.DepResource, amountToDeposit)
	endAgentInv.AdjustAmount(d.DepResource, -amountToDeposit)

	return end
}

// Cost implements Action.Cost, and returns the energy ActionCost of depositing the Resource.
func (d *Deposit) Cost(_ core.Agent) float64 {
	return d.ActionCost // TODO: more dynamic ActionCost
}

// Description implements Action.Description, and returns a string representation of the Action.
func (d *Deposit) Description() string {
	return fmt.Sprintf("deposit %d %s at %s", d.Amount, d.DepResource.Name, d.ActionLocation)
}

// GetChanges generates a list of state changes representing the effects of depositing a resource on the agent and location.
func (d *Deposit) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{
		{
			Entity:     agent.Name(),
			EntityType: core.AgentEntity,
			Resource:   d.DepResource,
			Amount:     -d.Amount,
		},
		{
			Entity:     d.ActionLocation.Name,
			EntityType: core.LocationEntity,
			Resource:   d.DepResource,
			Amount:     d.Amount,
		},
	}
}

// Location returns the Location where the deposit action takes place.
func (d *Deposit) Location() *core.Location {
	return d.ActionLocation
}

// Resource returns the Resource associated with the Deposit.
func (d *Deposit) Resource() *core.Resource {
	return d.DepResource
}
