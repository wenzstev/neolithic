package agent

import (
	"testing"

	"Neolithic/internal/core"
)

func TestNewCoordPath(t *testing.T) {
	tests := []struct {
		name     string
		coords   []core.Coord
		expected *CoordPath
	}{
		{
			name:   "empty Path",
			coords: []core.Coord{},
			expected: &CoordPath{
				coords: []core.Coord{},
				index:  1,
			},
		},
		{
			name:   "single coordinate",
			coords: []core.Coord{{X: 1, Y: 2}},
			expected: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}},
				index:  1,
			},
		},
		{
			name:   "multiple coordinates",
			coords: []core.Coord{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
			expected: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
				index:  1,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewCoordPath(tt.coords)
			if got.index != tt.expected.index {
				t.Errorf("NewCoordPath() index = %v, want %v", got.index, tt.expected.index)
			}
			if len(got.coords) != len(tt.expected.coords) {
				t.Errorf("NewCoordPath() coords length = %v, want %v", len(got.coords), len(tt.expected.coords))
			}
			for i := range got.coords {
				if got.coords[i] != tt.expected.coords[i] {
					t.Errorf("NewCoordPath() coords[%d] = %v, want %v", i, got.coords[i], tt.expected.coords[i])
				}
			}
		})
	}
}

func TestCoordPath_NextCoord(t *testing.T) {
	tests := []struct {
		name     string
		path     *CoordPath
		expected []core.Coord
		wantErr  bool
	}{
		{
			name: "valid Path with multiple coordinates",
			path: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
				index:  1,
			},
			expected: []core.Coord{{X: 3, Y: 4}, {X: 5, Y: 6}},
			wantErr:  false,
		},
		{
			name: "empty Path",
			path: &CoordPath{
				coords: []core.Coord{},
				index:  1,
			},
			expected: []core.Coord{},
			wantErr:  true,
		},
		{
			name: "single coordinate Path",
			path: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}},
				index:  1,
			},
			expected: []core.Coord{},
			wantErr:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for i, expected := range tt.expected {
				got := tt.path.NextCoord()
				if got != expected {
					t.Errorf("NextCoord() = %v, want %v", got, expected)
				}
				if tt.path.index != i+2 { // +2 because we start at index 1 and have moved i+1 times
					t.Errorf("index after NextCoord() = %v, want %v", tt.path.index, i+2)
				}
			}

			// Test that calling NextCoord after Path is complete panics
			if !tt.wantErr {
				defer func() {
					if r := recover(); r == nil {
						t.Error("NextCoord() should have panicked but didn't")
					}
				}()
				tt.path.NextCoord()
			}
		})
	}
}

func TestCoordPath_IsComplete(t *testing.T) {
	tests := []struct {
		name     string
		path     *CoordPath
		expected bool
	}{
		{
			name: "empty Path",
			path: &CoordPath{
				coords: []core.Coord{},
				index:  1,
			},
			expected: true,
		},
		{
			name: "single coordinate Path",
			path: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}},
				index:  1,
			},
			expected: true,
		},
		{
			name: "multiple coordinates Path - not complete",
			path: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
				index:  1,
			},
			expected: false,
		},
		{
			name: "multiple coordinates Path - complete",
			path: &CoordPath{
				coords: []core.Coord{{X: 1, Y: 2}, {X: 3, Y: 4}, {X: 5, Y: 6}},
				index:  3,
			},
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.path.IsComplete()
			if got != tt.expected {
				t.Errorf("Complete() = %v, want %v", got, tt.expected)
			}
		})
	}
}
