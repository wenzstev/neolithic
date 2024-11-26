package camera

// Camera represents the camera
type Camera struct {
	X, Y, Zoom float64
}

// NewCamera initializes a new camera with necessary values
func NewCamera() *Camera {
	return &Camera{
		X:    0,
		Y:    0,
		Zoom: 1.0,
	}
}

// Move moves the camera around
func (c *Camera) Move(deltaX, deltaY float64) {
	c.X += deltaX / c.Zoom
	c.Y += deltaY / c.Zoom
}

// ZoomAt zooms the camera in or out, centered around screenWidth and screenHeight
func (c *Camera) ZoomAt(factor, screenWidth, screenHeight float64) {
	oldZoom := c.Zoom
	c.Zoom *= factor
	if c.Zoom > 5.0 {
		c.Zoom = 5.0
	}
	if c.Zoom < 0.1 {
		c.Zoom = 0.1
	}

	c.X += (screenWidth / 2) * (1/oldZoom - 1/c.Zoom)
	c.Y += (screenHeight / 2) * (1/oldZoom - 1/c.Zoom)
}
