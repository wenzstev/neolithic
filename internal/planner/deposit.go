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
	end := start.Copy()

	endLoc, ok := end.Locations[d.locName]
	if !ok {
		return nil // error, no location of that type in State
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // error, no agent of that type in State
	}

	amountOnAgent := endAgent.GetAmount(d.resource)
	if amountOnAgent == 0 {
		return nil
	}

	amountToDeposit := d.amount
	if amountOnAgent < d.amount {
		amountToDeposit = amountOnAgent
	}

	locInv := endLoc.Inventory

	_, ok = locInv[d.resource]
	if !ok {
		locInv[d.resource] = 0
	}

	locInv[d.resource] += amountToDeposit
	newAgent := endAgent.AdjustInventory(d.resource, amountToDeposit)
	end.Agents[agent.Name()] = newAgent

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

// GetStateChange returns the diff in State. It knows nothing about the actual State that might be. Instead it produces State values
// as a diff. So the inventory amount could be negative.
func (d *Deposit) GetStateChange(agent core.Agent) *core.WorldState {
	newAgent := agent.AdjustInventory(d.resource, -d.amount)

	return &core.WorldState{
		Locations: map[string]core.Location{
			d.locName: {
				Name: d.locName,
				Inventory: core.Inventory{
					d.resource: d.amount,
				},
			},
		},

		Agents: map[string]core.Agent{
			newAgent.Name(): newAgent,
		},
	}
}
