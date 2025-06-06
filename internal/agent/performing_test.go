package agent

import (
	"testing"

	"Neolithic/internal/core"
	"Neolithic/internal/logging"
	"github.com/stretchr/testify/require"
)

const deltaTime = 1.0 / 60

func TestPerforming_Execute(t *testing.T) {
	type testCase struct {
		timeLeft                    float64
		action                      core.Action
		plan                        Plan
		startWorldState             *core.WorldState
		expectedAmountInEndLocation int
		startAgent                  *Agent
		expectedAgent               *Agent
		expectedAction              core.Action
		expectedTimeLeft            float64
		expectedError               error
		nilEndState                 bool
	}

	basicAgent := &Agent{
		name: "basicAgent",
		Behavior: &Behavior{
			CurPlan: &MockPlan{
				NextAction: &mockAction{},
			},
			CurState: &Moving{},
		},
	}

	basicTimeAgent := &Agent{
		name: "basicTimeAgent",
		Behavior: &Behavior{
			CurPlan: &MockPlan{
				NextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
		},
	}

	basicNullAgent := &Agent{
		name: "basicNullAgent",
		Behavior: &Behavior{
			CurPlan: &MockPlan{
				NextAction: &mockNullAction{},
				Complete:   true,
			},
			CurState: &Idle{},
		},
	}

	tests := map[string]testCase{
		"can perform instant action": {
			plan: &MockPlan{
				NextAction: &mockAction{},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]*core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{
					basicAgent.Name(): basicAgent,
				},
			},
			expectedAmountInEndLocation: 1,
			startAgent:                  basicAgent,
			expectedAgent:               basicAgent,
			expectedAction:              &mockAction{},
			expectedTimeLeft:            0,
		},
		"will tick down time": {
			plan: &MockPlan{
				NextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]*core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			startAgent:       basicTimeAgent,
			expectedAgent:    basicTimeAgent,
			expectedAction:   &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft: 1.0 - deltaTime,
			nilEndState:      true,
		},
		"action fails, reset to idle": {
			plan: &MockPlan{
				NextAction: &mockNullAction{},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]*core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{
					basicNullAgent.Name(): basicNullAgent,
				},
			},
			startAgent:       basicNullAgent,
			expectedAgent:    basicNullAgent,
			expectedAction:   &mockNullAction{},
			expectedTimeLeft: 0,
			nilEndState:      true,
		},
		"action succeeds, time left is zero": {
			plan: &MockPlan{
				NextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			action:   &mockActionWithTime{timeNeeded: 1.0},
			timeLeft: 0,
			startWorldState: &core.WorldState{
				Locations: map[string]*core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{
					basicTimeAgent.Name(): basicTimeAgent,
				},
			},
			expectedAmountInEndLocation: 1,
			startAgent:                  basicTimeAgent,
			expectedAgent:               basicTimeAgent,
			expectedAction:              &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft:            0,
		},
		"action succeeds, plan is complete": {
			plan: &MockPlan{
				NextAction: &mockAction{},
				Complete:   true,
			},
			startWorldState: &core.WorldState{
				Locations: map[string]*core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{
					basicAgent.Name(): basicAgent,
				},
			},
			expectedAction:              &mockAction{},
			expectedAmountInEndLocation: 1,
			startAgent:                  basicAgent,
			expectedAgent:               basicAgent,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testPerforming := &Performing{
				timeLeft: tc.timeLeft,
				action:   tc.action,
				agent:    tc.startAgent,
				logger:   logging.NewLogger("info"),
			}
			output, err := testPerforming.Execute(tc.startWorldState, 1.0/60.0)
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
				return
			}
			require.Equal(t, tc.expectedAgent, testPerforming.agent)
			require.Equal(t, tc.expectedAction, testPerforming.action)
			require.Equal(t, tc.expectedTimeLeft, testPerforming.timeLeft)
			if tc.nilEndState {
				require.Nil(t, output)
				return
			}
			outputLocation, exists := output.GetLocation("testLocation")
			require.True(t, exists)
			require.Equal(t, tc.expectedAmountInEndLocation, outputLocation.Inventory.GetAmount(testResource))
		})
	}
}
