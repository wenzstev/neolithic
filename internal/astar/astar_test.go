package astar

import (
	"errors"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// testError is used to test various error states in the search
var testError = errors.New("test error")

type dummyNode struct {
	name           string
	neighbors      []*dummyNode
	heuristicError error
	idError        error
	successorError error
}

var _ Node = (*dummyNode)(nil)

func (d *dummyNode) Heuristic(goal Node) (float64, error) {
	if d.heuristicError != nil {
		return 0, d.heuristicError
	}
	goalDummy, ok := goal.(*dummyNode)
	if !ok {
		return 0, errors.New("goal is not a dummyNode")
	}
	if d.name == goalDummy.name {
		return 0, nil
	}
	return 1, nil
}

func (d *dummyNode) ID() (string, error) {
	if d.idError != nil {
		return "", d.idError
	}
	return d.name, nil
}

func (d *dummyNode) Cost(prev Node) float64 {
	return 1
}

func (d *dummyNode) GetSuccessors() ([]Node, error) {
	if d.successorError != nil {
		return nil, d.successorError
	}
	successors := make([]Node, len(d.neighbors))
	for i, neighbor := range d.neighbors {
		successors[i] = neighbor
	}
	return successors, nil
}

func TestSearchState_RunIterations(t *testing.T) {
	type testCase struct {
		setupFunc                  func() (*dummyNode, *dummyNode, []Node)
		numRuns                    int
		expectedStartError         error
		expectedRunIterationsError error
		expectedCost               float64
		expectedIsBest             bool
	}

	testCases := map[string]testCase{
		"simple path": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B"}
				C := &dummyNode{name: "C"}

				A.neighbors = []*dummyNode{B}
				B.neighbors = []*dummyNode{C}

				path := []Node{A, B, C}

				return A, C, path
			},
			numRuns:                    10,
			expectedStartError:         nil,
			expectedRunIterationsError: nil,
			expectedCost:               2,
			expectedIsBest:             true,
		},
		"more complex path": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B"}
				C := &dummyNode{name: "C"}
				D := &dummyNode{name: "D"}
				E := &dummyNode{name: "E"}

				A.neighbors = []*dummyNode{B, C}
				B.neighbors = []*dummyNode{D}
				D.neighbors = []*dummyNode{E}
				C.neighbors = []*dummyNode{D, B}

				path := []Node{A, B, D, E}
				return A, E, path
			},
			numRuns:                    10,
			expectedStartError:         nil,
			expectedRunIterationsError: nil,
			expectedCost:               3,
			expectedIsBest:             true,
		},
		"fail due to heuristic error": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B", heuristicError: testError}
				C := &dummyNode{name: "C"}

				A.neighbors = []*dummyNode{B}
				B.neighbors = []*dummyNode{C}

				return A, C, nil
			},
			numRuns:                    10,
			expectedStartError:         nil,
			expectedRunIterationsError: testError,
			expectedCost:               0,
			expectedIsBest:             false,
		},
		"fail due to id error": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A", idError: testError}
				B := &dummyNode{name: "B"}
				C := &dummyNode{name: "C"}

				A.neighbors = []*dummyNode{B}
				B.neighbors = []*dummyNode{C}

				return A, C, nil
			},
			numRuns:                    10,
			expectedStartError:         testError,
			expectedRunIterationsError: nil,
			expectedCost:               0,
			expectedIsBest:             false,
		},
		"fail due to successors error": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B", successorError: testError}
				C := &dummyNode{name: "C"}

				A.neighbors = []*dummyNode{B}
				B.neighbors = []*dummyNode{C}

				return A, C, nil
			},
			numRuns:                    10,
			expectedStartError:         nil,
			expectedRunIterationsError: testError,
			expectedCost:               0,
			expectedIsBest:             false,
		},
		"fail due to no path": {
			setupFunc: func() (*dummyNode, *dummyNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B"}
				C := &dummyNode{name: "C"}
				D := &dummyNode{name: "D"}
				E := &dummyNode{name: "E"}

				A.neighbors = []*dummyNode{B, C}
				B.neighbors = []*dummyNode{D}
				C.neighbors = []*dummyNode{D, B}

				return A, E, nil
			},
			numRuns:                    10,
			expectedStartError:         nil,
			expectedRunIterationsError: ErrNoPath,
			expectedCost:               0,
			expectedIsBest:             false,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			startNode, endNode, expectedPath := tc.setupFunc()

			search, err := NewSearch(startNode, endNode)
			if tc.expectedStartError != nil {
				assert.Equal(t, tc.expectedStartError, err)
				return
			}
			assert.NoError(t, err)

			err = search.RunIterations(tc.numRuns)
			if tc.expectedRunIterationsError != nil {
				assert.Equal(t, tc.expectedRunIterationsError, err)
				return
			}
			assert.NoError(t, err)

			assert.Equal(t, tc.expectedCost, search.BestCost)
			assert.Equal(t, expectedPath, search.CurrentBestPath())
			assert.Equal(t, tc.expectedIsBest, search.FoundBest)
		})
	}
}

func TestNewSearch(t *testing.T) {
	type testCase struct {
		start, end          *dummyNode
		expectedSearchState *SearchState
	}

	start := &dummyNode{name: "start"}
	end := &dummyNode{name: "end"}

	testCases := map[string]testCase{
		"can initialize new search": {
			start: start,
			end:   end,
			expectedSearchState: &SearchState{
				Start:    start,
				Goal:     end,
				BestCost: math.Inf(1),
				openSet: &PriorityQueue{
					&searchNode{
						gCost:     0,
						hCost:     1,
						parent:    nil,
						nodeState: start,
					},
				},
				closedSet: map[string]bool{},
				openSetMap: map[string]*searchNode{
					"start": {
						gCost:     0,
						hCost:     1,
						parent:    nil,
						nodeState: start,
					},
				},
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			search, err := NewSearch(tc.start, tc.end)
			assert.NoError(t, err)
			assert.Equal(t, tc.expectedSearchState, search)
		})
	}
}

func TestSearchState_CurrentBest(t *testing.T) {
	type testCase struct {
		setupFunc func() (*searchNode, []Node)
	}

	testCases := map[string]testCase{
		"can get basic path": {
			setupFunc: func() (*searchNode, []Node) {
				A := &dummyNode{name: "A"}
				B := &dummyNode{name: "B"}
				C := &dummyNode{name: "C"}

				searchNodeA := &searchNode{nodeState: A}
				searchNodeB := &searchNode{nodeState: B, parent: searchNodeA}
				searchNodeC := &searchNode{nodeState: C, parent: searchNodeB}

				return searchNodeC, []Node{A, B, C}
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			start, expected := tc.setupFunc()
			searchState := &SearchState{
				bestSolution: start,
			}

			assert.Equal(t, expected, searchState.CurrentBestPath())
		})
	}
}
