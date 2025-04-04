package agent

import (
	"Neolithic/internal/core"
	"Neolithic/internal/logging"
	"testing"

	"Neolithic/internal/astar"
	"Neolithic/internal/planner"

	"github.com/stretchr/testify/require"
)

func TestIdle_Execute(t *testing.T) {
	type testCase struct {
		iterationsPerCall  int
		planner            *astar.SearchState
		worldState         *core.WorldState
		expectedIterations int
		expectedPlan       Plan
	}

	// Create separate locations for start and end states
	startLocation := core.Location{
		Name:      "testLocation",
		Inventory: core.NewInventory(),
	}

	endLocation := core.Location{
		Name:      "testLocation",
		Inventory: core.NewInventory(),
	}
	endLocation.Inventory.AdjustAmount(testResource, 3)

	testStart := &core.WorldState{
		Locations: map[string]core.Location{
			"testLocation": startLocation,
		},
		Agents: map[string]core.Agent{},
	}

	testEnd := &core.WorldState{
		Locations: map[string]core.Location{
			"testLocation": endLocation,
		},
		Agents: map[string]core.Agent{},
	}

	agentBehavior := &Behavior{
		PossibleActions: []planner.Action{
			&mockAction{},
		},
		Goal: testEnd,
	}

	testAgent := &Agent{
		Behavior: agentBehavior,
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
				IterationsPerCall: test.iterationsPerCall,
				planner:           test.planner,
				agent:             testAgent,
				logger:            logging.NewLogger("info"),
			}
			_, err := testIdle.Execute(test.worldState, 0)
			require.NoError(t, err)

			require.Equal(t, test.expectedIterations, testIdle.planner.Iterations)
			require.Equal(t, expectedStart, testIdle.planner.Start)
			require.Equal(t, expectedGoal, testIdle.planner.Goal)
			require.Equal(t, testAgent.Behavior.CurPlan, testIdle.agent.Behavior.CurPlan)

		})
	}
}
