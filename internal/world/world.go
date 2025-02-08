package world

import (
	"Neolithic/internal/camera"
	"Neolithic/internal/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

// World represents the game world
type World struct {
	Villagers []*Villager
	Grid      *grid.Grid
}

// New creates a new instance of World
func New(width, height, cellSize int) *World {
	world := &World{
		Villagers: make([]*Villager, 0),
		Grid:      grid.New(width, height, cellSize),
	}

	world.Grid.Initialize(makeTile)
	return world
}

// makeTile returns a new Grass tile to populate the world grid
func makeTile() (grid.Tile, error) {
	ground, err := NewGrassGround()
	if err != nil {
		return nil, err
	}

	return &Tile{
		Ground: ground,
	}, nil
}

// Draw draws the world
func (w *World) Draw(screen *ebiten.Image, viewport *camera.Viewport, camera *camera.Camera) {
	w.Grid.Draw(screen, viewport, camera)

	transform := viewport.GetTransform()
	for _, villager := range w.Villagers {
		villager.Draw(screen, &transform, w.Grid.CellSize)
	}

}
