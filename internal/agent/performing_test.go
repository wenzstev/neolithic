package agent

import (
	"testing"

	"Neolithic/internal/planner"
	"github.com/stretchr/testify/require"
)

const deltaTime = 1.0 / 60

func TestPerforming_Execute(t *testing.T) {
	type testCase struct {
		timeLeft         float64
		action           planner.Action
		plan             Plan
		startWorldState  *planner.State
		endWorldState    *planner.State
		expectedAgent    *mockAgent
		expectedAction   planner.Action
		expectedTimeLeft float64
		expectedError    error
	}

	testWorldState := &planner.State{
		Locations: map[*planner.Location]planner.Inventory{
			testLocation: {},
		},
		Agents: map[planner.Agent]planner.Inventory{},
	}

	tests := map[string]testCase{
		"can perform instant action": {
			plan: &mockPlan{
				nextAction: &mockAction{},
			},
			startWorldState: testWorldState,
			endWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedAgent: &mockAgent{
				plan: &mockPlan{
					nextAction: &mockAction{},
				},
				curState: &Moving{},
			},
			expectedAction:   &mockAction{},
			expectedTimeLeft: 0,
		},
		"will tick down time": {
			plan: &mockPlan{
				nextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			startWorldState: testWorldState,
			expectedAgent: &mockAgent{
				plan: &mockPlan{
					nextAction: &mockActionWithTime{timeNeeded: 1.0},
				},
			},
			expectedAction:   &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft: 1.0 - deltaTime,
		},
		"action fails, reset to idle": {
			plan: &mockPlan{
				nextAction: &mockNullAction{},
			},
			startWorldState: testWorldState,
			expectedAgent: &mockAgent{
				plan: &mockPlan{
					nextAction: &mockNullAction{},
				},
				curState: &Idle{},
			},
			expectedAction:   &mockNullAction{},
			expectedTimeLeft: 0,
		},
		"action succeeds, time left is zero": {
			plan: &mockPlan{
				nextAction: &mockActionWithTime{timeNeeded: 1.0},
			},
			action:          &mockActionWithTime{timeNeeded: 1.0},
			timeLeft:        0,
			startWorldState: testWorldState,
			endWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedAgent: &mockAgent{
				plan: &mockPlan{
					nextAction: &mockActionWithTime{timeNeeded: 1.0},
				},
				curState: &Moving{},
			},
			expectedAction:   &mockActionWithTime{timeNeeded: 1.0},
			expectedTimeLeft: 0,
		},
		"action succeeds, plan is complete": {
			plan: &mockPlan{
				nextAction: &mockAction{},
				isComplete: true,
			},
			startWorldState: testWorldState,
			endWorldState: &planner.State{
				Locations: map[*planner.Location]planner.Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[planner.Agent]planner.Inventory{},
			},
			expectedAgent: &mockAgent{
				plan: &mockPlan{
					nextAction: &mockAction{},
					isComplete: true,
				},
				curState: &Idle{},
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
				agent: &mockAgent{
					plan: tc.plan,
				},
			}
			output, err := testPerforming.Execute(tc.startWorldState, 1.0/60.0)
			if tc.expectedError != nil {
				require.Equal(t, tc.expectedError, err)
				return
			}
			require.Equal(t, tc.endWorldState, output)
			require.Equal(t, tc.expectedAgent, testPerforming.agent)
			require.Equal(t, tc.expectedAction, testPerforming.action)
			require.Equal(t, tc.expectedTimeLeft, testPerforming.timeLeft)
		})
	}
}
