package planner

import (
	"Neolithic/internal/core"
	"encoding/gob"
)

var (
	testLocation  = core.Location{Name: "testLocation", Inventory: core.NewInventory()}
	testLocation2 = core.Location{Name: "testLocation2", Inventory: core.NewInventory()}

	testAgent = &mockAgent{
		N:         "testAgent",
		inventory: core.NewInventory(),
	}

	testResource = &core.Resource{
		Name: "testResource",
	}

	gatherTest = &Gather{
		Res:            testResource,
		Amount:         10,
		ActionLocation: &core.Location{Name: "testLocation"},
		ActionCost:     10.0,
	}

	gatherTest2 = &Gather{
		Res:            testResource,
		Amount:         10,
		ActionLocation: &core.Location{Name: "testLocation2"},
		ActionCost:     10.0,
	}

	depositTest = &Deposit{
		DepResource:    testResource,
		Amount:         20,
		ActionLocation: &core.Location{Name: "testLocation"},
		ActionCost:     1.0,
	}

	depositTest2 = &Deposit{
		DepResource:    testResource,
		Amount:         20,
		ActionLocation: &core.Location{Name: "testLocation2"},
		ActionCost:     1.0,
	}
)

func init() {
	gob.Register(mockAgent{})
}

type mockAgent struct {
	N         string
	inventory core.Inventory
}

func (m *mockAgent) String() string {
	return "mockAgent"
}

func (m *mockAgent) DeepCopy() core.Agent {
	return &mockAgent{
		N:         m.N,
		inventory: m.inventory.DeepCopy(),
	}
}

func (m *mockAgent) Inventory() core.Inventory {
	return m.inventory
}

func (m *mockAgent) Name() string {
	return m.N
}

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ core.Action = (*mockAction)(nil)

func (m *mockAction) Perform(start *core.WorldState, agent core.Agent) *core.WorldState {
	end := start.DeepCopy()
	endLocation, _ := end.GetLocation("testLocation")
	endLocation.Inventory.AdjustAmount(testResource, 1)
	return end
}

func (m *mockAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockAction) Description() string {
	return "a mock Action"
}

func (m *mockAction) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{
		{
			Entity:     "testLocation",
			EntityType: core.LocationEntity,
			Resource:   testResource,
			Amount:     1,
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ core.Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *core.WorldState, _ core.Agent) *core.WorldState {
	return nil
}

func (m *mockNullAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetChanges(agent core.Agent) []core.StateChange {
	return []core.StateChange{}
}
