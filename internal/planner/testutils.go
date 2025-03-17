package planner

import "Neolithic/internal/core"

var (
	testAgent = &mockAgent{
		N: "testAgent",
	}

	testResource = &core.Resource{
		Name: "testResource",
	}

	gatherTest = &Gather{
		resource: testResource,
		amount:   10,
		locName:  "testLocation",
		cost:     10.0,
	}

	gatherTest2 = &Gather{
		resource: testResource,
		amount:   10,
		locName:  "testLocation2",
		cost:     10.0,
	}

	depositTest = &Deposit{
		resource: testResource,
		amount:   20,
		locName:  "testLocation",
		cost:     1.0,
	}

	depositTest2 = &Deposit{
		resource: testResource,
		amount:   20,
		locName:  "testLocation2",
		cost:     1.0,
	}
)

type mockAgent struct {
	N string
}

func (m *mockAgent) Name() string {
	return m.N
}

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *core.State, agent Agent) *core.State {
	return start.Add(m.GetStateChange(agent), false)
}

func (m *mockAction) Cost(_ Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetStateChange(_ Agent) *core.State {
	return &core.State{
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

func (m *mockNullAction) Perform(_ *core.State, _ Agent) *core.State {
	return nil
}

func (m *mockNullAction) Cost(_ Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetStateChange(_ Agent) *core.State {
	return &core.State{}
}
