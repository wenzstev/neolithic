package planner

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestActions_AStar(t *testing.T) {
	testLocation2 := &Location{Name: "testLocation2"}

	gatherTest := &Gather{
		resource: testResource,
		amount:   10,
		location: testLocation,
		cost:     10.0,
	}

	gatherTest2 := &Gather{
		resource: testResource,
		amount:   10,
		location: testLocation2,
		cost:     10.0,
	}

	depositTest := &Deposit{
		resource: testResource,
		amount:   10,
		location: testLocation,
		cost:     1.0,
	}

	depositTest2 := &Deposit{
		resource: testResource,
		amount:   10,
		location: testLocation2,
		cost:     1.0,
	}

	startState := &State{
		Locations: map[*Location]Inventory{
			testLocation: {
				testResource: 50,
			},
			testLocation2: {},
		},
		Agents: map[*Agent]Inventory{
			testAgent: {},
		},
	}

	actionList := Actions{
		gatherTest,
		gatherTest2,
		depositTest,
		depositTest2,
	}

	type testCase struct {
		actions        Actions
		startState     *State
		goalState      *State
		agent          *Agent
		maxDistance    int
		expectedOutput *AStarOutput
		expectedError  error
	}

	tests := map[string]testCase{
		"can find gather path": {
			actions:    actionList,
			startState: startState,
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 20,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 10000,
			expectedOutput: &AStarOutput{
				actions: []Action{
					depositTest2,
					gatherTest,
					depositTest2,
					gatherTest,
				},
				totalCost: 22,
				expectedState: &State{
					Locations: map[*Location]Inventory{
						testLocation: {
							testResource: 30,
						},
						testLocation2: {
							testResource: 20,
						},
					},
					Agents: map[*Agent]Inventory{
						testAgent: {},
					},
				},
			},
		},
		"can move all resource to new location": {
			actions:    actionList,
			startState: startState,
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 50,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 10000,
			expectedOutput: &AStarOutput{
				actions: []Action{
					depositTest2,
					gatherTest,
					depositTest2,
					depositTest2,
					depositTest2,
					depositTest2,
					gatherTest,
					gatherTest,
					gatherTest,
					gatherTest,
				},
				totalCost: 55,
				expectedState: &State{
					Locations: map[*Location]Inventory{
						testLocation: {},
						testLocation2: {
							testResource: 50,
						},
					},
					Agents: map[*Agent]Inventory{
						testAgent: {},
					},
				},
			},
		},
		"will return error if no path": {
			actions:    actionList,
			startState: startState,
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 100,
					},
				},
			},
			agent:          testAgent,
			maxDistance:    10000,
			expectedOutput: nil,
			expectedError:  ErrNoPath,
		},
		"will return partial path": {
			actions:    actionList,
			startState: startState,
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 50,
					},
				},
			},
			agent:       testAgent,
			maxDistance: 2,
			expectedOutput: &AStarOutput{
				actions: []Action{
					depositTest2,
					gatherTest,
				},
				totalCost: 11,
				expectedState: &State{
					Locations: map[*Location]Inventory{
						testLocation: {
							testResource: 40,
						},
						testLocation2: {
							testResource: 10,
						},
					},
					Agents: map[*Agent]Inventory{
						testAgent: {},
					},
				},
			},
		},
		"will err if the agent is not in the state": {
			actions: actionList,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation2: {
						testResource: 50,
					},
				},
			},
			agent:          testAgent,
			maxDistance:    10000,
			expectedOutput: nil,
			expectedError:  ErrNoAgent,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output, err := tc.actions.AStar(tc.startState, tc.goalState, tc.agent, tc.maxDistance)
			if tc.expectedError != nil {
				assert.ErrorIs(t, err, tc.expectedError)
				assert.Nil(t, output)
				return
			}
			assert.NoError(t, err)

			// sort both lists because sometimes there will be differences in order
			sort.Slice(tc.expectedOutput.actions, func(i, j int) bool {
				actionI := tc.expectedOutput.actions[i]
				actionJ := tc.expectedOutput.actions[j]
				return actionI.Description() < actionJ.Description()
			})
			sort.Slice(output.actions, func(i, j int) bool {
				actionI := output.actions[i]
				actionJ := output.actions[j]
				return actionI.Description() < actionJ.Description()
			})

			assert.Equal(t, tc.expectedOutput, output)
		})
	}
}

func TestActions_Heuristic(t *testing.T) {
	type testCase struct {
		curState         *State
		goalState        *State
		agent            *Agent
		expectedDistance float64
	}

	testResource2 := &Resource{Name: "testResource2"}

	tests := map[string]testCase{
		"goal is reached": {
			curState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &State{
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
			curState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource:  50,
						testResource2: 20,
					},
				},
			},
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 0,
		},
		"state amount is less than goal": {
			curState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 20.0,
		},
		"state amount is greater than goal": {
			curState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {},
				},
			},
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 20.0,
		},
		"agent inventory is weighted in goal": {
			curState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 50,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 20,
					},
				},
			},
			goalState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 70,
					},
				},
			},
			agent:            testAgent,
			expectedDistance: 2.0, // 20 from location, minus 18 from agent
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := heuristic(tc.curState, tc.goalState, tc.agent)
			assert.Equal(t, tc.expectedDistance, output)
		})
	}
}

func TestActions_GetSuccessors(t *testing.T) {
	mockAction1 := &mockAction{}
	mockAction2 := &mockAction{}
	mockAction3 := &mockAction{}

	mockNullAction1 := &mockNullAction{}

	expectedEndState := &State{
		Locations: map[*Location]Inventory{
			testLocation: {
				testResource: 1,
			},
		},
		Agents: map[*Agent]Inventory{},
	}

	type testCase struct {
		actions            Actions
		startState         *State
		agent              *Agent
		expectedSuccessors []SuccessorState
	}

	tests := map[string]testCase{
		"get successors with one action": {
			actions: Actions{
				mockAction1,
			},
			startState: &State{},
			agent:      testAgent,
			expectedSuccessors: []SuccessorState{
				{
					seq:      mockAction1,
					newState: expectedEndState,
				},
			},
		},
		"get successors with multiple actions": {
			actions: Actions{
				mockAction1,
				mockAction2,
				mockAction3,
			},
			startState: &State{},
			agent:      testAgent,
			expectedSuccessors: []SuccessorState{
				{
					seq:      mockAction1,
					newState: expectedEndState,
				},
				{
					seq:      mockAction2,
					newState: expectedEndState,
				},
				{
					seq:      mockAction3,
					newState: expectedEndState,
				},
			},
		},
		"get successors with some null actions": {
			actions: Actions{
				mockNullAction1,
				mockAction1,
				mockAction2,
			},
			startState: &State{},
			agent:      testAgent,
			expectedSuccessors: []SuccessorState{
				{
					seq:      mockAction1,
					newState: expectedEndState,
				},
				{
					seq:      mockAction2,
					newState: expectedEndState,
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.actions.getSuccessors(tc.startState, tc.agent)
			assert.Equal(t, tc.expectedSuccessors, output)
		})
	}
}
