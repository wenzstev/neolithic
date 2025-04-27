package attributes

import (
	"Neolithic/internal/actions"
	"Neolithic/internal/core"
)

// CanDepositTo is an attribute that determines if a location can have resources deposited at it.
type CanDepositTo struct {
	// Amount is the amount that can be deposited in a given action
	Amount int
	// Cost is the cost of depositing
	Cost float64
	// Location is the location the CanDepositTo attribute is attached to.
	Location *core.Location
}

// NeedsLocation indicates if CanDepositTo needs a separate location (separate from the location the attribute is attached
// to) to function. It does NOT.
func (c *CanDepositTo) NeedsLocation() bool {
	return false
}

// NeedsResource indicates if CanDepositTo needs a resource. It DOES.
func (c *CanDepositTo) NeedsResource() bool {
	return true
}

// CreateAction provides a concrete action for depositing a resource at the location attached to the CanDepositTo attribute.
func (c *CanDepositTo) CreateAction(params core.CreateActionParams) core.Action {
	return &actions.Deposit{
		DepResource:    params.Resource,
		Amount:         c.Amount,
		ActionLocation: c.Location,
		ActionCost:     c.Cost,
	}
}
