package world

import "github.com/hajimehoshi/ebiten/v2"

type Resource struct {
	Name   string
	Image  *ebiten.Image
	Amount int
}
