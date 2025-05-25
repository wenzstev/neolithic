package core

import (
	"fmt"
	"sort"
	"strings"
)

// AttributeList defines an interface for a list of attributes.
// It provides methods for managing and accessing attributes.
type AttributeList interface {
	fmt.Stringer
	// UpsertAttribute adds a new attribute or updates an existing one of the same type.
	UpsertAttribute(Attribute)
	// RemoveAttribute removes an attribute of the specified type from the list.
	// The order of the remaining attributes is preserved.
	RemoveAttribute(AttributeType)
	// AttributeByType retrieves an attribute by its type using binary search.
	// It returns nil if no attribute of the specified type is found.
	AttributeByType(AttributeType) Attribute
	// Copy creates a deep copy of the AttributeList.
	Copy() AttributeList
	// List returns all attributes in the list
	List() []Attribute
}

// attributeList is the concrete implementation of the AttributeList interface.
// It stores attributes in a slice, kept sorted by AttributeType for efficient access.
type attributeList []Attribute

// NewAttributeList creates and returns a new, empty AttributeList.
func NewAttributeList() AttributeList {
	return &attributeList{}
}

// UpsertAttribute adds a new attribute to the list or updates an existing one
// if an attribute of the same type is already present.
// The list is kept sorted by AttributeType.
func (a *attributeList) UpsertAttribute(attr Attribute) {
	idx := sort.Search(len(*a), func(i int) bool { return (*a)[i].Type() >= attr.Type() })

	if idx < len(*a) && (*a)[idx].Type() == attr.Type() {
		(*a)[idx] = attr
		return
	}

	*a = append(*a, nil) // Grow the slice by one element (the value doesn't matter here)
	copy((*a)[idx+1:], (*a)[idx:])
	(*a)[idx] = attr
}

// RemoveAttribute removes the attribute of the specified AttributeType from the list.
// If no attribute of the given type is found, the list remains unchanged.
// This operation preserves the order of the remaining elements in the list.
func (a *attributeList) RemoveAttribute(attrType AttributeType) {
	for i, curAttr := range *a {
		if curAttr.Type() == attrType {
			// Remove the element at index i, preserving order
			*a = append((*a)[:i], (*a)[i+1:]...)
			return
		}
	}
}

// AttributeByType searches for and returns an attribute of the specified AttributeType
// using binary search, as the list is kept sorted by AttributeType.
// If no attribute of the given type is found, it returns nil.
func (a *attributeList) AttributeByType(attrType AttributeType) Attribute {
	idx := sort.Search(len(*a), func(i int) bool { return (*a)[i].Type() >= attrType })
	if idx < len(*a) && (*a)[idx].Type() == attrType {
		return (*a)[idx]
	}
	return nil
}

// String returns a string representation of the attributeList.
// It concatenates the string representation of each attribute in the list.
// Returns "{}" if the list is empty.
func (a *attributeList) String() string {
	if len(*a) == 0 {
		return "{}"
	}

	var sb strings.Builder
	for _, curAttr := range *a {
		sb.WriteString(curAttr.String())
	}
	return sb.String()
}

// Copy creates and returns a deep copy of the attributeList.
// Each attribute within the list is also copied.
func (a *attributeList) Copy() AttributeList {
	copyAttrList := make(attributeList, len(*a))
	for i := 0; i < len(*a); i++ {
		copyAttrList[i] = (*a)[i].Copy()
	}
	return &copyAttrList
}

// List returns a copy of all attributes in an AttributeList. It's a copy to prevent modification.
func (a *attributeList) List() []Attribute {
	attrList := make([]Attribute, len(*a))
	for i := 0; i < len(*a); i++ {
		attrList[i] = (*a)[i].Copy()
	}
	return attrList
}
