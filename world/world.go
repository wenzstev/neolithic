package world

import (
	"Neolithic/camera"
	"Neolithic/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

// Grid aliases to grid.Grid[*Tile]
type Grid = grid.Grid[*Tile]

// World represents the game world
type World struct {
	Villagers []*Villager
	Grid      *Grid
}

// New creates a new instance of World
func New(width, height, cellSize int) *World {
	world := &World{
		Villagers: make([]*Villager, 0),
		Grid:      grid.New[*Tile](width, height, cellSize),
	}

	world.Grid.Initialize(makeTile)
	return world
}

// makeTile returns a new Grass tile to populate the world grid
func makeTile() (*Tile, error) {
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
