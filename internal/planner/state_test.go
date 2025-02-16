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
		"basic copy with Agent": {
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
		"string with Agent": {
			testState: &State{
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 1,
					},
				},
			},
			expected: "State: \n  Locations: \n  Agents:\n   testAgent: \n      testResource: 1\n",
		},
		"string with Agent and location": {
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
			expected: "64824b7b2deede8f76a220c0ce4f455743f07a120872b0bfe7643d03d62ea8bd",
		},
		"id with multiple locations": {
			testState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
					testLocation2: {
						testResource: 1,
					},
				},
			},
			expected: "275be152e128ab24878fa164422441cf113fc78f34f5256f6a1b2565daa35f69",
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
