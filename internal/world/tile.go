package world

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Tile implements grid.Tile, and represents a single square of ground
type Tile struct {
	Ground   *Ground
	Resource *Resource
}

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
