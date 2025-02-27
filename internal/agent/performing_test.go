package agent

import (
	"Neolithic/internal/planner"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestPerforming_Execute(t *testing.T) {
	type testCase struct {
		testPerforming     *Performing
		testWorldState     *planner.State
		expectedWorldState *planner.State
		expectedPerforming *Performing
		expectedError      error
	}

	tests := map[string]testCase{
		"can perform instant action": {
			testPerforming: &Performing{
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockAction{},
						},
					},
				},
			},
			testWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedPerforming: &Performing{
				action: &mockAction{},
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockAction{},
						},
						curState: &Moving{},
					},
				},
			},
		},
		"can perform action after time": {
			testPerforming: &Performing{
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockActionWithTime{timeNeeded: 1.0},
						},
					},
				},
			},
			testWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedWorldState: nil,
			expectedPerforming: &Performing{
				action:   &mockActionWithTime{timeNeeded: 1.0},
				timeLeft: 1.0 - (1.0 / 60.0),
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockActionWithTime{timeNeeded: 1.0},
						},
					},
				},
			},
		},
		"action fails, reset to idle": {
			testPerforming: &Performing{
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockNullAction{},
						},
					},
				},
			},
			testWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedWorldState: nil,
			expectedPerforming: &Performing{
				action: &mockNullAction{},
				agent: &agent{
					name: "testAgent",
					behavior: &Behavior{
						curPlan: &mockPlan{
							isComplete: false,
							nextAction: &mockNullAction{},
						},
						curState: &Idle{},
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := tc.testPerforming.Execute(tc.testWorldState, 1.0/60.0)
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
			}
			require.Equal(t, tc.expectedWorldState, output)
			require.Equal(t, tc.expectedPerforming, tc.testPerforming)
		})
	}
}
