package core

var (
	testLocation = Location{
		Name: "testlocation1",
	}
	testLocation2 = Location{
		Name: "testlocation2",
	}
	testAgent = mockAgent{}
)

type mockAgent struct{}

func (m mockAgent) String() string {
	return "mockagentstring"
}

func (m mockAgent) Name() string {
	return "mockagent"
}

func (m mockAgent) DeepCopy() Agent {
	return m
}

func (m mockAgent) Inventory() Inventory {
	return &mockInventory{}
}

type mockInventory struct{}

func (m mockInventory) GetAmount(res *Resource) int {
	return 5
}

func (m mockInventory) AdjustAmount(res *Resource, amount int) {}

func (m mockInventory) DeepCopy() Inventory {
	return m
}

func (m mockInventory) String() string {
	return "mockInventorystring"
}
