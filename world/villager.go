package world

import (
	"Neolithic/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Villager struct {
	X, Y  int
	Image *ebiten.Image
}

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

func (v *Villager) Move(dx, dy int) {
	v.X += dx
	v.Y += dy
}

func (v *Villager) GetTile(grid grid.Grid) grid.Tile {
	return grid.Tiles[v.X][v.Y]
}
