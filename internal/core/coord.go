package core

import (
	"strconv"
	"strings"
)

// Coord represents a point in a 2D space with integer coordinates.
type Coord struct {
	X, Y int
}

// String returns a string representation of the Coord in the format "(x, y)".
func (c Coord) String() string {
	var sb strings.Builder
	sb.WriteString("(")
	sb.WriteString(strconv.Itoa(c.X))
	sb.WriteString(", ")
	sb.WriteString(strconv.Itoa(c.Y))
	sb.WriteString(")")
	return sb.String()
}

// IsWithin checks if the Coord is within a certain distance of another Coord.
func (c Coord) IsWithin(other Coord, distance int) bool {
	return c.X >= other.X-distance && c.X <= other.X+distance &&
		c.Y >= other.Y-distance && c.Y <= other.Y+distance
}
