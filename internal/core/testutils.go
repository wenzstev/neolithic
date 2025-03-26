package core

var (
	testLocation = Location{
		Name:      "testlocation1",
		Inventory: NewInventory(),
	}
	testLocation2 = Location{
		Name:      "testlocation2",
		Inventory: NewInventory(),
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
	return nil
}
