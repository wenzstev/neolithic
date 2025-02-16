package planner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testDeposit = &Deposit{
	resource: testResource,
	amount:   10,
	location: testLocation,
	cost:     1.0,
}

func TestDeposit_Perform(t *testing.T) {
	type testCase struct {
		testDeposit      *Deposit
		startState       *State
		agent            *Agent
		expectedEndState *State
	}

	tests := map[string]testCase{
		"can do basic deposit": {
			testDeposit: testDeposit,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 10,
					},
				},
			},
			agent: testAgent,
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 10,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {},
				},
			},
		},
		"deposit fails, nothing in agent inventory": {
			testDeposit: testDeposit,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 0,
					},
				},
				Agents: map[*Agent]Inventory{},
			},
			agent:            testAgent,
			expectedEndState: nil,
		},
		"partial deposit success": {
			testDeposit: testDeposit,
			startState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 0,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {
						testResource: 5,
					},
				},
			},
			agent: testAgent,
			expectedEndState: &State{
				Locations: map[*Location]Inventory{
					testLocation: {
						testResource: 5,
					},
				},
				Agents: map[*Agent]Inventory{
					testAgent: {},
				},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			endState := tc.testDeposit.Perform(tc.startState, tc.agent)
			assert.Equal(t, tc.expectedEndState, endState)
		})
	}
}

func TestDeposit_Cost(t *testing.T) {
	type testCase struct {
		testDeposit  *Deposit
		testAgent    *Agent
		expectedCost float64
	}

	tests := map[string]testCase{
		"cost works": {
			testDeposit:  testDeposit,
			testAgent:    testAgent,
			expectedCost: 1.0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedCost, testDeposit.Cost(tc.testAgent))
		})
	}
}

func TestDeposit_String(t *testing.T) {
	type testCase struct {
		testDeposit    *Deposit
		expectedString string
	}

	tests := map[string]testCase{
		"basic deposit message": {
			testDeposit:    testDeposit,
			expectedString: "deposit 10 testResource at testLocation",
		},
		"deposit message with different amount": {
			testDeposit: &Deposit{
				resource: testResource,
				amount:   100,
				location: testLocation,
			},
			expectedString: "deposit 100 testResource at testLocation",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.testDeposit.Description()
			assert.Equal(t, tc.expectedString, output)
		})
	}
}
