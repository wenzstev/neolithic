package core

import "fmt"

type Coord struct {
	X, Y int
}

func (c *Coord) String() string {
	return fmt.Sprintf("(%d, %d)", c.X, c.Y)
}

type Resource struct {
	Name string
}

func (r *Resource) String() string {
	return fmt.Sprintf("Resource: %s", r.Name)
}

type Location struct {
	Name      string
	Inventory Inventory
	Coord     Coord
}

func (l *Location) String() string {
	return fmt.Sprintf("Location: %s\nCoordinates: %v\nInventory: %v", l.Name, l.Coord, l.Inventory)
}

func (l *Location) DeepCopy() *Location {
	return &Location{
		Name:      l.Name,
		Inventory: l.Inventory.DeepCopy(),
		Coord:     l.Coord,
	}
}
