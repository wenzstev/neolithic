package planner

import "fmt"

// Gather implements Action, and represents the act of gathering a resource
type Gather struct {
	requires *Resource
	resource *Resource
	amount   int
	location *Location
	cost     float64
}

// Perform implements Action.Perform, and simulates the act of gathering a resource
func (g *Gather) Perform(start *State, agent *Agent) *State {

	end := start.Copy()
	locationInv, ok := end.Locations[g.location]
	if !ok {
		return nil // location must always be in state, this is an error
	}
	curResource, ok := locationInv[g.resource]
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

	agentInv, ok := end.Agents[agent]
	if !ok {
		return nil // fail, no agent in state
	}

	_, ok = agentInv[g.requires]
	if g.requires != nil && !ok {
		return nil // fail, does not have the necessary tool
	}

	_, ok = agentInv[g.resource]
	if !ok {
		agentInv[g.resource] = 0
	}

	agentInv[g.resource] += amountToGather
	locationInv[g.resource] -= amountToGather

	if locationInv[g.resource] == 0 {
		delete(locationInv, g.resource)
	}

	return end
}

// Cost implements Action.Cost, and returns the cost of the gather action
func (g *Gather) Cost(_ *Agent) float64 {
	return g.cost
}

// Description implements Action.Description, and provides a brief description of the gather action
func (g *Gather) Description() string {
	return fmt.Sprintf("gather %d %s from %s", g.amount, g.resource.Name, g.location.Name)
}
