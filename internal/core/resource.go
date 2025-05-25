package core

import (
	"strings"
)

// Resource represents a resource in the simulation world.
type Resource struct {
	// Name is the unique identifier for the resource.
	Name string
	// attributes is a list of Attributes of the resource.
	attributes AttributeList
}

// NewResource creates and returns a new Resource with the given name.
// It applies any provided ResourceOption functions to customize the new Resource.
// By default, a new resource will have an empty list of attributes.
func NewResource(name string, opts ...ResourceOption) *Resource {
	res := &Resource{
		Name:       name,
		attributes: NewAttributeList(),
	}
	for _, opt := range opts {
		opt(res)
	}
	return res
}

// ResourceOption is a functional option type for configuring a new Resource.
// It allows for flexible and extensible Resource creation.
type ResourceOption func(*Resource)

// WithResourceAttributes is a ResourceOption that initializes the Resource's attributes
// with the provided Attribute items.
func WithResourceAttributes(attributes ...Attribute) ResourceOption {
	return func(r *Resource) {
		// Ensure attributes list is initialized
		if r.attributes == nil {
			r.attributes = NewAttributeList()
		}
		for _, attr := range attributes {
			r.attributes.UpsertAttribute(attr)
		}
	}
}

// String returns a string representation of the Resource in the format
// "Resource: <name>\nAttributes: <attributes>".
func (r *Resource) String() string {
	var sb strings.Builder
	sb.WriteString("Resource: ")
	sb.WriteString(r.Name)
	sb.WriteString("\n")
	sb.WriteString("Attributes: ")
	// Ensure attributes is not nil before calling String() on it.
	if r.attributes != nil {
		sb.WriteString(r.attributes.String())
	} else {
		sb.WriteString("{}")
	}
	return sb.String()
}

// Attributes returns the AttributeList associated with the Resource.
func (r *Resource) Attributes() AttributeList {
	return r.attributes
}
