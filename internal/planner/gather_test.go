package planner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGather_Perform(t *testing.T) {
	testGather := &Gather{
		resource: testResource,
		amount:   5,
		location: testLocation,
		cost:     1,
	}

	testTool := &Resource{Name: "testTool"}

	testGatherRequires := &Gather{
		requires: testTool,
		resource: testResource,
		amount:   5,
		location: testLocation,
	}

	type testCase struct {
		testGather       *Gather
		testAgent        Agent
		startState       *State
		expectedEndState *State
	}

	testCases := map[string]testCase{
		"can do basic gather": {
			testGather: testGather,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 10,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 5,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testResource: 5,
					},
				},
			},
		},
		"gather partially succeeds": {
			testGather: testGather,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 2,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testResource: 2,
					},
				},
			},
		},
		"gather succeeds with tool": {
			testGather: testGatherRequires,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 10,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testTool: 1,
					},
				},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 5,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testResource: 5,
						testTool:     1,
					},
				},
			},
		},
		"gather fails, no resource in location": {
			testGather: testGather,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[Agent]Inventory{
					testAgent: {},
				},
			},
			expectedEndState: nil,
		},
		"gather fails, required tool not present": {
			testGather: testGatherRequires,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 10,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testTool: 1,
					},
				},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 5,
					},
				},
				Agents: map[Agent]Inventory{
					testAgent: {
						testTool:     1,
						testResource: 5,
					},
				},
			},
		},
		"gather fails, agent not in State": {
			testGather: testGather,
			testAgent:  testAgent,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 10,
					},
				},
			},
			expectedEndState: nil,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			endState := tc.testGather.Perform(tc.startState, tc.testAgent)
			assert.Equal(t, tc.expectedEndState, endState)
		})
	}
}
