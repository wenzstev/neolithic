package attributes

import (
	"Neolithic/internal/core"
	"errors"
	"strconv"
	"strings"
)

const (
	// WeightAttributeType is the attribute type that corresponds to the Weight attribute
	WeightAttributeType core.AttributeType = "weight"

	// defaultGatherAmount is the default amount that the created Gather action will have
	defaultGatherAmount = 1
)

// Weight is an Attribute that indicates the weight of a resource. It corresponds to the Gather action, indicating that
// a resource can be gathered and deposited
type Weight struct {
	Amount float64
}

// NeedsLocation indicates whether Weight requires an additional location to create an action.
func (w *Weight) NeedsLocation() bool {
	return true
}

// NeedsResource indicates whether Weight requires a resource to create an action.
func (w *Weight) NeedsResource() bool {
	return false
}

// CreateAction creates a specific action of gathering a specific resource from the location
func (w *Weight) CreateAction(holder core.AttributeHolder, params core.CreateActionParams) (core.Action, error) {
	res, ok := holder.(*core.Resource)
	if !ok {
		return nil, errors.New("weight can only be applied to a resource")
	}
	if params.Location == nil {
		return nil, errors.New("CreateAction was called for a resource with a weight but no location")
	}

	return &Gather{
		Res:            res,
		Amount:         defaultGatherAmount,
		ActionLocation: params.Location,
		ActionCost:     w.Amount,
	}, nil
}

// Type returns the WeightAttributeType for the Weight attribute
func (w *Weight) Type() core.AttributeType {
	return WeightAttributeType
}

// Copy returns a copy of the weight attribute
func (w *Weight) Copy() core.Attribute {
	return &Weight{Amount: w.Amount}
}

// String returns a string representation fo the weight attribute
func (w *Weight) String() string {
	var sb strings.Builder
	sb.WriteString("Weight: ")
	sb.WriteString(strconv.FormatFloat(w.Amount, 'f', -1, 64))
	return sb.String()
}
