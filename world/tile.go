package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Ground   Ground
	CellSize int
}

func (t *Tile) Draw(screen *ebiten.Image, transform *ebiten.GeoM) {
	t.Ground.Draw(screen, transform)
}
