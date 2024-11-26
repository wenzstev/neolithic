package world

import (
	"Neolithic/drawable"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
)

type Ground interface {
	drawable.Drawable
}

type RGBGround struct {
	Color color.RGBA
	Image *ebiten.Image
}

func NewRGBGround(col color.RGBA, cellSize int) *RGBGround {
	image := ebiten.NewImage(cellSize, cellSize)
	image.Fill(col)
	return &RGBGround{
		Color: col,
		Image: image,
	}
}

func (r *RGBGround) Draw(screen *ebiten.Image, transform *ebiten.GeoM) {
	op := &ebiten.DrawImageOptions{}
	op.GeoM = *transform
	screen.DrawImage(r.Image, op)
}

var _ Ground = &RGBGround{}

var GreenGround = NewRGBGround(color.RGBA{R: 0, G: 255, B: 0, A: 255}, 32)
var BrownGround = NewRGBGround(color.RGBA{150, 75, 0, 255}, 32)
