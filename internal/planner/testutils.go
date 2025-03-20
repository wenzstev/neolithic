package planner

import (
	"Neolithic/internal/core"
	"sort"
)

var (
	testLocation  = core.Location{Name: "testLocation"}
	testLocation2 = core.Location{Name: "testLocation2"}

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

type mockInventory struct {
	entries []core.InventoryEntry
}

func (m mockInventory) String() string {
	return "mockInventory"
}

func (m mockInventory) GetAmount(res *core.Resource) int {
	for _, entry := range m.entries {
		if entry.Resource == res {
			return entry.Amount
		}
	}
	return 0
}

func (m mockInventory) AdjustAmount(res *core.Resource, amount int) {
	for i := 0; i < len(m.entries); i++ {
		if m.entries[i].Resource == res {
			m.entries[i].Amount += amount
			return
		}
	}
	m.entries = append(m.entries, core.InventoryEntry{Resource: res, Amount: amount})
	sort.Slice(m.entries, func(i, j int) bool {
		return m.entries[i].Resource.Name < m.entries[j].Resource.Name
	})
}

func (m mockInventory) DeepCopy() core.Inventory {
	newEntries := make([]core.InventoryEntry, len(m.entries))
	copy(newEntries, m.entries)
	return &mockInventory{
		entries: newEntries,
	}
}

func (m mockInventory) Entries() []core.InventoryEntry {
	return m.entries
}

type mockAgent struct {
	N string
}

func (m *mockAgent) String() string {
	return "mockAgent"
}

func (m *mockAgent) DeepCopy() core.Agent {
	return &mockAgent{}
}

func (m *mockAgent) Inventory() core.Inventory {
	//TODO implement me
	panic("implement me")
}

func (m *mockAgent) Name() string {
	return m.N
}

// mockAction implements Action and is used for testing.
type mockAction struct{}

var _ Action = (*mockAction)(nil)

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

func (m *mockAction) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{
		{
			Entity:     "testLocation",
			EntityType: LocationEntity,
			Resource:   testResource,
			Amount:     1,
		},
	}
}

// mockNullAction implements Action and is used for testing. It always returns a null State.
type mockNullAction struct{}

var _ Action = (*mockNullAction)(nil)

func (m *mockNullAction) Perform(_ *core.WorldState, _ core.Agent) *core.WorldState {
	return nil
}

func (m *mockNullAction) Cost(_ core.Agent) float64 {
	return 10.0
}

func (m *mockNullAction) Description() string {
	return "a mock null Action"
}

func (m *mockNullAction) GetChanges(agent core.Agent) []StateChange {
	return []StateChange{}
}
