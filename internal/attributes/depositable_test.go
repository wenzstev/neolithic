package attributes

import (
	"testing"

	"Neolithic/internal/actions"
	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
)

// TestCanDeposit_NeedsLocation verifies the NeedsLocation method always returns false.
func TestCanDeposit_NeedsLocation(t *testing.T) {
	cd := CanDepositTo{}
	// Assert that NeedsLocation returns false using direct package call
	assert.False(t, cd.NeedsLocation(), "CanDepositTo.NeedsLocation() should return false")
}

// TestCanDeposit_NeedsResource verifies the NeedsResource method always returns true.
func TestCanDeposit_NeedsResource(t *testing.T) {
	cd := CanDepositTo{}
	// Assert that NeedsResource returns true using direct package call
	assert.True(t, cd.NeedsResource(), "CanDepositTo.NeedsResource() should return true")
}

// TestCanDeposit_CreateAction verifies the CreateAction method returns the correct Deposit action.
func TestCanDeposit_CreateAction(t *testing.T) {
	// Mock Location and Resource for testing
	mockLocation := &core.Location{}
	mockResource := &core.Resource{Name: "Wood"}
	mockEmptyResource := &core.Resource{} // Test case with nil-like resource

	// Define test cases
	testCases := map[string]struct {
		canDeposit CanDepositTo
		params     core.CreateActionParams
		wantAction *actions.Deposit
	}{
		"Standard deposit": {
			canDeposit: CanDepositTo{
				Amount:   10,
				Cost:     0.5,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Deposit{
				DepResource:    mockResource,
				Amount:         10,
				ActionLocation: mockLocation,
				ActionCost:     0.5,
			},
		},
		"Zero amount deposit": {
			canDeposit: CanDepositTo{
				Amount:   0,
				Cost:     0.1,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Deposit{
				DepResource:    mockResource,
				Amount:         0,
				ActionLocation: mockLocation,
				ActionCost:     0.1,
			},
		},
		"Zero cost deposit": {
			canDeposit: CanDepositTo{
				Amount:   5,
				Cost:     0.0,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Deposit{
				DepResource:    mockResource,
				Amount:         5,
				ActionLocation: mockLocation,
				ActionCost:     0.0,
			},
		},
		"Empty resource deposit": {
			canDeposit: CanDepositTo{
				Amount:   2,
				Cost:     1.0,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockEmptyResource,
			},
			wantAction: &actions.Deposit{
				DepResource:    mockEmptyResource,
				Amount:         2,
				ActionLocation: mockLocation,
				ActionCost:     1.0,
			},
		},
		"Nil location in CanDepositTo": {
			canDeposit: CanDepositTo{
				Amount:   3,
				Cost:     0.2,
				Location: nil,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Deposit{
				DepResource:    mockResource,
				Amount:         3,
				ActionLocation: nil,
				ActionCost:     0.2,
			},
		},
		"Nil resource in Params": {
			canDeposit: CanDepositTo{
				Amount:   4,
				Cost:     0.3,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: nil,
			},
			wantAction: &actions.Deposit{
				DepResource:    nil,
				Amount:         4,
				ActionLocation: mockLocation,
				ActionCost:     0.3,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotAction := tc.canDeposit.CreateAction(tc.params)

			// Assert that the returned action is of the expected type (*actions.Deposit)
			// assert.IsType checks if the types match and handles the nil case gracefully.
			assert.IsType(t, tc.wantAction, gotAction, "Returned action should be of type *actions.Deposit")

			// Assert that the returned action is deeply equal to the expected action
			// assert.Equal performs a deep comparison for structs and pointers.
			assert.Equal(t, tc.wantAction, gotAction, "Returned action does not match expected action")
		})
	}
}
