package core

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestWorldStateCopy(t *testing.T) {
	// Create a WorldState with one location and one agent.
	orig := &WorldState{
		Locations: map[string]Location{
			// using the Location’s Name field as the key
			testLocation.Name: testLocation,
		},
		Agents: map[string]Agent{
			// using Agent’s Name() as the key
			testAgent.Name(): testAgent,
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
				Locations: map[string]Location{
					testLocation.Name: testLocation,
				},
				Agents: map[string]Agent{},
			},
			expected: "WorldState :\n Locations:\n" + testLocation.String() + "\n Agents:\n",
		},
		"string with agent": {
			ws: &WorldState{
				Locations: map[string]Location{},
				Agents: map[string]Agent{
					testAgent.Name(): testAgent,
				},
			},
			expected: "WorldState :\n Locations:\n\n Agents:\n" + testAgent.String(),
		},
		"string with both agent and location": {
			ws: &WorldState{
				Locations: map[string]Location{
					testLocation.Name: testLocation,
				},
				Agents: map[string]Agent{
					testAgent.Name(): testAgent,
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
				Locations: map[string]Location{
					testLocation.Name: testLocation,
				},
				Agents: map[string]Agent{},
			},
			// Expected hash values depend on your gob encoding and sort order.
			// Update these expected values as needed.
			expected: "b761a3f00454a2955d59ba88ac4a7d7df1d79ee4af2a46a828711e3c056ca831",
		},
		"id with multiple locations": {
			ws: &WorldState{
				Locations: map[string]Location{
					testLocation.Name:  testLocation,
					testLocation2.Name: testLocation2,
				},
				Agents: map[string]Agent{},
			},
			expected: "a45402744f3a265d2f4519e878780ff959a454033f8b1850a135feb18dbc8071",
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
