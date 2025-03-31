package core

import (
	"fmt"

	"Neolithic/internal/astar"
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

// Grid is an interface that represents a grid of cells
type Grid interface {
	// CellAt returns the cell at the given coordinate
	CellAt(coord Coord) Cell
}

// Cell is an interface that represents a cell in a grid. It is expected to implement the astar.Node interface for pathfinding.
type Cell interface {
	// require implementing astar.Node
	astar.Node
	// Coord returns the coordinate of the cell
	Coord() Coord
}
