package planner

import (
	"Neolithic/internal/core"
	"encoding/gob"
	"fmt"
)

func init() {
	gob.Register(Gather{})
}

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
var _ Action = (*Gather)(nil)

// Perform implements Action.Perform, and simulates the act of gathering a Resource
func (g *Gather) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()
	endLocation, ok := end.Locations[g.ActionLocation.Name]
	if !ok {
		return nil // location must always be in State, this is an error
	}

	amountToGather := minInt(g.Amount, endLocation.Inventory.GetAmount(g.Res))
	if amountToGather <= 0 {
		return nil // fail, no DepResource to gather
	}

	endAgent, ok := end.Agents[agent.Name()]
	if !ok {
		return nil // fail, no agent in State
	}

	endAgentInv := endAgent.Inventory()

	if g.Requires != nil && endAgentInv.GetAmount(g.Requires) <= 0 {
		return nil // fail, does not have the necessary tool
	}

	endAgentInv.AdjustAmount(g.Res, amountToGather)
	endLocation.Inventory.AdjustAmount(g.Res, -amountToGather)
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
func (g *Gather) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{
		{
			Entity:     agent.Name(),
			EntityType: AgentEntity,
			Resource:   g.Res,
			Amount:     g.Amount,
		},
		{
			Entity:     g.ActionLocation.Name,
			EntityType: LocationEntity,
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
