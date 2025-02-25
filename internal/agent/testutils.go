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
