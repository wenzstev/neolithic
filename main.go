package main

import (
	"log"
	"math"

	"Neolithic/camera"
	"Neolithic/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	Grid     grid.Grid
	Camera   *camera.Camera
	Viewport *camera.Viewport
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyW) {
		g.Camera.Move(0, -10)
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) {
		g.Camera.Move(0, 10)
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) {
		g.Camera.Move(10, 0)
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) {
		g.Camera.Move(-10, 0)
	}

	if ebiten.IsKeyPressed(ebiten.KeyQ) {
		g.Camera.ZoomAt(1.05, float64(g.Viewport.Width), float64(g.Viewport.Height))
	}
	if ebiten.IsKeyPressed(ebiten.KeyE) {
		g.Camera.ZoomAt(.95, float64(g.Viewport.Width), float64(g.Viewport.Height))
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	transform := g.Viewport.GetTransform()

	screenWidth, screenHeight := g.Viewport.Width, g.Viewport.Height
	cellSize := float64(g.Grid.CellSize)

	invZoom := 1.0 / g.Camera.Zoom

	leftWorld := g.Camera.X
	rightWorld := g.Camera.X + float64(screenWidth)*invZoom
	topWorld := g.Camera.Y
	bottomWorld := g.Camera.Y + float64(screenHeight)*invZoom

	left := int(math.Floor(leftWorld / cellSize))
	right := int(math.Ceil(rightWorld / cellSize))
	top := int(math.Floor(topWorld / cellSize))
	bottom := int(math.Ceil(bottomWorld / cellSize))

	// Clamp indices to grid bounds
	if left < 0 {
		left = 0
	}
	if right > g.Grid.Width {
		right = g.Grid.Width
	}
	if top < 0 {
		top = 0
	}
	if bottom > g.Grid.Height {
		bottom = g.Grid.Height
	}

	for y := top; y < bottom; y++ {
		for x := left; x < right; x++ {
			if x >= 0 && x < g.Grid.Width && y >= 0 && y < g.Grid.Height {
				g.Grid.DrawCell(screen, x, y, &transform)
			}
		}
	}

}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.Viewport.Width, g.Viewport.Height
}

func main() {
	cam := camera.NewCamera()
	vp := camera.NewViewport(cam, 800, 600)

	game := &Game{
		Grid: grid.Grid{
			Width:    32,
			Height:   32,
			CellSize: 16,
		},
		Camera:   cam,
		Viewport: vp,
	}

	game.Grid.Initialize()

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
