package attributes

import (
	"Neolithic/internal/core"
	"fmt"
)

// Gather implements Action, and represents the act of gathering a Resource
type Gather struct {
	// Requires is an optional Resource that is required to perform the gather
	Requires *core.Resource
	// Res is the Resource being gathered
	Res *core.Resource
	// Amount is the Amount of the Resource being gathered
	Amount int
	// ActionLocation is the Location where the Resource is being gathered
	ActionLocation *core.Location
	// ActionCost is the cost of taking the Action
	ActionCost float64
}

// Force Gather to implement Action
var _ core.Action = (*Gather)(nil)

// Perform implements Action.Perform, and simulates the act of gathering a Resource
func (g *Gather) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	gatherLocation, ok := start.GetLocation(g.ActionLocation.Name)
	if !ok {
		return nil
	}
	amountToGather := minInt(g.Amount, gatherLocation.Inventory.GetAmount(g.Res))
	if amountToGather <= 0 {
		return nil // fail, no DepResource to gather
	}

	if g.Requires != nil && agent.Inventory().GetAmount(g.Requires) <= 0 {
		return nil // fail, does not have the necessary tool
	}

	startAgent, ok := start.GetAgent(agent.Name())
	if !ok {
		return nil
	}

	endAgent := startAgent.DeepCopy()
	endAgentInv := endAgent.Inventory()
	endLocation := gatherLocation.DeepCopy()

	endAgentInv.AdjustAmount(g.Res, amountToGather)
	endLocation.Inventory.AdjustAmount(g.Res, -amountToGather)

	end := start.ShallowCopy()
	end.Locations[endLocation.Name] = endLocation
	end.Agents[endAgent.Name()] = endAgent

	return end
}

// Cost implements Action.Cost, and returns the ActionCost of the gather Action
func (g *Gather) Cost(_ core.Agent) float64 {
	return g.ActionCost
}

// Description implements Action.Description, and provides a brief description of the gather Action
func (g *Gather) Description() string {
	return fmt.Sprintf("gather %d %s from %s", g.Amount, g.Res.Name, g.ActionLocation)
}

// GetChanges generates a list of state changes resulting from gathering a resource by a specified agent.
func (g *Gather) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{
		{
			Entity:     agent.Name(),
			EntityType: core.AgentEntity,
			Resource:   g.Res,
			Amount:     g.Amount,
		},
		{
			Entity:     g.ActionLocation.Name,
			EntityType: core.LocationEntity,
			Resource:   g.Res,
			Amount:     -g.Amount,
		},
	}
}

// minInt returns the smaller of two integer values a and b.
func minInt(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Location retrieves the location where the resource is being gathered.
func (g *Gather) Location() *core.Location {
	return g.ActionLocation
}

// Resource returns the DepResource associated with the Gather action.
func (g *Gather) Resource() *core.Resource {
	return g.Res
}
