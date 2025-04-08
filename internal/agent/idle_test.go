package agent

import (
	"testing"

	"Neolithic/internal/astar"
	"Neolithic/internal/core"
	"Neolithic/internal/goalengine"
	"Neolithic/internal/logging"
	"Neolithic/internal/planner"

	"github.com/stretchr/testify/require"
)

var testChunkerFunc goalengine.ChunkerFunc = func(location *core.Location, resource *core.Resource) *core.WorldState {
	goalLocation := core.Location{
		Name:      location.Name,
		Inventory: core.NewInventory(),
	}
	goalLocation.Inventory.AdjustAmount(resource, 3)

	return &core.WorldState{
		Locations: map[string]core.Location{
			location.Name: goalLocation,
		},
	}
}

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
		GoalEngine: &goalengine.GoalEngine{
			Goal: goalengine.Goal{
				Name: "testGoal",
				Logic: goalengine.GoalLogic{
					Chunker:      testChunkerFunc,
					Fallback:     goalengine.FallbackChunkFunc,
					ShouldGiveUp: goalengine.GiveUpIfLessThanFive,
				},
				Location: &core.Location{
					Name: "testLocation",
				},
				Resource: testResource,
			},
		},
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
		"can create planner and run 2 iterations": {
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

func TestIdle_Execute2(t *testing.T) {
	type testCase struct {
		iterationsPerCall  int
		planner            *astar.SearchState
		startLocation      core.Location
		goalEngine         *goalengine.GoalEngine
		possibleActions    []planner.Action
		expectedIterations int
		expectedPlan       Plan
		expectedRetries    int
		expectedError      error
	}

	defaultGoalEngine := &goalengine.GoalEngine{
		Goal: goalengine.Goal{
			Name: "testGoal",
			Logic: goalengine.GoalLogic{
				Chunker:      testChunkerFunc,
				Fallback:     goalengine.FallbackChunkFunc,
				ShouldGiveUp: goalengine.GiveUpIfLessThanFive,
			},
			Location: &core.Location{
				Name: "testLocation",
			},
			Resource: testResource,
		},
	}

	// Create common inventories
	emptyInventory := core.NewInventory()
	endInventory := core.NewInventory()
	endInventory.AdjustAmount(testResource, 3)

	tests := map[string]testCase{
		"can create planner and run 2 iterations": {
			iterationsPerCall: 2,
			startLocation: core.Location{
				Name:      "testLocation",
				Inventory: emptyInventory,
			},
			goalEngine:         defaultGoalEngine,
			possibleActions:    []planner.Action{&mockAction{}},
			expectedIterations: 2,
			expectedRetries:    0,
		},
		"can create planner and complete plan": {
			iterationsPerCall: 15,
			startLocation: core.Location{
				Name:      "testLocation",
				Inventory: emptyInventory,
			},
			goalEngine:         defaultGoalEngine,
			possibleActions:    []planner.Action{&mockAction{}},
			expectedIterations: 4,
			expectedRetries:    0,
			expectedPlan: &plan{
				Actions: []planner.Action{
					&mockAction{},
					&mockAction{},
					&mockAction{},
				},
			},
		},
		"increments retries when no path found": {
			iterationsPerCall: 100,
			startLocation: core.Location{
				Name:      "testLocation",
				Inventory: emptyInventory,
			},
			goalEngine: &goalengine.GoalEngine{
				Goal: goalengine.Goal{
					Name: "testGoal",
					Logic: goalengine.GoalLogic{
						Chunker:      testChunkerFunc,
						Fallback:     goalengine.FallbackChunkFunc,
						ShouldGiveUp: goalengine.GiveUpIfLessThanFive,
					},
					Location: &core.Location{
						Name: "testLocation",
					},
					Resource: testResource,
				},
			},
			possibleActions:    []planner.Action{&mockNilAction{}},
			expectedIterations: 1,
			expectedRetries:    1,
			expectedError:      astar.ErrNoPath,
		},
		"increments retries when iterations run out": {
			iterationsPerCall: 10,
			startLocation: core.Location{
				Name:      "testLocation",
				Inventory: emptyInventory,
			},
			goalEngine: &goalengine.GoalEngine{
				Goal: goalengine.Goal{
					Name: "testGoal",
					Logic: goalengine.GoalLogic{
						Chunker:      goalengine.AddToLocation,
						Fallback:     goalengine.FallbackChunkFunc,
						ShouldGiveUp: goalengine.GiveUpIfLessThanFive,
					},
					Location: &core.Location{
						Name: "testLocation",
					},
					Resource: testResource,
				},
			},
			possibleActions:    []planner.Action{&mockAction{}},
			expectedIterations: 10,
			expectedRetries:    1,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Create world states for this specific test
			testStart := &core.WorldState{
				Locations: map[string]core.Location{
					"testLocation": test.startLocation,
				},
				Agents: map[string]core.Agent{},
			}

			// Create agent behavior for this specific test
			agentBehavior := &Behavior{
				PossibleActions: test.possibleActions,
				GoalEngine:      test.goalEngine,
			}

			testAgent := &Agent{
				Behavior: agentBehavior,
			}

			expectedStart := &planner.GoapNode{
				State: testStart,
				GoapRunInfo: &planner.GoapRunInfo{
					Agent:               testAgent,
					PossibleNextActions: test.possibleActions,
				},
			}

			testIdle := &Idle{
				IterationsPerCall: test.iterationsPerCall,
				planner:           test.planner,
				agent:             testAgent,
				logger:            logging.NewLogger("info"),
			}

			_, err := testIdle.Execute(testStart, 0)

			if test.expectedError != nil {
				require.ErrorIs(t, err, astar.ErrNoPath)
				require.Equal(t, test.expectedRetries, testIdle.numRetries)
			} else {
				require.NoError(t, err)
				require.Equal(t, test.expectedIterations, testIdle.planner.Iterations)
				require.Equal(t, expectedStart, testIdle.planner.Start)
				require.Equal(t, testAgent.Behavior.CurPlan, testIdle.agent.Behavior.CurPlan)
			}
		})
	}
}
