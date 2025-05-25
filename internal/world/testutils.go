package world

import (
	"Neolithic/internal/core"
)

var (
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
	attrType         core.AttributeType
}

func (m *mockAttribute) NeedsLocation() bool { return m.requiresLocation }
func (m *mockAttribute) NeedsResource() bool { return m.requiresResource }
func (m *mockAttribute) CreateAction(holder core.AttributeHolder, params core.CreateActionParams) (core.Action, error) {
	return mockActionWithResourceAndLocationCreateFunc(params), nil
}
func (m *mockAttribute) String() string {
	return "mockAttribute"
}
func (m *mockAttribute) Type() core.AttributeType {
	return m.attrType
}
func (m *mockAttribute) Copy() core.Attribute {
	return &mockAttribute{
		requiresLocation: m.requiresLocation,
		requiresResource: m.requiresResource,
		action:           m.action,
	}
}

func createTestAttribute(needsLoc, needsRes bool, action core.Action, attrType core.AttributeType) core.Attribute {
	return &mockAttribute{
		requiresLocation: needsLoc,
		requiresResource: needsRes,
		action:           action,
		attrType:         attrType,
	}
}

var (
	mockAttributeNoNeeds   = createTestAttribute(false, false, &mockAction{}, "noNeeds")
	mockAttributeNeedsLoc  = createTestAttribute(true, false, &mockAction{}, "needsLoc")
	mockAttributeNeedsRes  = createTestAttribute(false, true, &mockAction{}, "needsRes")
	mockAttributeNeedsBoth = createTestAttribute(true, true, &mockAction{}, "needsBoth")
)
