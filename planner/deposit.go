package planner

import "fmt"

// Deposit implements Action, and represents the act of depositing a resource at a location
type Deposit struct {
	resource *Resource
	amount   int
	location *Location
}

// Force Deposit to implement Action
var _ Action = (*Deposit)(nil)

// Perform implements Action.Perform, and simulates the act of depositing a resource in a location
func (d *Deposit) Perform(start *State, agent *Agent) *State {
	end := start.Copy()

	locationInv, ok := end.Locations[d.location]
	if !ok {
		return nil // error, no location of that type in state
	}

	agentInv, ok := end.Agents[agent]
	if !ok {
		return nil // error, no agent of that type in state
	}

	amountOnAgent, ok := agentInv[d.resource]
	if !ok {
		return nil // error, agent doesn't have the required resource
	}

	amountToDeposit := d.amount
	if amountOnAgent < d.amount {
		amountToDeposit = amountOnAgent
	}

	_, ok = locationInv[d.resource]
	if !ok {
		locationInv[d.resource] = 0
	}

	locationInv[d.resource] += amountToDeposit
	agentInv[d.resource] -= amountToDeposit

	if agentInv[d.resource] == 0 {
		delete(agentInv, d.resource)
	}

	return end
}

// Cost implements Action.Cost, and returns the energy cost of depositing the resource.
func (d *Deposit) Cost(_ *Agent) float64 {
	return float64(1) // TODO: more dynamic cost
}

// Description implements Action.Description, and returns a string representation of the action.
func (d *Deposit) Description() string {
	return fmt.Sprintf("deposit %d %s at %s", d.amount, d.resource.Name, d.location.Name)
}
