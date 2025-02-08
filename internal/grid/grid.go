package grid

import (
	"Neolithic/internal/camera"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"math"
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
	Draw(*ebiten.Image, *ebiten.GeoM)
}

// New creates a new instance of Grid
func New(width, height, cellSize int) *Grid {
	grid := &Grid{
		Width:    width,
		Height:   height,
		CellSize: cellSize,
	}
	return grid
}

// Initialize initializes a new grid with necessary values. Takes in a function
// for making a tile
func (g *Grid) Initialize(MakeTile func() (Tile, error)) {
	g.Tiles = make([][]Tile, g.Width)
	for i := 0; i < g.Width; i++ {
		g.Tiles[i] = make([]Tile, g.Height)
		for j := 0; j < g.Height; j++ {
			tile, err := MakeTile()
			if err != nil {
				fmt.Printf("%s", err.Error())
				continue
			}

			g.Tiles[i][j] = tile

		}
	}
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
