package agent

import (
	"Neolithic/internal/core"
	"testing"

	"Neolithic/internal/planner"

	"github.com/stretchr/testify/require"
)

const deltaTime = 1.0 / 60

func TestPerforming_Execute(t *testing.T) {
	type testCase struct {
		timeLeft                    float64
		action                      planner.Action
		plan                        Plan
		startWorldState             *core.WorldState
		expectedAmountInEndLocation int
		expectedAgent               *Agent
		expectedAction              planner.Action
		expectedTimeLeft            float64
		expectedError               error
		nilEndState                 bool
	}

	tests := map[string]testCase{
		"can perform instant action": {
			plan: &mockPlan{
				nextAction: &mockAction{},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			expectedAmountInEndLocation: 1,
			expectedAgent: &Agent{
				Behavior: &Behavior{
					CurPlan: &mockPlan{
						nextAction: &mockAction{},
					},
					curState: &Moving{},
				},
			},
			expectedAction:   &mockAction{},
			expectedTimeLeft: 0,
		},
		"will tick down time": {
			plan: &mockPlan{
				nextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": core.Location{
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			expectedAgent: &Agent{
				Behavior: &Behavior{
					CurPlan: &mockPlan{
						nextAction: &mockActionWithTime{timeNeeded: 1.0},
					},
				},
			},
			expectedAction:   &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft: 1.0 - deltaTime,
			nilEndState:      true,
		},
		"action fails, reset to idle": {
			plan: &mockPlan{
				nextAction: &mockNullAction{},
			},
			startWorldState: &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": core.Location{
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			expectedAgent: &Agent{
				Behavior: &Behavior{
					CurPlan: &mockPlan{
						nextAction: &mockNullAction{},
					},
					curState: &Idle{},
				},
			},
			expectedAction:   &mockNullAction{},
			expectedTimeLeft: 0,
			nilEndState:      true,
		},
		"action succeeds, time left is zero": {
			plan: &mockPlan{
				nextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			action:   &mockActionWithTime{timeNeeded: 1.0},
			timeLeft: 0,
			startWorldState: &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			expectedAmountInEndLocation: 1,
			expectedAgent: &Agent{
				Behavior: &Behavior{
					CurPlan: &mockPlan{
						nextAction: &mockActionWithTime{timeNeeded: 1.0},
					},
					curState: &Moving{},
				},
			},
			expectedAction:   &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft: 0,
		},
		"action succeeds, plan is complete": {
			plan: &mockPlan{
				nextAction: &mockAction{},
				isComplete: true,
			},
			startWorldState: &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": {
						Name:      "testLocation",
						Inventory: core.NewInventory(),
					},
				},
				Agents: map[string]core.Agent{},
			},
			expectedAmountInEndLocation: 1,
			expectedAgent: &Agent{
				Behavior: &Behavior{
					CurPlan: &mockPlan{
						nextAction: &mockAction{},
						isComplete: true,
					},
					curState: &Idle{},
				},
			},
			expectedAction:   &mockAction{},
			expectedTimeLeft: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testPerforming := &Performing{
				timeLeft: tc.timeLeft,
				action:   tc.action,
				agent: &Agent{
					Behavior: &Behavior{CurPlan: tc.plan},
				},
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
			require.Equal(t, tc.expectedAmountInEndLocation, output.Locations["testLocation"].Inventory.GetAmount(testResource))
		})
	}
}
