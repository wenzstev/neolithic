package planner

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStateCopy(t *testing.T) {
	type testCase struct {
		startState *State
	}

	tests := map[string]testCase{
		"basic copy": {
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[*Agent]Inventory{},
			},
		},
		"basic copy with agent": {
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 1,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			copiedState := tc.startState.Copy()
			assert.True(t, reflect.DeepEqual(tc.startState, copiedState), "expected states to have same values")
			assert.False(t, tc.startState == copiedState, "expected copied state to have different memory address")
		})
	}

}

func TestStateString(t *testing.T) {
	type testCase struct {
		testState *State
		expected  string
	}

	tests := map[string]testCase{
		"string with location": {
			testState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n   testLocation: \n      testResource: 1\n  Agents:\n",
		},
		"string with agent": {
			testState: &State{
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n  Agents:\n   testAgent: \n      testResource: 1\n",
		},
		"string with agent and location": {
			testState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n   testLocation: \n      testResource: 1\n  Agents:\n   testAgent: \n      testResource: 1\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.testState.String()
			assert.Equal(t, tc.expected, output, "State strings do not match")
		})
	}
}

func TestStateID(t *testing.T) {
	type testCase struct {
		testState *State
		expected  string
	}

	tests := map[string]testCase{
		"id with location": {
			testState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
			},
			expected: "2e022d8ce9a99d6071b691bd04668dc027a89f059abec00e1300a7096d98ced7",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := tc.testState.ID()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, id)
		})
	}
}
