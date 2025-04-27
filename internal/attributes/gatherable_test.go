package attributes

import (
	"testing"

	"Neolithic/internal/actions"
	"Neolithic/internal/core"

	"github.com/stretchr/testify/assert"
)

// TestCanGather_NeedsLocation verifies the NeedsLocation method always returns false.
func TestCanGather_NeedsLocation(t *testing.T) {
	cg := CanGather{}
	// Assert that NeedsLocation returns false using direct package call
	assert.False(t, cg.NeedsLocation(), "CanGather.NeedsLocation() should return false")
}

// TestCanGather_NeedsResource verifies the NeedsResource method always returns true.
func TestCanGather_NeedsResource(t *testing.T) {
	cg := CanGather{}
	// Assert that NeedsResource returns true using direct package call
	assert.True(t, cg.NeedsResource(), "CanGather.NeedsResource() should return true")
}

// TestCanGather_CreateAction verifies the CreateAction method returns the correct Gather action.
func TestCanGather_CreateAction(t *testing.T) {
	// Mock Location and Resource for testing
	mockLocation := &core.Location{}
	mockResource := &core.Resource{Name: "Stone"}
	mockEmptyResource := &core.Resource{}

	// Define test cases
	testCases := map[string]struct {
		canGather  CanGather
		params     core.CreateActionParams
		wantAction *actions.Gather
	}{
		"Standard gather": {
			canGather: CanGather{
				Amount:   15,
				Cost:     1.0,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Gather{
				Res:            mockResource,
				Amount:         15,
				ActionLocation: mockLocation,
				ActionCost:     1.0,
			},
		},
		"Zero amount gather": {
			canGather: CanGather{
				Amount:   0,
				Cost:     0.2,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Gather{
				Res:            mockResource,
				Amount:         0,
				ActionLocation: mockLocation,
				ActionCost:     0.2,
			},
		},
		"Zero cost gather": {
			canGather: CanGather{
				Amount:   8,
				Cost:     0.0,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Gather{
				Res:            mockResource,
				Amount:         8,
				ActionLocation: mockLocation,
				ActionCost:     0.0,
			},
		},
		"Empty resource gather": {
			canGather: CanGather{
				Amount:   3,
				Cost:     0.5,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: mockEmptyResource,
			},
			wantAction: &actions.Gather{
				Res:            mockEmptyResource,
				Amount:         3,
				ActionLocation: mockLocation,
				ActionCost:     0.5,
			},
		},
		"Nil location in CanGather": {
			canGather: CanGather{
				Amount:   5,
				Cost:     0.7,
				Location: nil,
			},
			params: core.CreateActionParams{
				Resource: mockResource,
			},
			wantAction: &actions.Gather{
				Res:            mockResource,
				Amount:         5,
				ActionLocation: nil,
				ActionCost:     0.7,
			},
		},
		"Nil resource in Params": {
			canGather: CanGather{
				Amount:   6,
				Cost:     0.9,
				Location: mockLocation,
			},
			params: core.CreateActionParams{
				Resource: nil,
			},
			wantAction: &actions.Gather{
				Res:            nil,
				Amount:         6,
				ActionLocation: mockLocation,
				ActionCost:     0.9,
			},
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotAction := tc.canGather.CreateAction(tc.params)

			// Assert that the returned action is of the expected type (*actions.Gather)
			assert.IsType(t, tc.wantAction, gotAction, "Returned action should be of type *actions.Gather")

			// Assert that the returned action is deeply equal to the expected action
			assert.Equal(t, tc.wantAction, gotAction, "Returned action does not match expected action")
		})
	}
}
