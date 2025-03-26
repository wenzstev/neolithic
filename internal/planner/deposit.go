package planner

import (
	"Neolithic/internal/core"
	"fmt"
)

// Deposit implements Action, and represents the act of depositing a resource at a location
type Deposit struct {
	// resource is the Resource being deposited
	resource *core.Resource
	// amount is the amount of resource being deposited
	amount int
	// location is the Location the resource is being deposited
	locName string
	// cost is the cost of taking the Action
	cost float64
}

// Force Deposit to implement Action
var _ Action = (*Deposit)(nil)

// Perform implements Action.Perform, and simulates the act of depositing a resource in a location
func (d *Deposit) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()

	endLoc, ok := end.Locations[d.locName]
	if !ok {
		return nil // error, no location of that type in State
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // error, no agent of that type in State
	}

	endAgentInv := endAgent.Inventory()
	amountToDeposit := minInt(endAgentInv.GetAmount(d.resource), d.amount)

	if amountToDeposit <= 0 {
		return nil // fail, no resource to deposit
	}

	endLoc.Inventory.AdjustAmount(d.resource, amountToDeposit)
	endAgentInv.AdjustAmount(d.resource, -amountToDeposit)

	return end
}

// Cost implements Action.Cost, and returns the energy cost of depositing the resource.
func (d *Deposit) Cost(_ core.Agent) float64 {
	return d.cost // TODO: more dynamic cost
}

// Description implements Action.Description, and returns a string representation of the Action.
func (d *Deposit) Description() string {
	return fmt.Sprintf("deposit %d %s at %s", d.amount, d.resource.Name, d.locName)
}

func (d *Deposit) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{
		{
			Entity:     agent.Name(),
			EntityType: AgentEntity,
			Resource:   d.resource,
			Amount:     -d.amount,
		},
		{
			Entity:     d.locName,
			EntityType: LocationEntity,
			Resource:   d.resource,
			Amount:     d.amount,
		},
	}
}
