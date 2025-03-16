package core

import (
	"Neolithic/internal/planner"
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
					planner.testLocation: {
						planner.testResource: 1,
					},
				},
				Agents: map[Agent]Inventory{},
			},
		},
		"basic copy with Agent": {
			startState: &State{
				Locations: map[*Location]Inventory{
					planner.testLocation: {
						planner.testResource: 1,
					},
				},
				Agents: map[Agent]Inventory{
					planner.testAgent: {
						planner.testResource: 1,
					},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			copiedState := tc.startState.Copy()
			assert.True(t, reflect.DeepEqual(tc.startState, copiedState), "expected states to have same values")
			assert.False(t, tc.startState == copiedState, "expected copied State to have different memory address")
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
					planner.testLocation: {
						planner.testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n   testLocation: \n      testResource: 1\n  Agents:\n",
		},
		"string with Agent": {
			testState: &State{
				Agents: map[Agent]Inventory{
					planner.testAgent: {
						planner.testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n  Agents:\n   testAgent: \n      testResource: 1\n",
		},
		"string with Agent and location": {
			testState: &State{
				Locations: map[*Location]Inventory{
					planner.testLocation: {
						planner.testResource: 1,
					},
				},
				Agents: map[Agent]Inventory{
					planner.testAgent: {
						planner.testResource: 1,
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
					planner.testLocation: {
						planner.testResource: 1,
					},
				},
			},
			expected: "b761a3f00454a2955d59ba88ac4a7d7df1d79ee4af2a46a828711e3c056ca831",
		},
		"id with multiple locations": {
			testState: &State{
				Locations: map[*Location]Inventory{
					planner.testLocation: {
						planner.testResource: 1,
					},
					planner.testLocation2: {
						planner.testResource: 1,
					},
				},
			},
			expected: "a45402744f3a265d2f4519e878780ff959a454033f8b1850a135feb18dbc8071",
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
