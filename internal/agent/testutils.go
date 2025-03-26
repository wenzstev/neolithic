package agent

import (
	"Neolithic/internal/core"
	"Neolithic/internal/planner"
)

var (
	testResource = &core.Resource{
		Name: "testResource",
	}
)

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ planner.Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()
	end.Locations["testLocation"].Inventory.AdjustAmount(testResource, 1)
	return end
}

func (m *mockAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetChanges(agent core.Agent) []planner.StateChange {
	return []planner.StateChange{
		{
			EntityType: planner.LocationEntity,
			Entity:     "testLocation",
			Resource:   testResource,
			Amount:     1,
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ planner.Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *core.WorldState, _ core.Agent) *core.WorldState {
	return nil
}

func (m *mockNullAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetChanges(agent core.Agent) []planner.StateChange {
	return []planner.StateChange{}
}

type mockActionWithTime struct {
	mockAction
	timeNeeded float64
}

func (m *mockActionWithTime) TimeNeeded() float64 {
	return m.timeNeeded
}

type mockPlan struct {
	isComplete bool
	nextAction planner.Action
}

func (m *mockPlan) IsComplete() bool {
	return m.isComplete
}

func (m *mockPlan) PeekAction() planner.Action {
	return m.nextAction
}

func (m *mockPlan) PopAction() planner.Action {
	return m.nextAction
}
