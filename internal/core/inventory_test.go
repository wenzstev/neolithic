package core

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

var testResource = &Resource{
	Name: "testResource",
}

func TestInventory_AdjustAmount(t *testing.T) {
	type testCase struct {
		initialState  *inventory
		resource      *Resource
		amount        int
		expectedState *inventory
	}

	tests := map[string]testCase{
		"add new resource": {
			initialState: &inventory{},
			resource:     &Resource{Name: "newResource"},
			amount:       5,
			expectedState: &inventory{
				{Resource: &Resource{Name: "newResource"}, Amount: 5},
			},
		},
		"add to existing resource": {
			initialState: &inventory{
				{Resource: &Resource{Name: "existingResource"}, Amount: 10},
			},
			resource: &Resource{Name: "existingResource"},
			amount:   5,
			expectedState: &inventory{
				{Resource: &Resource{Name: "existingResource"}, Amount: 15},
			},
		},
		"remove from existing resource": {
			initialState: &inventory{
				{Resource: &Resource{Name: "existingResource"}, Amount: 10},
			},
			resource: &Resource{Name: "existingResource"},
			amount:   -5,
			expectedState: &inventory{
				{Resource: &Resource{Name: "existingResource"}, Amount: 5},
			},
		},
		"remove resource completely": {
			initialState: &inventory{
				{Resource: &Resource{Name: "existingResource"}, Amount: 10},
			},
			resource:      &Resource{Name: "existingResource"},
			amount:        -10,
			expectedState: &inventory{},
		},
		"attempt to add negative amount": {
			initialState:  &inventory{},
			resource:      &Resource{Name: "newResource"},
			amount:        -5,
			expectedState: &inventory{},
		},
		"maintain sorted order": {
			initialState: &inventory{
				{Resource: &Resource{Name: "aResource"}, Amount: 5},
				{Resource: &Resource{Name: "cResource"}, Amount: 5},
			},
			resource: &Resource{Name: "bResource"},
			amount:   5,
			expectedState: &inventory{
				{Resource: &Resource{Name: "aResource"}, Amount: 5},
				{Resource: &Resource{Name: "bResource"}, Amount: 5},
				{Resource: &Resource{Name: "cResource"}, Amount: 5},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.initialState.AdjustAmount(tc.resource, tc.amount)
			assert.True(t, reflect.DeepEqual(tc.initialState, tc.expectedState))
		})
	}
}

func TestInventory_DeepCopy(t *testing.T) {
	type testCase struct {
		name     string
		initial  *inventory
		mutateOp func(*inventory)
	}

	tests := map[string]testCase{
		"copy empty inventory": {
			name:    "empty inventory",
			initial: &inventory{},
			mutateOp: func(i *inventory) {
				i.AdjustAmount(&Resource{Name: "test"}, 5)
			},
		},
		"copy multiple resources": {
			name: "multiple resources",
			initial: &inventory{
				{Resource: &Resource{Name: "resource1"}, Amount: 10},
				{Resource: &Resource{Name: "resource2"}, Amount: 20},
				{Resource: &Resource{Name: "resource3"}, Amount: 30},
			},
			mutateOp: func(i *inventory) {
				i.AdjustAmount(&Resource{Name: "resource1"}, 5)
			},
		},
		"large inventory": {
			name: "large inventory",
			initial: &inventory{
				{Resource: &Resource{Name: "r1"}, Amount: 1},
				{Resource: &Resource{Name: "r2"}, Amount: 2},
				{Resource: &Resource{Name: "r3"}, Amount: 3},
				{Resource: &Resource{Name: "r4"}, Amount: 4},
				{Resource: &Resource{Name: "r5"}, Amount: 5},
			},
			mutateOp: func(i *inventory) {
				i.AdjustAmount(&Resource{Name: "r6"}, 6)
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			original := tc.initial.DeepCopy()
			tc.mutateOp(tc.initial)

			assert.False(t, reflect.DeepEqual(original, tc.initial), "original and mutated inventory should not be equal")

			copied := original.DeepCopy()
			assert.True(t, reflect.DeepEqual(original, copied), "copy should equal original")
			assert.NotSame(t, original, copied, "copy should not have same pointer")

		})
	}
}

func TestInventory_String(t *testing.T) {
	type testCase struct {
		inventory Inventory
		amount    int
		expected  string
	}

	tests := map[string]testCase{
		"can produce expected string": {
			inventory: &inventory{},
			amount:    10,
			expected:  "  testResource: 10\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			tc.inventory.AdjustAmount(testResource, tc.amount)
			output := tc.inventory.String()
			assert.Equal(t, tc.expected, output)
		})
	}
}
