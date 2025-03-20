package planner

import (
	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDeposit = &Deposit{
	resource: testResource,
	amount:   10,
	locName:  "testLocation",
	cost:     1.0,
}

func TestDeposit_Perform(t *testing.T) {
	type testCase struct {
		testDeposit              *Deposit
		startLocation            core.Location
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
			startLocation:            testLocation,
			startAmountInLocation:    0,
			agent:                    testAgent,
			startAmountInAgent:       10,
			expectedAmountInLocation: 10,
			expectedAmountInAgent:    0,
		},
		"deposit fails, nothing in agent inventory": {
			testDeposit:              testDeposit,
			startLocation:            testLocation,
			startAmountInLocation:    0,
			agent:                    testAgent,
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    0,
			expectNil:                true,
		},
		"partial deposit success": {
			testDeposit:              testDeposit,
			startLocation:            testLocation,
			startAmountInLocation:    0,
			agent:                    testAgent,
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
				Locations: map[string]core.Location{"testLocation": tc.startLocation},
				Agents:    map[string]core.Agent{"testAgent": tc.agent},
			}

			endState := tc.testDeposit.Perform(startState, tc.agent)
			if tc.expectNil {
				assert.Nil(t, endState)
				return
			}
			endAgent := endState.Agents["testAgent"]
			endLocation := endState.Locations["testLocation"]
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
		"cost works": {
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
			expectedString: "deposit 10 testResource at testLocation",
		},
		"deposit message with different amount": {
			testDeposit: &Deposit{
				resource: testResource,
				amount:   100,
				locName:  "testLocation",
			},
			expectedString: "deposit 100 testResource at testLocation",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.testDeposit.Description()
			assert.Equal(t, tc.expectedString, output)
		})
	}
}
