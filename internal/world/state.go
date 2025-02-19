package world

import (
	"Neolithic/internal/agent"
	"Neolithic/internal/grid"
	"Neolithic/internal/planner"
	"errors"
)

type Coord struct {
	X, Y int
}

func (c Coord) GetTile(grid *grid.Grid) (grid.Tile, error) {
	if len(grid.Tiles) <= c.X {
		return nil, errors.New("no X column for coord")
	}
	if len(grid.Tiles[c.X]) <= c.Y {
		return nil, errors.New("no Y column for coord")
	}

	return grid.Tiles[c.X][c.Y], nil
}

func (c Coord) IsWithin(other *Coord, distance float64) bool {
	return false // todo implement me
}

type Inventory map[*Resource]int

type Location struct {
	name      string
	loc       Coord
	inventory Inventory
}

type Agent struct {
	name      string
	loc       Coord
	inventory Inventory
	Behavior  *Behavior
}

type Behavior struct {
	AvailableActions *[]planner.Action
	Goal             State
	CurState         agent.State
	Plan             *agent.Plan
}

type State struct {
	Agents    []Agent
	Locations []Location
}

type GameWorld struct {
	CurState State
	grid     *grid.Grid
}

func (g *GameWorld) Tick() error {
	for _, agent := range g.CurState.Agents {
		newState, err := agent.curState.Execute(&g.CurState)
		if err != nil {
			return err
		}
		if newState != nil {
			g.CurState = *newState
		}
	}
	return nil
}
