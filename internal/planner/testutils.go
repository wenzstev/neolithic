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

	testLocation2 = &Location{Name: "testLocation2"}

	gatherTest = &Gather{
		resource: testResource,
		amount:   10,
		location: testLocation,
		cost:     10.0,
	}

	gatherTest2 = &Gather{
		resource: testResource,
		amount:   10,
		location: testLocation2,
		cost:     10.0,
	}

	depositTest = &Deposit{
		resource: testResource,
		amount:   20,
		location: testLocation,
		cost:     1.0,
	}

	depositTest2 = &Deposit{
		resource: testResource,
		amount:   20,
		location: testLocation2,
		cost:     1.0,
	}
)

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *State, agent *Agent) *State {
	return start.Add(m.GetStateChange(agent), false)
}

func (m *mockAction) Cost(_ *Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetStateChange(_ *Agent) *State {
	return &State{
		Locations: map[*Location]Inventory{
			testLocation: {
				testResource: 1,
			},
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *State, _ *Agent) *State {
	return nil
}

func (m *mockNullAction) Cost(_ *Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetStateChange(_ *Agent) *State {
	return &State{}
}
