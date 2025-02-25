package planner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testMockAction = &mockAction{}

func TestCompoundAction_Perform(t *testing.T) {
	type testCase struct {
		compoundAction   *CompoundAction
		startState       *State
		expectedEndState *State
	}

	tests := map[string]testCase{
		"compound Action of one Action": {
			compoundAction: &CompoundAction{
				testMockAction,
			},
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[Agent]Inventory{},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 1,
					},
				},
				Agents: map[Agent]Inventory{},
			},
		},
		"compound Action of two actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
			},
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[Agent]Inventory{},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 2,
					},
				},
				Agents: map[Agent]Inventory{},
			},
		},
		"compound Action of three actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
				testMockAction,
			},
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[Agent]Inventory{},
			},
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 3,
					},
				},
				Agents: map[Agent]Inventory{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			endState := tc.compoundAction.Perform(tc.startState, testAgent)
			assert.Equal(t, tc.expectedEndState, endState)
		})
	}
}

func TestCompoundAction_Cost(t *testing.T) {
	type testCase struct {
		compoundAction *CompoundAction
		expectedCost   float64
	}

	tests := map[string]testCase{
		"compound Action of one Action": {
			compoundAction: &CompoundAction{
				testMockAction,
			},
			expectedCost: 10.0,
		},
		"compound Action of two actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
			},
			expectedCost: 20.0,
		},
		"compound Action of three actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
				testMockAction,
			},
			expectedCost: 30.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCost, tc.compoundAction.Cost(testAgent))
		})
	}
}

func TestCompoundAction_Description(t *testing.T) {
	type testCase struct {
		compoundAction *CompoundAction
		expectedDesc   string
	}

	tests := map[string]testCase{
		"compound Action of one Action": {
			compoundAction: &CompoundAction{
				testMockAction,
			},
			expectedDesc: "Sequence:\n  a mock Action\n",
		},
		"compound Action of two actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
			},
			expectedDesc: "Sequence:\n  a mock Action\n  a mock Action\n",
		},
		"compound Action of three actions": {
			compoundAction: &CompoundAction{
				testMockAction,
				testMockAction,
				testMockAction,
			},
			expectedDesc: "Sequence:\n  a mock Action\n  a mock Action\n  a mock Action\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedDesc, tc.compoundAction.Description())
		})
	}
}
