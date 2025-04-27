package world

import (
	"errors"
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

var (
	// ErrAgentAlreadyExists is thrown when an agent with a duplicate name is added to the world.
	ErrAgentAlreadyExists = errors.New("agent already exists")
	// ErrLocationAlreadyExists is thrown when a location with a duplicate name is added to the world.
	ErrLocationAlreadyExists = errors.New("location already exists")
)

// Engine is the main struct that holds the world state and the images for the villager and location.
type Engine struct {
	// World is the main world state
	World *core.WorldState
	// Registry holds all actions, resources, and locations and creates actions when new resources and locations are provided.
	Registry *Registry
	// villagerImage is the sprite used to represent a villager
	villagerImage *ebiten.Image
	// locationImage is the sprite used to represent a location
	locationImage *ebiten.Image
	// logger is the logger
	logger *slog.Logger
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
		Locations: []core.Location{},
		Agents:    []core.Agent{},
	}

	return &Engine{
		World: world,
		Registry: &Registry{
			ActionRegistry: []*ActionRegistryEntry{},
			Actions:        []core.Action{},
			Locations:      []*core.Location{},
			Resources:      []*core.Resource{},
		},
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

// AddLocation adds a new location to the world and registers it in the registry. It returns an error if registration fails.
func (e *Engine) AddLocation(location *core.Location) error {
	_, exists := e.World.GetLocation(location.Name)
	if exists {
		return ErrLocationAlreadyExists
	}
	if err := e.Registry.RegisterLocation(location); err != nil {
		return err
	}
	e.World.Locations = append(e.World.Locations, *location)
	return nil
}

// AddResource registers a resource in the registry of the engine and returns an error if the operation fails.
func (e *Engine) AddResource(resource *core.Resource) error {
	return e.Registry.RegisterResource(resource)
}

// RegisterAction registers a new action in the engine's registry with a specified name, action, and creation function.
func (e *Engine) RegisterAction(name string, action core.Action, createFunc ActionCreator) error {
	return e.Registry.RegisterAction(name, action, createFunc)
}

// AddAgent adds a new agent to the world and updates its possible actions. Returns an error if the agent already exists.
func (e *Engine) AddAgent(agent *agent.Agent) error {
	_, exists := e.World.GetAgent(agent.Name())
	if exists {
		return ErrAgentAlreadyExists
	}

	agent.Behavior.PossibleActions = e.Registry.Actions
	e.World.Agents = append(e.World.Agents, agent)
	return nil
}
