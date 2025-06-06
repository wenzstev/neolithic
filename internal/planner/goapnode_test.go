package planner

import (
	"Neolithic/internal/logging"
	"math"
	"testing"

	"Neolithic/internal/astar"
	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestActions_AStar(t *testing.T) {

	actionList := []core.Action{
		gatherTest,
		gatherTest2,
		depositTest,
		depositTest2,
	}

	type testCase struct {
		startLocationAmount int
		goalLocationAmount  int
		expectedError       error
		expectedActionList  []core.Action
		expectedIterations  int
		expectedCost        float64
	}

	tests := map[string]testCase{
		"can find gather path": {
			startLocationAmount: 100,
			goalLocationAmount:  20,
			expectedActionList: []core.Action{
				nil,
				gatherTest,
				gatherTest,
				depositTest2,
			},
			expectedIterations: 5,
			expectedCost:       21,
		},

		"can move all resource to new location": {
			startLocationAmount: 50,
			goalLocationAmount:  50,
			expectedActionList: []core.Action{
				nil,
				gatherTest,
				gatherTest,
				gatherTest,
				depositTest2,
				gatherTest,
				gatherTest,
				depositTest2,
				depositTest2,
			},
			expectedIterations: 17,
			expectedCost:       53,
		},
		"will return error if no path": {
			startLocationAmount: 20,
			goalLocationAmount:  100,
			expectedError:       astar.ErrNoPath,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			startState := &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name:  testLocation.DeepCopy(),
					testLocation2.Name: testLocation2.DeepCopy(),
				},
				Agents: map[string]core.Agent{
					testAgent.Name(): testAgent.DeepCopy(),
				},
			}

			startLoc, exists := startState.GetLocation("testLocation")
			require.True(t, exists)
			startLoc.Inventory.AdjustAmount(testResource, tc.startLocationAmount)

			goalState := &core.WorldState{
				Locations: map[string]*core.Location{testLocation2.Name: testLocation2.DeepCopy()},
				Agents:    map[string]core.Agent{},
			}
			goalLoc, exists := goalState.GetLocation("testLocation2")
			require.True(t, exists)
			goalLoc.Inventory.AdjustAmount(testResource, tc.goalLocationAmount)

			runInfo := &GoapRunInfo{
				Agent:               testAgent,
				PossibleNextActions: actionList,
			}

			startNode := &GoapNode{
				State:       startState,
				GoapRunInfo: runInfo,
			}

			endNode := &GoapNode{
				State:       goalState,
				GoapRunInfo: runInfo,
			}

			search, err := astar.NewSearch(startNode, endNode, astar.WithLogger(logging.NewLogger("debug")))
			assert.NoError(t, err)

			err = search.RunIterations(1000)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCost, search.BestCost)
			assert.Equal(t, tc.expectedIterations, search.Iterations)

			solutionActions := make([]core.Action, 0)
			solution := search.CurrentBestPath()
			for _, node := range solution {
				goapNode, ok := node.(*GoapNode)
				assert.True(t, ok)
				solutionActions = append(solutionActions, goapNode.Action)
			}
			assert.Equal(t, tc.expectedActionList, solutionActions)
			assert.Equal(t, tc.expectedIterations, search.Iterations)
			assert.Equal(t, tc.expectedCost, search.BestCost)
		})
	}
}

func TestActions_Heuristic(t *testing.T) {
	type testCase struct {
		amountInCurState   int
		amountInGoalState  int
		amountInStartAgent int
		expectedDistance   float64
	}

	tests := map[string]testCase{
		"goal is reached": {
			amountInCurState:  10,
			amountInGoalState: 10,
			expectedDistance:  0.0,
		},
		"impossible to reach goal": {
			amountInCurState:  0,
			amountInGoalState: 100,
			expectedDistance:  math.Inf(1),
		},
		"State amount is less than goal": {
			amountInCurState:   50,
			amountInGoalState:  70,
			amountInStartAgent: 20,
			expectedDistance:   1.0,
		},
		"State amount is greater than goal": {
			amountInCurState:  70,
			amountInGoalState: 50,
			expectedDistance:  20.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			curState := &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name:  testLocation.DeepCopy(),
					testLocation2.Name: testLocation2.DeepCopy(), // extra location to make sure we ignore it
				},
				Agents: map[string]core.Agent{
					testAgent.Name(): testAgent.DeepCopy(),
				},
			}
			curLoc, exists := curState.GetLocation("testLocation")
			require.True(t, exists)
			curLoc.Inventory.AdjustAmount(testResource, tc.amountInCurState)
			curAg, exists := curState.GetAgent("testAgent")
			require.True(t, exists)
			curAg.Inventory().AdjustAmount(testResource, tc.amountInStartAgent)

			goalState := &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name: testLocation.DeepCopy(),
				},
				Agents: map[string]core.Agent{},
			}
			goalLoc, exists := goalState.GetLocation("testLocation")
			require.True(t, exists)
			goalLoc.Inventory.AdjustAmount(testResource, tc.amountInGoalState)

			testStats := &GoapRunInfo{
				Agent: testAgent,
				PossibleNextActions: []core.Action{
					gatherTest,
					gatherTest2,
					depositTest,
					depositTest2,
				},
			}
			testNode := &GoapNode{
				State:       curState,
				GoapRunInfo: testStats,
			}

			goalNode := &GoapNode{
				State:       goalState,
				GoapRunInfo: testStats,
			}

			val, err := testNode.Heuristic(goalNode)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedDistance, val)
		})
	}
}

func TestActions_GetSuccessors(t *testing.T) {
	mockAction1 := &mockAction{}
	mockAction2 := &mockAction{}
	mockAction3 := &mockAction{}
	mockNullAction1 := &mockNullAction{}

	expectedEndState := &core.WorldState{
		Locations: map[string]*core.Location{
			testLocation.Name: testLocation.DeepCopy(),
		},
		Agents: map[string]core.Agent{
			testAgent.Name(): testAgent.DeepCopy(),
		},
	}
	expectedLoc, exists := expectedEndState.GetLocation("testLocation")
	require.True(t, exists)
	expectedLoc.Inventory.AdjustAmount(testResource, 1)

	type testCase struct {
		actions            []core.Action
		startState         *core.WorldState
		agent              core.Agent
		expectedSuccessors []*GoapNode
	}

	tests := map[string]testCase{
		"single action should generate one successor": {
			actions: []core.Action{mockAction1},
			startState: &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name: testLocation.DeepCopy(),
				},
				Agents: map[string]core.Agent{
					testAgent.Name(): testAgent.DeepCopy(),
				},
			},
			agent: testAgent,
			expectedSuccessors: []*GoapNode{
				{
					Action: mockAction1,
					State:  expectedEndState,
				},
			},
		},
		"multiple actions should generate multiple successors": {
			actions: []core.Action{
				mockAction1,
				mockAction2,
				mockAction3,
			},
			startState: &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name: testLocation.DeepCopy(),
				},
				Agents: map[string]core.Agent{
					testAgent.Name(): testAgent.DeepCopy(),
				},
			},
			agent: testAgent,
			expectedSuccessors: []*GoapNode{
				{
					Action: mockAction1,
					State:  expectedEndState,
				},
				{
					Action: mockAction2,
					State:  expectedEndState,
				},
				{
					Action: mockAction3,
					State:  expectedEndState,
				},
			},
		},
		"null actions should be filtered out": {
			actions: []core.Action{
				mockNullAction1,
				mockAction1,
				mockAction2,
			},
			startState: &core.WorldState{
				Locations: map[string]*core.Location{
					testLocation.Name: testLocation.DeepCopy(),
				},
				Agents: map[string]core.Agent{
					testAgent.Name(): testAgent.DeepCopy(),
				},
			},
			agent: testAgent,
			expectedSuccessors: []*GoapNode{
				{
					Action: mockAction1,
					State:  expectedEndState,
				},
				{
					Action: mockAction2,
					State:  expectedEndState,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testNode := &GoapNode{
				Action: nil,
				State:  tc.startState,
				GoapRunInfo: &GoapRunInfo{
					Agent:               tc.agent,
					PossibleNextActions: tc.actions,
				},
			}
			output, err := testNode.GetSuccessors()

			successorList := make([]*GoapNode, 0)
			for _, successor := range output {
				goapSuccessor, ok := successor.(*GoapNode)
				assert.True(t, ok)
				goapSuccessor.GoapRunInfo = nil
				successorList = append(successorList, goapSuccessor)
			}

			assert.NoError(t, err)
			assert.Equal(t, len(tc.expectedSuccessors), len(successorList))
			for i, expected := range tc.expectedSuccessors {
				assert.Equal(t, expected.Action, successorList[i].Action)
				expectedStateLoc, ok := expected.State.GetLocation("testLocation")
				assert.True(t, ok)
				successorLoc, ok := successorList[i].State.GetLocation("testLocation")
				assert.True(t, ok)
				assert.Equal(t, expectedStateLoc.Inventory.GetAmount(testResource),
					successorLoc.Inventory.GetAmount(testResource))
			}
		})
	}
}
