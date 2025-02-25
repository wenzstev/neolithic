package agent

import (
	"Neolithic/internal/grid"
	"Neolithic/internal/world"
)

type Agent struct {
	name     string
	grid     *grid.Grid
	behavior *Behavior
	loc      *world.Coord
}

func (a *Agent) Name() string {
	return a.name
}

func (a *Agent) SetLocation(loc *world.Coord) {
	a.loc = loc
}
