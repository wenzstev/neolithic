package core

import (
	"fmt"
)

// CreateActionParams are the parameters used for the CreateAction function.
type CreateActionParams struct {
	// Location is the location to be passed to the CreateAction func
	Location *Location
	// Resource is the resource to be passed to the CreatAction func
	Resource *Resource
}

// AttributeType is a string representing the type of an Attribute.
type AttributeType string

// Attribute is a representation of a unique aspect of a Resource or a Location
type Attribute interface {
	fmt.Stringer
	// Type returns the AttributeType of the attribute.
	Type() AttributeType
	// CreateAction provides a way to create an action from an interface
	CreateAction(AttributeHolder, CreateActionParams) (Action, error)
	// NeedsLocation determines if an attribute needs a separate location to be relevant
	NeedsLocation() bool
	// NeedsResource determines if an attribute needs a separate resource to be relevant
	NeedsResource() bool
	// Copy returns a copy of the Attribute
	Copy() Attribute
}

// AttributeHolder is the interface representing the holder of an attribute.
// Any type that can hold attributes should implement this interface.
type AttributeHolder interface {
	// Attributes returns the AttributeList held by the AttributeHolder.
	Attributes() AttributeList
}
