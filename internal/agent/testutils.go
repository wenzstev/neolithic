package agent

import (
	"Neolithic/internal/astar"
	"Neolithic/internal/core"
	"fmt"
)

var (
	testResource = &core.Resource{
		Name: "testResource",
	}
)

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ core.Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()
	endLoc, _ := end.GetLocation("testLocation")
	endLoc.Inventory.AdjustAmount(testResource, 1)
	return end
}

func (m *mockAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{
		{
			EntityType: core.LocationEntity,
			Entity:     "testLocation",
			Resource:   testResource,
			Amount:     1,
		},
	}
}

type mockNilAction struct {
	mockAction
}

func (m *mockNilAction) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	return nil
}

type mockLocationAction struct {
	mockAction
	location *core.Location
}

func (m *mockLocationAction) Location() *core.Location {
	return &core.Location{
		Name: "testLocation",
		Coord: core.Coord{
			X: 3,
			Y: 3,
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ core.Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *core.WorldState, _ core.Agent) *core.WorldState {
	return nil
}

func (m *mockNullAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{}
}

type mockActionWithTime struct {
	mockAction
	timeNeeded float64
}

func (m *mockActionWithTime) TimeNeeded() float64 {
	return m.timeNeeded
}

type MockPlan struct {
	Complete   bool
	NextAction core.Action
}

func (m *MockPlan) IsComplete() bool {
	return m.Complete
}

func (m *MockPlan) PeekAction() core.Action {
	return m.NextAction
}

func (m *MockPlan) PopAction() core.Action {
	return m.NextAction
}

type mockGrid struct{}

func (m mockGrid) CellAt(coord core.Coord) core.Cell {
	return &mockTile{coord.X, coord.Y}
}

type mockTile struct {
	X, Y int
}

func (m *mockTile) Heuristic(goal astar.Node) (float64, error) {
	if m == goal.(*mockTile) {
		return 0.0, nil
	}
	return 1.0, nil
}

func (m *mockTile) ID() (string, error) {
	return fmt.Sprintf("%d,%d", m.X, m.Y), nil

}

func (m *mockTile) Cost(prev astar.Node) float64 {
	return 1
}

func (m *mockTile) GetSuccessors() ([]astar.Node, error) {
	directions := []struct{ dx, dy int }{
		{-1, -1}, {-1, 0}, {-1, 1}, // Top-left, Top, Top-right
		{0, -1}, {0, 1}, // Left,        Right
		{1, -1}, {1, 0}, {1, 1}, // Bottom-left, Bottom, Bottom-right
	}

	var adjacentTiles []astar.Node

	for _, d := range directions {
		newX, newY := m.X+d.dx, m.Y+d.dy
		adjacentTiles = append(adjacentTiles, &mockTile{newX, newY})
	}

	return adjacentTiles, nil
}

func (m *mockTile) Coord() core.Coord {
	return core.Coord{X: m.X, Y: m.Y}
}
