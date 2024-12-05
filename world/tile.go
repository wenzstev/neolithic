package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

type Tile struct {
	Ground   *Ground
	Resource *Resource
}

func (t *Tile) Draw(screen *ebiten.Image, transform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *transform
	screen.DrawImage(t.Ground.Image, op)
	// TODO: once add resources, draw them here

	if t.Resource != nil {
		screen.DrawImage(t.Resource.Image, op)
	}
}
