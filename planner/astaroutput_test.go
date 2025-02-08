package planner

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAStarOutput_String(t *testing.T) {
	type testCase struct {
		astarOutput *AStarOutput
		expected    string
	}

	testCases := map[string]testCase{
		"has expected string": {
			astarOutput: &AStarOutput{
				actions: []Action{
					testMockAction,
					testMockAction,
				},
				totalCost: 30.0,
				expectedState: &State{
					Locations: map[*Location]Inventory{
						testLocation: {
							testResource: 10,
						},
					},
					Agents: map[*Agent]Inventory{
						testAgent: {
							testResource: 20,
						},
					},
				},
			},
			expected: "Total Cost: 30.000000\nActions: \na mock action\na mock action\n",
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expected, tc.astarOutput.String())
		})
	}
}
