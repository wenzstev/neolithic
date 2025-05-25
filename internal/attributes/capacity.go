package attributes

import (
	"errors"
	"strconv"
	"strings"

	"Neolithic/internal/core"
)

// CapacityAttributeType defines the string identifier for the Capacity attribute.
// It is used to register and retrieve the Capacity attribute from attribute collections.
const CapacityAttributeType core.AttributeType = "capacity"

// defaultDepositAmount is the standard amount of a resource considered deposited
// when a Deposit action is created. This is an internal constant.
const defaultDepositAmount = 1.0

// Capacity is an attribute that determines if a location can have resources deposited at it
// and defines the maximum weight of resources it can hold.
type Capacity struct {
	// Size represents the maximum weight capacity of the location.
	// Resources heavier than this size cannot be deposited.
	Size float64
}

// NeedsLocation indicates if the attribute's primary action
// requires a separate target location distinct from the location the attribute is attached to.
// For Capacity, this is false, as resources are deposited *at* the location with capacity.
func (c *Capacity) NeedsLocation() bool {
	return false
}

// NeedsResource indicates if the attribute's primary action
// requires a resource to be involved. For Capacity, this is true, as the action is to deposit a resource.
func (c *Capacity) NeedsResource() bool {
	return true
}

// CreateAction generates a specific Action instance for depositing a resource at the location
// to which this Capacity attribute is attached.
// It checks if the holder is a valid Location, if a resource is provided, if the resource has weight,
// and if the resource's weight does not exceed the location's capacity.
// If all conditions are met, a Deposit action is returned. Otherwise, it may return nil and an error,
// or nil and no error if the action is simply not possible (e.g., resource too heavy).
func (c *Capacity) CreateAction(holder core.AttributeHolder, params core.CreateActionParams) (core.Action, error) {
	loc, ok := holder.(*core.Location)
	if !ok {
		return nil, errors.New("capacity can only be applied to a location")
	}

	if params.Resource == nil {
		return nil, errors.New("CreateAction was called for a location with capacity, but no resource was provided")
	}

	resAttrs := params.Resource.Attributes()
	weightAttr := resAttrs.AttributeByType(WeightAttributeType)
	if weightAttr == nil {
		return nil, nil
	}
	weight, ok := weightAttr.(*Weight) // Assumes Weight struct is defined elsewhere
	if !ok {
		return nil, errors.New("resource attribute WeightAttributeType is not of type Weight")
	}

	if weight.Amount > c.Size {
		return nil, nil
	}

	return &Deposit{
		DepResource:    params.Resource,
		Amount:         defaultDepositAmount,
		ActionLocation: loc,
		ActionCost:     weight.Amount,
	}, nil
}

// Type returns the specific AttributeType for the Capacity attribute.
func (c *Capacity) Type() core.AttributeType {
	return CapacityAttributeType
}

// Copy creates a new instance of the Capacity attribute with the same Size.
func (c *Capacity) Copy() core.Attribute {
	return &Capacity{Size: c.Size}
}

// String provides a human-readable string representation of the Capacity attribute,
func (c *Capacity) String() string {
	var sb strings.Builder
	sb.WriteString("Capacity: ")
	sb.WriteString(strconv.FormatFloat(c.Size, 'f', -1, 64)) // 'f' for standard decimal notation
	return sb.String()
}
