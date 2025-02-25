package world

type Coord struct {
	X, Y int
}

// IsWithin returns whether the other Coord is within the given distance
func (c *Coord) IsWithin(other *Coord, dis float64) bool {
	// todo implement me
	panic("unimplemented")
}

// Locatable is an interface that provides a location.
type Locatable interface { // TODO find better location
	Location() Coord
}
