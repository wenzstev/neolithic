package attributes

import (
	"Neolithic/internal/actions"
	"Neolithic/internal/core"
)

// CanGatherFrom is an Attribute that indiciates if a location can be gathered from.
type CanGatherFrom struct {
	// Amount is the amount that can be gathered
	Amount int // TODO should amount be a function of the resource attributes?
	// Cost is the cost to gather
	Cost float64 // TODO should cost be a function of the resource attributes?
	// Location is the location that will be gathered from. TODO move location to register logic?
	Location *core.Location
}

// NeedsLocation indicates whether CanGatherFrom requires an additional location to create an action. It does NOT.
func (c *CanGatherFrom) NeedsLocation() bool {
	return false
}

// NeedsResource indicates whether CanGatherFrom requires a resource to create an action. It DOES.
func (c *CanGatherFrom) NeedsResource() bool {
	return true
}

// CreateAction creates a specific action of gathering a specific resource from the location
func (c *CanGatherFrom) CreateAction(params core.CreateActionParams) core.Action {
	return &actions.Gather{
		Res:            params.Resource,
		Amount:         c.Amount,
		ActionLocation: c.Location,
		ActionCost:     c.Cost,
	}
}
