package world

import (
	"Neolithic/internal/core"
	"Neolithic/internal/planner"
)

// mockAction implements planner.Action for testing
type mockAction struct {
	testVal string
}

func (m *mockAction) Perform(world *core.WorldState, agent core.Agent) *core.WorldState {
	return world.DeepCopy()
}

func (m *mockAction) Cost(_ core.Agent) float64 {
	return 1.0
}

func (m *mockAction) Description() string {
	return "mock action"
}

func (m *mockAction) GetChanges(_ core.Agent) []planner.StateChange {
	return nil
}
