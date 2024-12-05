package world

import (
	"Neolithic/camera"
	"Neolithic/grid"
	"github.com/hajimehoshi/ebiten/v2"
)

type World struct {
	Villagers []*Villager
	Grid      *grid.Grid
}

func New(width, height, cellSize int) *World {
	world := &World{
		Villagers: make([]*Villager, 0),
		Grid:      grid.New(width, height, cellSize),
	}

	world.Grid.Initialize(makeTile)
	return world
}

func makeTile() (grid.Tile, error) {
	ground, err := NewGrassGround()
	if err != nil {
		return nil, err
	}

	return &Tile{
		Ground: ground,
	}, nil
}

func (w *World) Draw(screen *ebiten.Image, viewport *camera.Viewport, camera *camera.Camera) {
	w.Grid.Draw(screen, viewport, camera)

	transform := viewport.GetTransform()
	for _, villager := range w.Villagers {
		villager.Draw(screen, &transform, w.Grid.CellSize)
	}

}
