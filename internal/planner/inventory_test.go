package planner

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestInventory_Copy(t *testing.T) {
	type testCase struct {
		inventory Inventory
	}

	tests := map[string]testCase{
		"basic copy": {
			inventory: Inventory{
				testResource: 10,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			newInventory := tc.inventory.Copy()
			assert.True(t, reflect.DeepEqual(newInventory, tc.inventory))
			assert.False(t, &tc.inventory == &newInventory)
		})
	}
}

func TestInventory_String(t *testing.T) {
	type testCase struct {
		inventory Inventory
		expected  string
	}

	tests := map[string]testCase{
		"can produce expected string": {
			inventory: Inventory{
				testResource: 10,
			},
			expected: "      testResource: 10\n",
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			output := tc.inventory.String()
			assert.Equal(t, tc.expected, output)
		})
	}
}
