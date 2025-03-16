package planner

import (
	"Neolithic/internal/astar"
	"Neolithic/internal/core"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActions_AStar(t *testing.T) {

	startState := &core.State{
		Locations: map[*Location]Inventory{
			testLocation: {
				testResource: 50,
			},
			testLocation2: {},
		},
		Agents: map[Agent]Inventory{
			testAgent: {},
		},
	}

	actionList := &[]Action{
		gatherTest,
		gatherTest2,
		depositTest,
		depositTest2,
	}

	type testCase struct {
		actions            *[]Action
		startState         *core.State
		goalState          *core.State
		agent              Agent
		maxDistance        int
		expectedOutput     *astar.SearchState
		expectedActionList []Action
		expectedError      error
	}

	tests := map[string]testCase{
		"can find gather path": {
			actions:    actionList,
			startState: startState,
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 20,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 10000,
			expectedOutput: &astar.SearchState{
				Start: &GoapNode{
					State: startState,
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				Goal: &GoapNode{
					State: &core.State{
						Locations: map[*Location]Inventory{
							testLocation2: {
								testResource: 20,
							},
						},
					},
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				BestCost:   21,
				Iterations: 5,
				FoundBest:  true,
			},
			expectedActionList: []Action{
				nil,
				gatherTest,
				gatherTest,
				depositTest2,
			},
		},

		"can move all resource to new location": {
			actions:    actionList,
			startState: startState,
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 50,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 10000,
			expectedOutput: &astar.SearchState{
				Start: &GoapNode{
					State: startState,
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				Goal: &GoapNode{
					State: &core.State{
						Locations: map[*Location]Inventory{
							testLocation2: {
								testResource: 50,
							},
						},
					},
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				Iterations: 17,
				FoundBest:  true,
				BestCost:   53,
			},
			expectedActionList: []Action{
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
		},
		"will return error if no path": {
			actions:    actionList,
			startState: startState,
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 100,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 10000,
			expectedOutput: &astar.SearchState{
				Start: &GoapNode{
					State: startState,
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				Goal: &GoapNode{
					State: &core.State{
						Locations: map[*Location]Inventory{
							testLocation: {
								testResource: 100,
							},
						},
					},
					GoapRunInfo: &GoapRunInfo{
						Agent:               testAgent,
						PossibleNextActions: actionList,
					},
				},
				Iterations: 21,
				FoundBest:  false,
				BestCost:   math.Inf(1),
			},
			expectedError: astar.ErrNoPath,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			runInfo := &GoapRunInfo{
				Agent:               tc.agent,
				PossibleNextActions: tc.actions,
			}

			startNode := &GoapNode{
				State:       tc.startState,
				GoapRunInfo: runInfo,
			}

			endNode := &GoapNode{
				State:       tc.goalState,
				GoapRunInfo: runInfo,
			}

			search, err := astar.NewSearch(startNode, endNode)
			assert.NoError(t, err)

			err = search.RunIterations(tc.maxDistance)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				return
			}
			assert.NoError(t, err)

			areNodesEqual(t, tc.expectedOutput.Start, search.Start)
			areNodesEqual(t, tc.expectedOutput.Goal, search.Goal)
			assert.Equal(t, tc.expectedOutput.BestCost, search.BestCost)
			assert.Equal(t, tc.expectedOutput.Iterations, search.Iterations)
			assert.Equal(t, tc.expectedOutput.FoundBest, search.FoundBest)

			solutionActions := make([]Action, 0)
			solution := search.CurrentBestPath()
			for _, node := range solution {
				goapNode, ok := node.(*GoapNode)
				assert.True(t, ok)
				solutionActions = append(solutionActions, goapNode.Action)
			}
			assert.Equal(t, tc.expectedActionList, solutionActions)
		})
	}
}

func areNodesEqual(t *testing.T, nodeA, nodeB astar.Node) {
	goapNodeA, ok := nodeA.(*GoapNode)
	assert.True(t, ok)

	goapNodeB, ok := nodeB.(*GoapNode)
	assert.True(t, ok)

	aId, err := goapNodeA.State.ID()
	assert.NoError(t, err)

	bId, err := goapNodeB.State.ID()
	assert.NoError(t, err)

	assert.Equal(t, aId, bId)
	assert.Equal(t, goapNodeA.Action, goapNodeB.Action)
	assert.Equal(t, goapNodeA.GoapRunInfo.Agent, goapNodeB.GoapRunInfo.Agent)
	assert.Equal(t, goapNodeA.GoapRunInfo.PossibleNextActions, goapNodeB.GoapRunInfo.PossibleNextActions)

}

func TestActions_Heuristic(t *testing.T) {
	type testCase struct {
		curState         *core.State
		goalState        *core.State
		agent            Agent
		expectedDistance float64
	}

	testResource2 := &Resource{Name: "testResource2"}

	tests := map[string]testCase{
		"goal is reached": {
			curState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 0,
		},
		"goal is reached, other nonimportant states present": {
			curState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource:  50,
						testResource2: 20,
					},
				},
			},
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 0,
		},
		"impossible to reach goal": {
			curState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: math.Inf(1),
		},
		"State amount is less than goal": {
			curState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testResource: 50,
					},
				},
			},
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 1,
		},
		"State amount is greater than goal": {
			curState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &core.State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 20.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testStats := &GoapRunInfo{
				Agent: tc.agent,
				PossibleNextActions: &[]Action{
					gatherTest,
					gatherTest2,
					depositTest,
					depositTest2,
				},
			}
			testNode := &GoapNode{
				State:       tc.curState,
				GoapRunInfo: testStats,
			}

			goalNode := &GoapNode{
				State:       tc.goalState,
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

	expectedEndState := &core.State{
		Locations: map[*Location]Inventory{
			testLocation: {
				testResource: 1,
			},
		},
		Agents: map[Agent]Inventory{},
	}

	type testCase struct {
		actions            []Action
		startState         *core.State
		agent              Agent
		expectedSuccessors []*GoapNode
	}

	tests := map[string]testCase{
		"get successors with one Action": {
			actions: []Action{
				mockAction1,
			},
			startState: &core.State{},
			agent:      testAgent,
			expectedSuccessors: []*GoapNode{
				{
					Action: mockAction1,
					State:  expectedEndState,
				},
			},
		},
		"get successors with multiple actions": {
			actions: []Action{
				mockAction1,
				mockAction2,
				mockAction3,
			},
			startState: &core.State{},
			agent:      testAgent,
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
		"get successors with some null actions": {
			actions: []Action{
				mockNullAction1,
				mockAction1,
				mockAction2,
			},
			startState: &core.State{},
			agent:      testAgent,
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
					PossibleNextActions: &tc.actions,
				},
			}
			output, err := testNode.GetSuccessors()

			successorList := make([]*GoapNode, 0)
			for _, successor := range output {
				goapSuccessor, ok := successor.(*GoapNode)
				assert.True(t, ok)
				goapSuccessor.GoapRunInfo = nil // zero out the run info to make test easier to compare
				successorList = append(successorList, goapSuccessor)
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedSuccessors, successorList)
		})
	}
}
