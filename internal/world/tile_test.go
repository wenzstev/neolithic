package world

import (
	"testing"

	"Neolithic/internal/astar"
	"Neolithic/internal/grid"
	"github.com/stretchr/testify/assert"
)

func testMakeTile(X, Y int, grid *grid.Grid) (grid.Tile, error) {
	return &Tile{
		X:    X,
		Y:    Y,
		grid: grid,
	}, nil
}

func TestTile_AStar(t *testing.T) {

	type testCase struct {
		startX, startY, endX, endY int
		gridSize                   int
		expectedPathString         []string
		expectedCost               float64
	}

	tests := map[string]testCase{
		"can find diagnoal path": {
			startX: 0, startY: 0, endX: 4, endY: 4,
			gridSize: 5,
			expectedPathString: []string{
				"0,0",
				"1,1",
				"2,2",
				"3,3",
				"4,4",
			},
			expectedCost: 5.6,
		},
		"can find straight path": {
			startX: 0, startY: 0, endX: 4, endY: 0,
			gridSize: 5,
			expectedPathString: []string{
				"0,0",
				"1,0",
				"2,0",
				"3,0",
				"4,0",
			},
			expectedCost: 4,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testGrid := grid.New(tc.gridSize, tc.gridSize, 1)
			err := testGrid.Initialize(testMakeTile)
			assert.NoError(t, err)

			start, ok := testGrid.Tiles[tc.startX][tc.startY].(*Tile)
			assert.True(t, ok)

			end, ok := testGrid.Tiles[tc.endX][tc.endY].(*Tile)
			assert.True(t, ok)

			search, err := astar.NewSearch(start, end)
			assert.NoError(t, err)

			err = search.RunIterations(10000)
			assert.NoError(t, err)

			var pathString []string
			nodeList := search.CurrentBestPath()
			for _, node := range nodeList {
				nodeID, err := node.ID()
				assert.NoError(t, err)
				pathString = append(pathString, nodeID)
			}
			assert.Equal(t, tc.expectedPathString, pathString)
			assert.Equal(t, tc.expectedCost, search.BestCost)

		})
	}
}

func TestTile_Heuristic(t *testing.T) {
	type testCase struct {
		startTile         *Tile
		endTile           *Tile
		expectedHeuristic float64
	}

	tests := map[string]testCase{
		"adjacent to goal": {
			startTile: &Tile{
				X: 0,
				Y: 0,
			},
			endTile: &Tile{
				X: 1,
				Y: 0,
			},
			expectedHeuristic: 1,
		},
		"is goal": {
			startTile: &Tile{
				X: 0,
				Y: 0,
			},
			endTile: &Tile{
				X: 0,
				Y: 0,
			},
			expectedHeuristic: 0,
		},
		"farther from goal": {
			startTile: &Tile{
				X: 0,
				Y: 0,
			},
			endTile: &Tile{
				X: 1,
				Y: 4,
			},
			expectedHeuristic: 4.123105625617661,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			hVal, err := tc.startTile.Heuristic(tc.endTile)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedHeuristic, hVal)
		})
	}
}

func TestTile_GetSuccessors(t *testing.T) {
	type testCase struct {
		X, Y               int
		gridSize           int
		expectedSuccessors []Tile
	}

	tests := map[string]testCase{
		"8 successors, middle tile": {
			X:        1,
			Y:        1,
			gridSize: 3,
			expectedSuccessors: []Tile{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 0, Y: 2},
				{X: 1, Y: 0},
				{X: 1, Y: 2},
				{X: 2, Y: 0},
				{X: 2, Y: 1},
				{X: 2, Y: 2},
			},
		},
		"5 successors, edge tile": {
			X:        1,
			Y:        0,
			gridSize: 3,
			expectedSuccessors: []Tile{
				{X: 0, Y: 0},
				{X: 0, Y: 1},
				{X: 1, Y: 1},
				{X: 2, Y: 0},
				{X: 2, Y: 1},
			},
		},
		"3 successors, corner tile": {
			X:        0,
			Y:        0,
			gridSize: 3,
			expectedSuccessors: []Tile{
				{X: 0, Y: 1},
				{X: 1, Y: 0},
				{X: 1, Y: 1},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			grid := grid.New(tc.gridSize, tc.gridSize, 1)
			err := grid.Initialize(testMakeTile)
			assert.NoError(t, err)

			testTile, ok := grid.Tiles[tc.X][tc.Y].(*Tile)
			assert.True(t, ok)

			successors, err := testTile.GetSuccessors()
			assert.NoError(t, err)

			var tileSlice []Tile
			for _, successor := range successors {
				tileSuccessor := *successor.(*Tile)
				tileSuccessor.grid = nil
				tileSlice = append(tileSlice, tileSuccessor)
			}

			assert.Equal(t, tc.expectedSuccessors, tileSlice)
		})
	}
}

func TestTile_ID(t *testing.T) {
	type testCase struct {
		testTile   *Tile
		expectedID string
	}

	tests := map[string]testCase{
		"1, 1": {
			testTile: &Tile{
				X: 1,
				Y: 1,
			},
			expectedID: "1,1",
		},
		"6, 2": {
			testTile: &Tile{
				X: 6,
				Y: 2,
			},
			expectedID: "6,2",
		},
		"423, 41": {
			testTile: &Tile{
				X: 423,
				Y: 41,
			},
			expectedID: "423,41",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			id, err := tc.testTile.ID()
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedID, id)
		})
	}
}
