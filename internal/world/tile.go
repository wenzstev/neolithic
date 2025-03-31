package world

import (
	"Neolithic/internal/core"
	"errors"
	"fmt"
	"math"

	"Neolithic/internal/astar"
	"Neolithic/internal/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

// Tile implements grid.Tile, and represents a single square of ground
type Tile struct {
	Ground   *Ground
	Resource *Resource
	X, Y     int
	grid     *grid.Grid
}

// Ensure Tile implements grid.Tile
var _ grid.Tile = (*Tile)(nil)

// Draw draws a tile on the screen
func (t *Tile) Draw(screen *ebiten.Image, transform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *transform
	screen.DrawImage(t.Ground.Image, op)
	// TODO: once add resources, draw them here

	if t.Resource != nil {
		screen.DrawImage(t.Resource.Image, op)
	}
}

// Heuristic implements astar.Node and provides a best guess for how far the tile is from the goal tile. Calculates
// based on distance as the crow flies.
func (t *Tile) Heuristic(goal astar.Node) (float64, error) {
	gNode, ok := goal.(*Tile)
	if !ok {
		return 0, fmt.Errorf("heuristic called on non-Tile %T", goal)
	}

	xDistanceSquared := math.Pow(float64(gNode.X-t.X), 2)
	yDistanceSquared := math.Pow(float64(gNode.Y-t.Y), 2)

	return math.Sqrt(xDistanceSquared + yDistanceSquared), nil
}

// ID implements astar.Node and returns a unique string id for the node. It takes the form of "X,Y"
func (t *Tile) ID() (string, error) {
	return fmt.Sprintf("%d,%d", t.X, t.Y), nil
}

// Cost implements astar.Node and returns the cost for moving onto the node. It returns 1 if the node is adjacent and
// 1.4 if the node is diagonal.
func (t *Tile) Cost(prev astar.Node) float64 {
	prevTile := prev.(*Tile)
	if isDiagonallyAdjacent(prevTile, t) {
		return 1.4
	}
	return 1
}

// GetSuccessors implements astar.Node and returns the nodes that are adjacent to the given node.
func (t *Tile) GetSuccessors() ([]astar.Node, error) {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1}, // Top-left, Top, Top-right
		{0, -1}, {0, 1}, // Left,        Right
		{1, -1}, {1, 0}, {1, 1}, // Bottom-left, Bottom, Bottom-right
	}

	var adjacentTiles []astar.Node
	rows := len(t.grid.Tiles)
	cols := len(t.grid.Tiles[0])

	for _, d := range directions {
		newX, newY := t.X+d.dx, t.Y+d.dy
		if newX >= 0 && newX < rows && newY >= 0 && newY < cols {
			adjacentTile, ok := t.grid.Tiles[newX][newY].(*Tile)
			if !ok {
				return nil, errors.New("grid tile is not instance of ground tile") // todo better error
			}
			adjacentTiles = append(adjacentTiles, adjacentTile)
		}
	}

	return adjacentTiles, nil
}

// isDiagonallyAdjacent is a helper function that determines whether two tiles are diagonally adjacent to each other.
func isDiagonallyAdjacent(tile1, tile2 *Tile) bool {
	return math.Abs(float64(tile1.X-tile2.X)) == 1 && math.Abs(float64(tile1.Y-tile2.Y)) == 1
}

func (t *Tile) Coord() core.Coord {
	return core.Coord{X: t.X, Y: t.Y}
}
