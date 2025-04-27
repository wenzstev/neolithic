package core

import (
	"fmt"
	"strings"
)

// Resource represents a resource in the simulation world.
type Resource struct {
	Name       string
	Attributes []Attribute
}

// String returns a string representation of the Resource in the format "Resource: <name>".
func (r *Resource) String() string {
	return fmt.Sprintf("Resource: %s", r.Name)
}

// Location represents a location in the simulation world.
type Location struct {
	Name       string
	Inventory  Inventory
	Coord      Coord
	Attributes []Attribute
}

// String returns a string representation of the Location in the format "Location: <name>\nCoordinates: <coordinates>\nInventory: <inventory>".
func (l *Location) String() string {
	var sb strings.Builder
	sb.WriteString("Location: ")
	sb.WriteString(l.Name)
	sb.WriteString("\n")
	sb.WriteString("Coordinates: ")
	sb.WriteString(l.Coord.String())
	sb.WriteString("\n")
	sb.WriteString("Inventory: ")
	sb.WriteString(l.Inventory.String())
	return sb.String()
}

// DeepCopy creates a deep copy of the Location.
func (l *Location) DeepCopy() *Location {
	return &Location{
		Name:      l.Name,
		Inventory: l.Inventory.DeepCopy(),
		Coord:     l.Coord,
	}
}
