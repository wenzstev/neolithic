package core

type WorldState struct {
	Locations map[string]Location
	Agents    map[string]Agent
}

func (*WorldState) ID() (string, error) {
	panic("implement me")
}

func (w *WorldState) Copy() *WorldState {
	panic("implement me")
}

type Inventory map[*Resource]int

type Resource struct {
	Name string
}

type Location struct {
	Name      string
	Inventory Inventory
	Coord     Coord
}

type Agent interface {
	Name() string
	AdjustInventory(*Resource, int) Agent
	GetAmount(*Resource) int
}

type Coord struct {
	X, Y int
}
