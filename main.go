package main

import (
	"Neolithic/camera"
	"Neolithic/world"
	"github.com/hajimehoshi/ebiten/v2"
	"image/color"
	"log"
)

type Game struct {
	World    *world.World
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

	if ebiten.IsKeyPressed(ebiten.KeyT) {
		g.World.Villagers[0].Move(1, 1)
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.World.Draw(screen, g.Viewport, g.Camera)
}

func (g *Game) Layout(_, _ int) (screenWidth, screenHeight int) {
	return g.Viewport.Width, g.Viewport.Height
}

func main() {
	cam := camera.NewCamera()
	vp := camera.NewViewport(cam, 800, 600)
	width, height := 32, 32
	cellSize := 16

	game := &Game{
		World:    world.New(width, height, cellSize),
		Camera:   cam,
		Viewport: vp,
	}

	villagerImg := ebiten.NewImage(8, 8)
	villagerImg.Fill(color.RGBA{
		R: 70,
		G: 80,
		B: 100,
		A: 255,
	})

	game.World.Villagers = append(game.World.Villagers, &world.Villager{
		X:     0,
		Y:     1,
		Image: villagerImg,
	})

	ebiten.SetWindowSize(800, 600)
	ebiten.SetWindowTitle("Hello, World!")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
