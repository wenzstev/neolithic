package agent

import "Neolithic/internal/core"

// Path represents a path of core.Coord that an agent can follow.
type Path interface {
	// NextCoord returns the next coordinate in the path
	NextCoord() core.Coord
	// IsComplete returns true if the path is complete
	IsComplete() bool
}

// CoordPath implements Path
type CoordPath struct {
	// coords is the list of coordinates in the path
	coords []core.Coord
	// index is the index of the next coordinate to return
	index int
}

// NewCoordPath creates a new CoordPath
func NewCoordPath(coords []core.Coord) *CoordPath {
	return &CoordPath{
		coords: coords,
		index:  1, // start at 1 because the agent is on the first coordinate
	}
}

// NextCoord returns the next coordinate in the path
func (p *CoordPath) NextCoord() core.Coord {
	if p.IsComplete() {
		panic("attempting to get next coordinate from completed path")
	}
	coord := p.coords[p.index]
	p.index++
	return coord
}

// IsComplete returns true if the path is complete
func (p *CoordPath) IsComplete() bool {
	return p.index >= len(p.coords)
}
