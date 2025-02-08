package camera

import (
	"github.com/hajimehoshi/ebiten/v2"
)

// Viewport represents the view that the player is seeing
type Viewport struct {
	Camera        *Camera
	Width, Height int
}

// NewViewport creates a new viewport
func NewViewport(camera *Camera, width, height int) *Viewport {
	return &Viewport{camera, width, height}
}

// GetTransform returns a GeoM struct with the camera's transform values
func (v *Viewport) GetTransform() ebiten.GeoM {
	var geo ebiten.GeoM
	geo.Translate(-v.Camera.X, -v.Camera.Y)
	geo.Scale(v.Camera.Zoom, v.Camera.Zoom)
	return geo
}
