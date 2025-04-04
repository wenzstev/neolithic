package grid

import (
	"fmt"
	"math"

	"Neolithic/internal/camera"
	"Neolithic/internal/core"
	"github.com/hajimehoshi/ebiten/v2"
)

// Grid represents the map, divided into Width by Height tiles.
type Grid struct {
	Width    int
	Height   int
	CellSize int
	Tiles    [][]Tile
}

// Tile represents a single square in the Grid
type Tile interface {
	core.Cell
	Draw(*ebiten.Image, *ebiten.GeoM)
}

// New creates a new instance of Grid
func New(width, height, cellSize int) (*Grid, error) {
	if width <= 0 || height <= 0 {
		return nil, fmt.Errorf("width and height must be positive")
	}
	grid := &Grid{
		Width:    width,
		Height:   height,
		CellSize: cellSize,
	}
	return grid, nil
}

// Initialize initializes a new grid with necessary values. Takes in a function
// for making a tile
func (g *Grid) Initialize(MakeTile func(X, Y int, grid *Grid) (Tile, error)) error {
	g.Tiles = make([][]Tile, g.Width)
	for i := 0; i < g.Width; i++ {
		g.Tiles[i] = make([]Tile, g.Height)
		for j := 0; j < g.Height; j++ {
			tile, err := MakeTile(i, j, g)
			if err != nil {
				return err
			}

			g.Tiles[i][j] = tile

		}
	}
	return nil
}

// DrawCell draws a grid cell
func (g *Grid) drawCell(screen *ebiten.Image, x, y int, transform *ebiten.GeoM) {
	worldX := float64(x * g.CellSize)
	worldY := float64(y * g.CellSize)

	var cellTransform ebiten.GeoM

	cellTransform.Reset()
	cellTransform.Translate(worldX, worldY)
	cellTransform.Concat(*transform)

	g.Tiles[x][y].Draw(screen, &cellTransform)
}

// Draw draws the Grid, based on the viewport and camera location
func (g *Grid) Draw(screen *ebiten.Image, viewport *camera.Viewport, camera *camera.Camera) {
	transform := viewport.GetTransform()

	screenWidth, screenHeight := viewport.Width, viewport.Height
	cellSize := float64(g.CellSize)

	invZoom := 1.0 / camera.Zoom

	leftWorld := camera.X
	rightWorld := camera.X + float64(screenWidth)*invZoom
	topWorld := camera.Y
	bottomWorld := camera.Y + float64(screenHeight)*invZoom

	left := int(math.Floor(leftWorld / cellSize))
	right := int(math.Ceil(rightWorld / cellSize))
	top := int(math.Floor(topWorld / cellSize))
	bottom := int(math.Ceil(bottomWorld / cellSize))

	// Clamp indices to grid bounds
	if left < 0 {
		left = 0
	}
	if right > g.Width {
		right = g.Width
	}
	if top < 0 {
		top = 0
	}
	if bottom > g.Height {
		bottom = g.Height
	}

	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			if x >= 0 && x < g.Width && y >= 0 && y < g.Height {
				g.drawCell(screen, x, y, &transform)
			}
		}
	}
}

// CellAt returns the cell at the given coordinate
func (g *Grid) CellAt(coord core.Coord) core.Cell {
	x := coord.X
	y := coord.Y

	if x < 0 || x >= g.Width || y < 0 || y >= g.Height {
		return nil
	}
	return g.Tiles[x][y]
}
