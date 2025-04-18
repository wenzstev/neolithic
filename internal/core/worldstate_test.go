package core

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldStateCopy(t *testing.T) {
	// Create a WorldState with one location and one agent.
	orig := &WorldState{
		Locations: []Location{
			testLocation,
		},
		Agents: []Agent{
			testAgent,
		},
	}

	copied := orig.DeepCopy()
	// The values should be deeply equal...
	assert.True(t, reflect.DeepEqual(orig, copied), "expected states to have same values")
	// ...but the copied pointer should be different.
	assert.False(t, orig == copied, "expected copied WorldState to have a different memory address")
}

func TestWorldStateString(t *testing.T) {
	tests := map[string]struct {
		ws       *WorldState
		expected string
	}{
		"string with location": {
			ws: &WorldState{
				Locations: []Location{
					testLocation,
				},
				Agents: []Agent{},
			},
			expected: "WorldState :\n Locations:\n" + testLocation.String() + "\n Agents:\n",
		},
		"string with agent": {
			ws: &WorldState{
				Locations: []Location{},
				Agents: []Agent{
					testAgent,
				},
			},
			expected: "WorldState :\n Locations:\n\n Agents:\n" + testAgent.String(),
		},
		"string with both agent and location": {
			ws: &WorldState{
				Locations: []Location{
					testLocation,
				},
				Agents: []Agent{
					testAgent,
				},
			},
			expected: "WorldState :\n Locations:\n" + testLocation.String() + "\n Agents:\n" + testAgent.String(),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.ws.String()
			assert.Equal(t, tc.expected, output, "WorldState strings do not match")
		})
	}
}

func TestWorldStateID(t *testing.T) {
	tests := map[string]struct {
		ws       *WorldState
		expected string
	}{
		"id with location": {
			ws: &WorldState{
				Locations: []Location{
					testLocation,
				},
				Agents: []Agent{},
			},
			// Expected hash values depend on your gob encoding and sort order.
			// Update these expected values as needed.
			expected: "7d5db8bb28ee0afa",
		},
		"id with multiple locations": {
			ws: &WorldState{
				Locations: []Location{
					testLocation,
					testLocation2,
				},
				Agents: []Agent{},
			},
			expected: "ecbc75bfdcf137ae",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := tc.ws.ID()
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, id)
		})
	}
}
