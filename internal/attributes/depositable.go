package attributes

import (
	"Neolithic/internal/actions"
	"Neolithic/internal/core"
)

// CanDeposit is an attribute that determines if a location can have resources deposited at it.
type CanDeposit struct {
	// Amount is the amount that can be deposited in a given action
	Amount int
	// Cost is the cost of depositing
	Cost float64
	// Location is the location the CanDeposit attribute is attached to.
	Location *core.Location
}

// NeedsLocation indicates if CanDeposit needs a separate location (separate from the location the attribute is attached
// to) to function. It does NOT.
func (c *CanDeposit) NeedsLocation() bool {
	return false
}

// NeedsResource indicates if CanDeposit needs a resource. It DOES.
func (c *CanDeposit) NeedsResource() bool {
	return true
}

// CreateAction provides a concrete action for depositing a resource at the location attached to the CanDeposit attribute.
func (c *CanDeposit) CreateAction(params core.CreateActionParams) core.Action {
	return &actions.Deposit{
		DepResource:    params.Resource,
		Amount:         c.Amount,
		ActionLocation: c.Location,
		ActionCost:     c.Cost,
	}
}
