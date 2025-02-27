package agent

import "Neolithic/internal/planner"

var (
	testLocation = &planner.Location{
		Name: "testLocation",
	}

	testResource = &planner.Resource{
		Name: "testResource",
	}
)

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ planner.Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *planner.State, agent planner.Agent) *planner.State {
	return start.Add(m.GetStateChange(agent), false)
}

func (m *mockAction) Cost(_ planner.Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetStateChange(_ planner.Agent) *planner.State {
	return &planner.State{
		Locations: map[*planner.Location]planner.Inventory{
			testLocation: {
				testResource: 1,
			},
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ planner.Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *planner.State, _ planner.Agent) *planner.State {
	return nil
}

func (m *mockNullAction) Cost(_ planner.Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetStateChange(_ planner.Agent) *planner.State {
	return &planner.State{}
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

type mockAgent struct { // todo: remove this? not really necessary...
	name     string
	behavior *Behavior
}

func (m *mockAgent) Name() string {
	return m.name
}

func (m *mockAgent) SetCurState(state State) {
	m.behavior.curState = state
}

func (m *mockAgent) Behavior() *Behavior {
	return m.behavior
}
