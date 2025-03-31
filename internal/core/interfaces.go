package core

import (
	"Neolithic/internal/astar"
	"fmt"
)

// Agent is the core interface for all packages that require an agent
type Agent interface {
	fmt.Stringer
	// Name returns the name of the agent
	Name() string
	// DeepCopy returns a deep copy of the agent
	DeepCopy() Agent
	// Inventory returns the agent's current inventory
	Inventory() Inventory
}

// Locatable is an interface that represents anything with a location associated with it
type Locatable interface {
	// Location returns the location of the entity
	Location() *Location
}

type Grid interface {
	CellAt(coord Coord) Cell
}

type Cell interface {
	astar.Node
	Coord() Coord
}
