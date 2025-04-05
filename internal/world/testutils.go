package world

import (
	"Neolithic/internal/core"
	"Neolithic/internal/planner"
)

var (
	mockActionCreateFunc = func(params ActionCreatorParams) planner.Action {
		return &mockAction{}
	}

	mockActionWithLocationCreateFunc = func(params ActionCreatorParams) planner.Action {
		return &mockActionWithLocation{
			location: params.Location,
		}
	}

	mockActionWithResourceCreateFunc = func(params ActionCreatorParams) planner.Action {
		return &mockActionWithResource{resource: params.Resource}
	}

	mockActionWithResourceAndLocationCreateFunc = func(params ActionCreatorParams) planner.Action {
		return &mockActionWithLocationAndResource{
			resource: params.Resource,
			location: params.Location,
		}
	}
)

// mockAction implements planner.Action for testing
type mockAction struct {
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

type mockActionWithLocation struct {
	mockAction
	location *core.Location
}

func (m *mockActionWithLocation) Location() *core.Location {
	return m.location
}

type mockActionWithResource struct {
	mockAction
	resource *core.Resource
}

func (m *mockActionWithResource) Resource() *core.Resource {
	return m.resource
}

type mockActionWithLocationAndResource struct {
	mockAction
	resource *core.Resource
	location *core.Location
}

func (m *mockActionWithLocationAndResource) Location() *core.Location {
	return m.location
}

func (m *mockActionWithLocationAndResource) Resource() *core.Resource {
	return m.resource
}
