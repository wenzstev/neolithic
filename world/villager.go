package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Villager represents a single person in the village.
type Villager struct {
	X, Y  int
	Image *ebiten.Image
}

// Draw draws a villager on the screen, on the tile with the matching X Y coordinates
func (v *Villager) Draw(screen *ebiten.Image, transform *ebiten.GeoM, cellSize int) {
	size := v.Image.Bounds().Size().X // assuming villager is square

	worldX := float64(v.X*cellSize + size/2)
	worldY := float64(v.Y*cellSize + size/2)

	var cellTransform ebiten.GeoM
	cellTransform.Reset()
	cellTransform.Translate(worldX, worldY)
	cellTransform.Concat(*transform)

	op := &ebiten.DrawImageOptions{}
	op.GeoM = cellTransform
	screen.DrawImage(v.Image, op)
}

// Move moves a Villager ot a different tile
func (v *Villager) Move(dx, dy int) {
	v.X += dx
	v.Y += dy
}

// GetTile returns the tile the villager is on
func (v *Villager) GetTile(grid Grid) *Tile {
	return grid.Tiles[v.X][v.Y]
}
