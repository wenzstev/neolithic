package astar

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

var testNode = &searchNode{}

var greaterNode = &searchNode{
	hCost: 5,
	gCost: 5,
}

var lessNode = &searchNode{
	hCost: 1,
	gCost: 1,
}

func TestPriorityQueue_Len(t *testing.T) {
	type testCase struct {
		pq          PriorityQueue
		expectedVal int
	}

	tests := map[string]testCase{
		"length of 3": {
			pq: PriorityQueue{
				testNode,
				testNode,
				testNode,
			},
			expectedVal: 3,
		},
		"length of 1": {
			pq: PriorityQueue{
				testNode,
			},
			expectedVal: 1,
		},
		"length of 0": {
			pq:          PriorityQueue{},
			expectedVal: 0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedVal, tc.pq.Len())
		})
	}
}

func TestPriorityQueue_Less(t *testing.T) {
	type testCase struct {
		iNode       *searchNode
		jNode       *searchNode
		expectedVal bool
	}

	tests := map[string]testCase{
		"i greater than j": {
			iNode:       greaterNode,
			jNode:       lessNode,
			expectedVal: false,
		},
		"j greater than i": {
			iNode:       lessNode,
			jNode:       greaterNode,
			expectedVal: true,
		},
		"i and j equal": {
			iNode:       lessNode,
			jNode:       lessNode,
			expectedVal: false,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := PriorityQueue{
				tc.iNode,
				tc.jNode,
			}

			assert.Equal(t, tc.expectedVal, pq.Less(0, 1))

		})
	}
}

func TestPriorityQueue_Swap(t *testing.T) {
	type testCase struct {
		iNode *searchNode
		jNode *searchNode
	}

	tests := map[string]testCase{
		"can swap": {
			iNode: greaterNode,
			jNode: lessNode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			pq := PriorityQueue{
				tc.iNode,
				tc.jNode,
			}
			pq.Swap(0, 1)
			assert.Equal(t, pq[0], tc.jNode)
			assert.Equal(t, pq[1], tc.iNode)
		})
	}
}

func TestPriorityQueue_Push(t *testing.T) {
	type testCase struct {
		pq      PriorityQueue
		newNode *searchNode
	}

	tests := map[string]testCase{
		"can push": {
			pq: PriorityQueue{
				testNode,
			},
			newNode: greaterNode,
		},
		"can push when empty": {
			pq:      PriorityQueue{},
			newNode: greaterNode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			originalLen := len(tc.pq)
			tc.pq.Push(tc.newNode)
			assert.Equal(t, len(tc.pq), originalLen+1)
			assert.Equal(t, tc.newNode, tc.pq[originalLen])
		})
	}
}

func TestPriorityQueue_Pop(t *testing.T) {
	type testCase struct {
		pq        PriorityQueue
		expectedN *searchNode
	}

	tests := map[string]testCase{
		"can pop": {
			pq: PriorityQueue{
				testNode,
				lessNode,
				greaterNode,
			},
			expectedN: greaterNode,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.expectedN, tc.pq.Pop())
		})
	}
}
