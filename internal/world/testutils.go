package world

import (
	"Neolithic/internal/core"
)

var (
	mockActionCreateFunc = func(params ActionCreatorParams) core.Action {
		return &mockAction{}
	}

	mockActionWithLocationCreateFunc = func(params ActionCreatorParams) core.Action {
		return &mockActionWithLocation{
			location: params.Location,
		}
	}

	mockActionWithResourceCreateFunc = func(params ActionCreatorParams) core.Action {
		return &mockActionWithResource{resource: params.Resource}
	}

	mockActionWithResourceAndLocationCreateFunc = func(params core.CreateActionParams) core.Action {
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

func (m *mockAction) GetChanges(_ core.Agent) []core.StateChange {
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

type mockAttribute struct {
	requiresLocation bool
	requiresResource bool
	action           core.Action
}

func (m *mockAttribute) NeedsLocation() bool { return m.requiresLocation }
func (m *mockAttribute) NeedsResource() bool { return m.requiresResource }
func (m *mockAttribute) CreateAction(params core.CreateActionParams) core.Action {
	return mockActionWithResourceAndLocationCreateFunc(params)
}

func createTestAttribute(needsLoc, needsRes bool, action core.Action) core.Attribute {
	return &mockAttribute{
		requiresLocation: needsLoc,
		requiresResource: needsRes,
		action:           action,
	}
}

var (
	mockAttributeNoNeeds   = createTestAttribute(false, false, &mockAction{})
	mockAttributeNeedsLoc  = createTestAttribute(true, false, &mockAction{})
	mockAttributeNeedsRes  = createTestAttribute(false, true, &mockAction{})
	mockAttributeNeedsBoth = createTestAttribute(true, true, &mockAction{})
)
