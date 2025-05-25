package core

import (
	"strings"
)

// Location represents a location in the simulation world.
type Location struct {
	// Name is the unique identifier for the location.
	Name string
	// Inventory holds the resources present at this location.
	Inventory Inventory
	// Coord represents the geographical coordinates of the location.
	Coord Coord
	// attributes is a list of special characteristics or properties of the location.
	attributes AttributeList
}

// LocationOption is a functional option type for configuring a new Location.
// It allows for flexible and extensible Location creation.
type LocationOption func(*Location)

// NewLocation creates and returns a new Location with the given name and coordinates.
// It applies any provided LocationOption functions to customize the new Location.
func NewLocation(name string, coord Coord, opts ...LocationOption) *Location {
	loc := &Location{
		Name:       name,
		Coord:      coord,
		Inventory:  NewInventory(),
		attributes: NewAttributeList(),
	}
	for _, opt := range opts {
		opt(loc)
	}
	return loc
}

// WithInventory is a LocationOption that initializes the Location's inventory
// with the provided InventoryEntry items.
func WithInventory(entries ...InventoryEntry) LocationOption {
	locInv := NewInventory()
	for _, entry := range entries {
		locInv.AdjustAmount(entry.Resource, entry.Amount)
	}
	return func(l *Location) {
		l.Inventory = locInv
	}
}

// WithAttributes is a LocationOption that initializes the Location's attributes
// with the provided Attribute items.
func WithAttributes(attributes ...Attribute) LocationOption {
	attrList := NewAttributeList()
	for _, attr := range attributes {
		attrList.UpsertAttribute(attr)
	}
	return func(l *Location) {
		l.attributes = attrList
	}
}

// String returns a string representation of the Location in the format
// "Location: <name>\nCoordinates: <coordinates>\nInventory: <inventory>\nAttributes: <attributes>".
func (l *Location) String() string {
	var sb strings.Builder
	sb.WriteString("Location: ")
	sb.WriteString(l.Name)
	sb.WriteString("\n")
	sb.WriteString("Coordinates: ")
	sb.WriteString(l.Coord.String())
	sb.WriteString("\n")
	sb.WriteString("Inventory: ")
	if l.Inventory != nil { // Guard against nil inventory if NewInventory wasn't called or was overwritten
		sb.WriteString(l.Inventory.String())
	} else {
		sb.WriteString("{}") // Or some other placeholder for nil inventory
	}
	sb.WriteString("\n")
	sb.WriteString("Attributes: ")
	if l.attributes != nil {
		sb.WriteString(l.attributes.String())
	} else {
		sb.WriteString("{}")
	}
	return sb.String()
}

// DeepCopy creates a deep copy of the Location.
// This includes copying the Name, Coordinates, Inventory, and Attributes.
func (l *Location) DeepCopy() *Location {
	// Ensure attributes and inventory are not nil before copying
	var copiedAttributes AttributeList
	if l.attributes != nil {
		copiedAttributes = l.attributes.Copy()
	} else {
		copiedAttributes = NewAttributeList() // Or handle as error/nil if appropriate
	}

	var copiedInventory Inventory
	if l.Inventory != nil {
		copiedInventory = l.Inventory.DeepCopy()
	} else {
		copiedInventory = NewInventory() // Or handle as error/nil if appropriate
	}

	return &Location{
		Name:       l.Name,
		Inventory:  copiedInventory,
		Coord:      l.Coord,
		attributes: copiedAttributes,
	}
}

// Attributes returns the AttributeList associated with the Location.
func (l *Location) Attributes() AttributeList {
	return l.attributes
}
