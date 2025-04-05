package world

import (
	"image/color"
	"log/slog"

	"Neolithic/internal/agent"
	"Neolithic/internal/camera"
	"Neolithic/internal/core"
	"Neolithic/internal/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

// cellSize is the size of the cells in the world grid
const cellSize = 16

// Engine is the main struct that holds the world state and the images for the villager and location.
type Engine struct {
	World         *core.WorldState
	villagerImage *ebiten.Image
	locationImage *ebiten.Image
	logger        *slog.Logger
}

// NewEngine creates a new Engine.
func NewEngine(grid *grid.Grid, logger *slog.Logger) (*Engine, error) {
	villagerImg := ebiten.NewImage(8, 8)
	villagerImg.Fill(color.RGBA{
		R: 70,
		G: 80,
		B: 100,
		A: 255,
	})

	locationImg := ebiten.NewImage(10, 10)
	locationImg.Fill(color.RGBA{
		R: 170,
		G: 80,
		B: 20,
		A: 255,
	})

	world := &core.WorldState{
		Grid:      grid,
		Locations: map[string]core.Location{},
		Agents:    map[string]core.Agent{},
	}

	return &Engine{
		World:         world,
		villagerImage: villagerImg,
		locationImage: locationImg,
		logger:        logger,
	}, nil
}

// Tick ticks the world state. Iterates through all agents and allows them to run their behavior based on their current state and the world state.
func (e *Engine) Tick(deltaTime float64) error {
	e.logger.Debug("engine tick", "deltaTime", deltaTime)
	for _, a := range e.World.Agents {
		aStruct := a.(*agent.Agent)
		newWorld, err := aStruct.Tick(e.World, deltaTime)
		if err != nil {
			e.logger.Error("agent tick error", "agent", aStruct.Name(), "error", err)
			return err
		}
		if newWorld != nil {
			newWorld.Grid = e.World.Grid
			e.World = newWorld
		}
	}
	return nil
}

// Draw draws the world state on the screen
func (e *Engine) Draw(screen *ebiten.Image, viewport *camera.Viewport, camera *camera.Camera) {
	e.World.Grid.(*grid.Grid).Draw(screen, viewport, camera)
	transform := viewport.GetTransform()

	for _, l := range e.World.Locations {
		DrawEntity(screen, &transform, 16, e.locationImage, l.Coord)
	}
	for _, a := range e.World.Agents {
		DrawEntity(screen, &transform, 16, e.villagerImage, a.(*agent.Agent).Position)
	}
}

// DrawEntity draws an entity on the screen at a given position. Entity can be an agent or a location
func DrawEntity(screen *ebiten.Image, transform *ebiten.GeoM, cellSize int, entityImg *ebiten.Image, position core.Coord) {
	size := entityImg.Bounds().Size().X // assuming villager is square
	worldX := float64(position.X*cellSize + size/2)
	worldY := float64(position.Y*cellSize + size/2)

	var cellTransform ebiten.GeoM
	cellTransform.Reset()
	cellTransform.Translate(worldX, worldY)
	cellTransform.Concat(*transform)

	op := &ebiten.DrawImageOptions{}
	op.GeoM = cellTransform
	screen.DrawImage(entityImg, op)
}
