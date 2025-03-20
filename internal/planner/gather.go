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
	end := start.DeepCopy()
	endLocation, ok := end.Locations[g.locName]
	if !ok {
		return nil // location must always be in State, this is an error
	}

	amountToGather := minInt(g.amount, endLocation.Inventory.GetAmount(g.resource))
	if amountToGather <= 0 {
		return nil // fail, no resource to gather
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // fail, no agent in State
	}

	endAgentInv := endAgent.Inventory()

	if endAgentInv.GetAmount(g.resource) <= 0 {
		return nil // fail, does not have the necessary tool
	}

	endAgentInv.AdjustAmount(g.resource, amountToGather)
	endLocation.Inventory.AdjustAmount(g.resource, -amountToGather)
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

func (g *Gather) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{
		{
			Entity:     agent.Name(),
			EntityType: AgentEntity,
			Resource:   g.resource,
			Amount:     g.amount,
		},
		{
			Entity:     g.locName,
			EntityType: LocationEntity,
			Resource:   g.resource,
			Amount:     -g.amount,
		},
	}
}

func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}
