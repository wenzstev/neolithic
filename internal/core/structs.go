package core

import "fmt"

// Coord represents a point in a 2D space with integer coordinates.
type Coord struct {
	X, Y int
}

// String returns a string representation of the Coord in the format "(x, y)".
func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

// Resource represents a resource in the simulation world.
type Resource struct {
	Name string
}

// String returns a string representation of the Resource in the format "Resource: <name>".
func (r *Resource) String() string {
	return fmt.Sprintf("Resource: %s", r.Name)
}

// Location represents a location in the simulation world.
type Location struct {
	Name      string
	Inventory Inventory
	Coord     Coord
}

// String returns a string representation of the Location in the format "Location: <name>\nCoordinates: <coordinates>\nInventory: <inventory>".
func (l *Location) String() string {
	return fmt.Sprintf("Location: %s\nCoordinates: %v\nInventory: %v", l.Name, l.Coord, l.Inventory)
}

// DeepCopy creates a deep copy of the Location.
func (l *Location) DeepCopy() *Location {
	return &Location{
		Name:      l.Name,
		Inventory: l.Inventory.DeepCopy(),
		Coord:     l.Coord,
	}
}
