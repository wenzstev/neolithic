package world

import "github.com/hajimehoshi/ebiten/v2"

// Resource represents a resource in the world, such as flint or clay.
type Resource struct {
	Name   string
	Image  *ebiten.Image
	Amount int
}
