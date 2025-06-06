package attributes

import (
	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDeposit = &Deposit{
	DepResource:    testResource,
	Amount:         10,
	ActionLocation: &core.Location{Name: "testLocation", Inventory: core.NewInventory()},
	ActionCost:     1.0,
}

func TestDeposit_Perform(t *testing.T) {
	type testCase struct {
		testDeposit              *Deposit
		startLocation            *core.Location
		startAmountInLocation    int
		agent                    core.Agent
		startAmountInAgent       int
		expectedAmountInLocation int
		expectedAmountInAgent    int
		expectNil                bool
	}

	tests := map[string]testCase{
		"can do basic deposit": {
			testDeposit:              testDeposit,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    0,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       10,
			expectedAmountInLocation: 10,
			expectedAmountInAgent:    0,
		},
		"deposit fails, nothing in agent inventory": {
			testDeposit:              testDeposit,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    0,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    0,
			expectNil:                true,
		},
		"partial deposit success": {
			testDeposit:              testDeposit,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    0,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       5,
			expectedAmountInLocation: 5,
			expectedAmountInAgent:    0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.startLocation.Inventory.AdjustAmount(testResource, tc.startAmountInLocation)
			tc.agent.Inventory().AdjustAmount(testResource, tc.startAmountInAgent)
			startState := &core.WorldState{
				Locations: map[string]*core.Location{tc.startLocation.Name: tc.startLocation},
				Agents:    map[string]core.Agent{tc.agent.Name(): tc.agent},
			}

			endState := tc.testDeposit.Perform(startState, tc.agent)
			if tc.expectNil {
				assert.Nil(t, endState)
				return
			}
			endAgent, exists := endState.GetAgent("testAgent")
			assert.True(t, exists)
			endLocation, exists := endState.GetLocation("testLocation")
			assert.True(t, exists)
			assert.Equal(t, tc.expectedAmountInAgent, endAgent.Inventory().GetAmount(testResource))
			assert.Equal(t, tc.expectedAmountInLocation, endLocation.Inventory.GetAmount(testResource))
		})
	}
}

func TestDeposit_Cost(t *testing.T) {
	type testCase struct {
		testDeposit  *Deposit
		testAgent    core.Agent
		expectedCost float64
	}

	tests := map[string]testCase{
		"ActionCost works": {
			testDeposit:  testDeposit,
			testAgent:    testAgent,
			expectedCost: 1.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCost, testDeposit.Cost(tc.testAgent))
		})
	}
}

func TestDeposit_String(t *testing.T) {
	type testCase struct {
		testDeposit    *Deposit
		expectedString string
	}

	tests := map[string]testCase{
		"basic deposit message": {
			testDeposit:    testDeposit,
			expectedString: "deposit 10 testResource at Location: testLocation\nCoordinates: (0, 0)\nInventory: {}\nAttributes: {}",
		},
		"deposit message with different Amount": {
			testDeposit: &Deposit{
				DepResource:    testResource,
				Amount:         100,
				ActionLocation: &core.Location{Name: "testLocation", Inventory: core.NewInventory()},
			},
			expectedString: "deposit 100 testResource at Location: testLocation\nCoordinates: (0, 0)\nInventory: {}\nAttributes: {}",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.testDeposit.Description()
			assert.Equal(t, tc.expectedString, output)
		})
	}
}
