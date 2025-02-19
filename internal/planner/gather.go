package planner

import "fmt"

// Gather implements Action, and represents the act of gathering a resource
type Gather struct {
	// requires is an optional resource that is required to perform the gather
	requires *Resource
	// resource is the Resource being gathered
	resource *Resource
	// amount is the amount of the resource being gathered
	amount int
	// location is the Location where the resource is being gathered
	location *Location
	// cost is the cost of taking the Action
	cost float64
}

// Force Gather to implement Action
var _ Action = (*Gather)(nil)

// Perform implements Action.Perform, and simulates the act of gathering a resource
func (g *Gather) Perform(start *State, agent *Agent) *State {
	end := start.Copy()
	locationInv, ok := end.Locations[g.location]
	if !ok {
		return nil // location must always be in State, this is an error
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
		return nil // fail, no agent in State
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

// Cost implements Action.Cost, and returns the cost of the gather Action
func (g *Gather) Cost(_ *Agent) float64 {
	return g.cost
}

// Description implements Action.Description, and provides a brief description of the gather Action
func (g *Gather) Description() string {
	return fmt.Sprintf("gather %d %s from %s", g.amount, g.resource.Name, g.location.Name)
}

// GetStateChange implements Action.GetStateChange and returns the change in State from the Action. For Gather, this
// means the designated amount is removed from the location and given to the agent.
func (g *Gather) GetStateChange(agent *Agent) *State {
	return &State{
		Locations: map[*Location]Inventory{
			g.location: {
				g.resource: g.amount * -1,
			},
		},
		Agents: map[*Agent]Inventory{
			agent: {
				g.resource: g.amount,
			},
		},
	}
}
