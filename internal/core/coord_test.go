package core

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCoord_String(t *testing.T) {
	type testCase struct {
		c    Coord
		want string
	}

	tests := map[string]testCase{
		"positive coordinates": {
			c:    Coord{X: 5, Y: 10},
			want: "(5, 10)",
		},
		"negative coordinates": {
			c:    Coord{X: -3, Y: -7},
			want: "(-3, -7)",
		},
		"zero coordinates": {
			c:    Coord{X: 0, Y: 0},
			want: "(0, 0)",
		},
		"mixed coordinates": {
			c:    Coord{X: -5, Y: 8},
			want: "(-5, 8)",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.want, tc.c.String())
		})
	}
}

func TestCoord_IsWithin(t *testing.T) {
	type testCase struct {
		c        Coord
		other    Coord
		distance int
		want     bool
	}

	tests := map[string]testCase{
		"same point": {
			c:        Coord{X: 5, Y: 5},
			other:    Coord{X: 5, Y: 5},
			distance: 0,
			want:     true,
		},
		"within distance": {
			c:        Coord{X: 5, Y: 5},
			other:    Coord{X: 7, Y: 6},
			distance: 3,
			want:     true,
		},
		"at exact boundary": {
			c:        Coord{X: 5, Y: 5},
			other:    Coord{X: 10, Y: 5},
			distance: 5,
			want:     true,
		},
		"outside distance": {
			c:        Coord{X: 5, Y: 5},
			other:    Coord{X: 11, Y: 5},
			distance: 5,
			want:     false,
		},
		"negative coordinates within": {
			c:        Coord{X: -5, Y: -5},
			other:    Coord{X: -7, Y: -7},
			distance: 3,
			want:     true,
		},
		"negative coordinates outside": {
			c:        Coord{X: -5, Y: -5},
			other:    Coord{X: -10, Y: -10},
			distance: 4,
			want:     false,
		},
		"mixed signs within": {
			c:        Coord{X: -2, Y: 3},
			other:    Coord{X: 1, Y: 1},
			distance: 4,
			want:     true,
		},
		"diagonal at exact distance": {
			c:        Coord{X: 0, Y: 0},
			other:    Coord{X: 3, Y: 3},
			distance: 3,
			want:     true,
		},
		"negative distance (invalid but should behave correctly)": {
			c:        Coord{X: 5, Y: 5},
			other:    Coord{X: 5, Y: 5},
			distance: -1,
			want:     false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			result := tc.c.IsWithin(tc.other, tc.distance)
			assert.Equal(t, tc.want, result)
		})
	}
}
