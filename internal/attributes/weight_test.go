package attributes

import (
	"errors"
	"testing"

	"Neolithic/internal/core" // Make sure this import path is correct for your project
	"github.com/stretchr/testify/assert"
)

// TestWeight_NeedsLocation verifies the NeedsLocation method returns true.
func TestWeight_NeedsLocation(t *testing.T) {
	w := Weight{}
	assert.True(t, w.NeedsLocation(), "Weight.NeedsLocation() should return true")
}

// TestWeight_NeedsResource verifies the NeedsResource method returns false.
func TestWeight_NeedsResource(t *testing.T) {
	w := Weight{}
	assert.False(t, w.NeedsResource(), "Weight.NeedsResource() should return false")
}

// TestWeight_Type verifies the Type method returns the correct attribute type.
func TestWeight_Type(t *testing.T) {
	w := Weight{}
	assert.Equal(t, WeightAttributeType, w.Type(), "Weight.Type() should return WeightAttributeType")
}

// TestWeight_Copy verifies the Copy method returns a distinct instance with the same value.
func TestWeight_Copy(t *testing.T) {
	w1 := &Weight{Amount: 75.25}
	w2Attribute := w1.Copy()

	// Assert that the type is correct
	w2, ok := w2Attribute.(*Weight)
	assert.True(t, ok, "Copy() should return a *Weight type")

	// Assert that it's a new instance, not the same pointer
	assert.NotSame(t, w1, w2, "Copy() should return a new instance")

	// Assert that the value is the same
	assert.Equal(t, w1.Amount, w2.Amount, "Copied instance should have the same Amount")

	// Modify the copy and check if the original is unchanged
	w2.Amount = 150.0
	assert.NotEqual(t, w1.Amount, w2.Amount, "Modifying the copy should not affect the original")
	assert.Equal(t, 75.25, w1.Amount, "Original amount should remain unchanged")
}

// TestWeight_String verifies the String method returns the expected format.
func TestWeight_String(t *testing.T) {
	w := Weight{Amount: 98.76}
	expected := "Weight: 98.76"
	assert.Equal(t, expected, w.String(), "Weight.String() should return the correct format for floats")

	w = Weight{Amount: 200}
	expected = "Weight: 200" // strconv.FormatFloat with 'f' and -1 precision will output "200" for whole numbers
	assert.Equal(t, expected, w.String(), "Weight.String() should handle whole numbers correctly")
}

// TestWeight_CreateAction verifies the CreateAction method handles various scenarios correctly.
func TestWeight_CreateAction(t *testing.T) {
	// Setup mock objects
	mockResource := core.NewResource("Iron Ore") // Assuming core.NewResource exists and returns *core.Resource
	mockLocation := &core.Location{}             // Assuming core.Location is a struct

	// Define test cases
	testCases := map[string]struct {
		weight     Weight
		holder     core.AttributeHolder
		params     core.CreateActionParams
		wantAction core.Action
		wantErr    error
	}{
		"Successful Gather": {
			weight: Weight{Amount: 15.0},
			holder: mockResource,
			params: core.CreateActionParams{Location: mockLocation},
			wantAction: &Gather{ // Expecting a Gather action as per the Weight.CreateAction implementation
				Res:            mockResource,
				Amount:         defaultGatherAmount, // From const in attributes package
				ActionLocation: mockLocation,
				ActionCost:     15.0, // Cost matches weight amount
			},
			wantErr: nil,
		},
		"Holder Not Resource": {
			weight:     Weight{Amount: 10.0},
			holder:     &core.Location{}, // Pass a holder that is not a *core.Resource
			params:     core.CreateActionParams{Location: mockLocation},
			wantAction: nil,
			wantErr:    errors.New("weight can only be applied to a resource"),
		},
		"Nil Location in Params": {
			weight:     Weight{Amount: 20.0},
			holder:     mockResource,
			params:     core.CreateActionParams{Location: nil}, // Location is nil
			wantAction: nil,
			wantErr:    errors.New("CreateAction was called for a resource with a weight but no location"),
		},
		"Zero Weight Amount": { // Test case for when weight amount is zero
			weight: Weight{Amount: 0.0},
			holder: mockResource,
			params: core.CreateActionParams{Location: mockLocation},
			wantAction: &Gather{
				Res:            mockResource,
				Amount:         defaultGatherAmount,
				ActionLocation: mockLocation,
				ActionCost:     0.0, // Cost is zero
			},
			wantErr: nil,
		},
	}

	// Run test cases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// It's the Weight attribute's CreateAction method we are testing
			gotAction, gotErr := tc.weight.CreateAction(tc.holder, tc.params)

			// Assert that the returned action is deeply equal to the expected action
			assert.Equal(t, tc.wantAction, gotAction, "Returned action does not match expected action")

			// Assert that the returned error matches the expected error
			// Using assert.EqualError for more precise error message checking if needed,
			// but assert.Equal works for comparing error objects (including nil).
			if tc.wantErr != nil {
				assert.EqualError(t, gotErr, tc.wantErr.Error(), "Error message does not match")
			} else {
				assert.NoError(t, gotErr, "Expected no error, but got one")
			}
		})
	}
}
