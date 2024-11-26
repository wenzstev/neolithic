package grid

import (
	"Neolithic/drawable"
	"Neolithic/world"
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
)

// Grid represents the map
type Grid struct {
	Width    int
	Height   int
	CellSize int
	Tiles    [][]drawable.Drawable
}

// Initialize initializes a new grid with necessary values
func (g *Grid) Initialize() {
	g.Tiles = make([][]drawable.Drawable, g.Width)
	for i := 0; i < g.Width; i++ {
		g.Tiles[i] = make([]drawable.Drawable, g.Height)
		for j := 0; j < g.Height; j++ {
			ground, err := world.NewGrassGround()
			if err != nil {
				fmt.Printf(err.Error())
				continue
			}
			g.Tiles[i][j] = &world.Tile{
				Ground: ground,
			}
		}
	}
}

// DrawCell draws a grid cell
func (g *Grid) DrawCell(screen *ebiten.Image, x, y int, transform *ebiten.GeoM) {
	worldX := float64(x * g.CellSize)
	worldY := float64(y * g.CellSize)

	var cellTransform ebiten.GeoM

	cellTransform.Reset()
	cellTransform.Translate(worldX, worldY)
	cellTransform.Concat(*transform)

	g.Tiles[x][y].Draw(screen, &cellTransform)
}
