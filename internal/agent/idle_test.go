package agent

import (
	"Neolithic/internal/core"
	"testing"

	"Neolithic/internal/astar"
	"Neolithic/internal/planner"
	"github.com/stretchr/testify/require"
)

func TestIdle_Execute(t *testing.T) {
	type testCase struct {
		iterationsPerCall  int
		planner            *astar.SearchState
		worldState         WorldState
		expectedIterations int
		expectedPlan       Plan
	}

	testStart := &core.State{
		Locations: map[*planner.Location]planner.Inventory{
			testLocation: {},
		},
		Agents: map[planner.Agent]planner.Inventory{},
	}

	testEnd := &core.State{
		Locations: map[*planner.Location]planner.Inventory{
			testLocation: {
				testResource: 3,
			},
		},
		Agents: map[planner.Agent]planner.Inventory{},
	}

	agentBehavior := &Behavior{
		PossibleActions: &[]planner.Action{
			&mockAction{},
		},
		Goal: testEnd,
	}

	testAgent := &mockAgent{
		behavior: agentBehavior,
	}

	expectedStart := &planner.GoapNode{
		State: testStart,
		GoapRunInfo: &planner.GoapRunInfo{
			Agent:               testAgent,
			PossibleNextActions: agentBehavior.PossibleActions,
		},
	}

	expectedGoal := &planner.GoapNode{
		State: testEnd,
		GoapRunInfo: &planner.GoapRunInfo{
			Agent:               testAgent,
			PossibleNextActions: agentBehavior.PossibleActions,
		},
	}

	tests := map[string]testCase{
		"can create planner and run 5 iterations": {
			iterationsPerCall:  2,
			worldState:         testStart,
			expectedIterations: 2,
		},
		"can create planner and complete plan": {
			iterationsPerCall:  15,
			worldState:         testStart,
			expectedIterations: 4,
			expectedPlan: &plan{
				Actions: []planner.Action{
					&mockAction{},
					&mockAction{},
					&mockAction{},
				},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			testIdle := &Idle{
				iterationsPerCall: test.iterationsPerCall,
				planner:           test.planner,
				agent:             testAgent,
			}
			_, err := testIdle.Execute(test.worldState, 0)
			require.NoError(t, err)

			require.Equal(t, test.expectedIterations, testIdle.planner.Iterations)
			require.Equal(t, expectedStart, testIdle.planner.Start)
			require.Equal(t, expectedGoal, testIdle.planner.Goal)
			require.Equal(t, testAgent.Behavior().CurPlan, testIdle.agent.Behavior().CurPlan)

		})
	}
}
