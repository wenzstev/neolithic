package attributes

import (
	"testing"

	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
)

func TestGather_Perform(t *testing.T) {
	testGather := &Gather{
		Res:            testResource,
		Amount:         5,
		ActionLocation: &core.Location{Name: "testLocation"},
		ActionCost:     1,
	}

	testTool := &core.Resource{Name: "testTool"}

	testGatherRequires := &Gather{
		Requires:       testTool,
		Res:            testResource,
		Amount:         5,
		ActionLocation: &core.Location{Name: "testLocation"},
		ActionCost:     1,
	}

	type testCase struct {
		testGather               *Gather
		startLocation            *core.Location
		startAmountInLocation    int
		agent                    core.Agent
		startAmountInAgent       int
		toolInAgent              *core.Resource
		expectedAmountInLocation int
		expectedAmountInAgent    int
		expectNil                bool
	}

	testCases := map[string]testCase{
		"can do basic gather": {
			testGather:               testGather,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    5,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    5,
		},
		"gather partially succeeds": {
			testGather:               testGather,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    2,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    2,
		},
		"gather succeeds with tool": {
			testGather:               testGatherRequires,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    5,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			toolInAgent:              testTool,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    5,
		},
		"gather fails, no resource in location": {
			testGather:               testGather,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    0,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    0,
			expectNil:                true,
		},
		"gather fails, required tool not present": {
			testGather:               testGatherRequires,
			startLocation:            testLocation.DeepCopy(),
			startAmountInLocation:    5,
			agent:                    testAgent.DeepCopy(),
			startAmountInAgent:       0,
			expectedAmountInLocation: 0,
			expectedAmountInAgent:    0,
			expectNil:                true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			tc.startLocation.Inventory.AdjustAmount(testResource, tc.startAmountInLocation)
			tc.agent.Inventory().AdjustAmount(testResource, tc.startAmountInAgent)
			if tc.toolInAgent != nil {
				tc.agent.Inventory().AdjustAmount(tc.toolInAgent, 1)
			}
			startState := &core.WorldState{
				Locations: []core.Location{*tc.startLocation},
				Agents:    []core.Agent{tc.agent},
			}

			endState := tc.testGather.Perform(startState, tc.agent)
			if tc.expectNil {
				assert.Nil(t, endState)
				return
			}
			endLoc, exists := endState.GetLocation(tc.startLocation.Name)
			assert.True(t, exists)
			assert.Equal(t, tc.expectedAmountInLocation, endLoc.Inventory.GetAmount(testResource))
			endAgent, exists := endState.GetAgent(tc.agent.Name())
			assert.True(t, exists)
			assert.Equal(t, tc.expectedAmountInAgent, endAgent.Inventory().GetAmount(testResource))
		})
	}
}
