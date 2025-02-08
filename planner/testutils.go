package planner

var (
	testLocation = &Location{
		Name: "testLocation",
	}

	testAgent = &Agent{
		Name: "testAgent",
	}

	testResource = &Resource{
		Name: "testResource",
	}
)

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ Action = (*mockAction)(nil)

func (m *mockAction) Perform(state *State, _ *Agent) *State {
	endState := state.Copy()
	_, ok := endState.Locations[testLocation]
	if !ok {
		endState.Locations[testLocation] = Inventory{}
	}

	endState.Locations[testLocation][testResource]++
	return endState
}

func (m *mockAction) Cost(_ *Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock action"
}

// mockNullAction implements Action and is used for testing. It always returns a null state.
type mockNullAction struct{}

var _ Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *State, _ *Agent) *State {
	return nil
}

func (m *mockNullAction) Cost(_ *Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null action"
}
