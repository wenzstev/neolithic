package planner

import (
	"Neolithic/internal/core"
	"fmt"
)

// Gather implements Action, and represents the act of gathering a resource
type Gather struct {
	// requires is an optional resource that is required to perform the gather
	requires *core.Resource
	// resource is the Resource being gathered
	resource *core.Resource
	// amount is the amount of the resource being gathered
	amount int
	// locName is the Location where the resource is being gathered
	locName string
	// cost is the cost of taking the Action
	cost float64
}

// Force Gather to implement Action
var _ Action = (*Gather)(nil)

// Perform implements Action.Perform, and simulates the act of gathering a resource
func (g *Gather) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.Copy()
	location, ok := end.Locations[g.locName]
	if !ok {
		return nil // location must always be in State, this is an error
	}
	curResource, ok := location.Inventory[g.resource]
	if !ok {
		return nil // fail, not possible to gather from this location as there is not resource
	}

	if curResource <= 0 {
		return nil // fail, no resource to gather
	}

	amountToGather := g.amount
	if curResource < g.amount {
		amountToGather = curResource
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // fail, no agent in State
	}

	if endAgent.GetAmount(g.resource) <= 0 {
		return nil // fail, does not have the necessary tool
	}

	newAgent := endAgent.AdjustInventory(g.resource, amountToGather)
	end.Agents[agent.Name()] = newAgent

	location.Inventory[g.resource] -= amountToGather

	if location.Inventory[g.resource] == 0 {
		delete(location.Inventory, g.resource)
	}

	return end
}

// Cost implements Action.Cost, and returns the cost of the gather Action
func (g *Gather) Cost(_ core.Agent) float64 {
	return g.cost
}

// Description implements Action.Description, and provides a brief description of the gather Action
func (g *Gather) Description() string {
	return fmt.Sprintf("gather %d %s from %s", g.amount, g.resource.Name, g.locName)
}

// GetStateChange implements Action.GetStateChange and returns the change in State from the Action. For Gather, this
// means the designated amount is removed from the location and given to the agent.
func (g *Gather) GetStateChange(agent core.Agent) *core.WorldState {
	newAgent := agent.AdjustInventory(g.resource, g.amount)

	return &core.WorldState{
		Locations: map[string]core.Location{
			g.locName: {
				Name: g.locName,
				Inventory: core.Inventory{
					g.resource: -g.amount,
				},
			},
		},
		Agents: map[string]core.Agent{
			agent.Name(): newAgent,
		},
	}
}
