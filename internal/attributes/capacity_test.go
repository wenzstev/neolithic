package attributes

import (
	"errors"
	"testing"

	"Neolithic/internal/core"
	"github.com/stretchr/testify/assert"
)

// TestCapacity_NeedsLocation verifies the NeedsLocation method returns false.
func TestCapacity_NeedsLocation(t *testing.T) {
	c := Capacity{}
	assert.False(t, c.NeedsLocation(), "Capacity.NeedsLocation() should return false")
}

// TestCapacity_NeedsResource verifies the NeedsResource method returns true.
func TestCapacity_NeedsResource(t *testing.T) {
	c := Capacity{}
	assert.True(t, c.NeedsResource(), "Capacity.NeedsResource() should return true")
}

// TestCapacity_Type verifies the Type method returns the correct attribute type.
func TestCapacity_Type(t *testing.T) {
	c := Capacity{}
	assert.Equal(t, CapacityAttributeType, c.Type(), "Capacity.Type() should return CapacityAttributeType")
}

// TestCapacity_Copy verifies the Copy method returns a distinct instance with the same value.
func TestCapacity_Copy(t *testing.T) {
	c1 := &Capacity{Size: 50.5}
	c2Attribute := c1.Copy()

	// Assert that the type is correct
	c2, ok := c2Attribute.(*Capacity)
	assert.True(t, ok, "Copy() should return a *Capacity type")

	// Assert that it's a new instance, not the same pointer
	assert.NotSame(t, c1, c2, "Copy() should return a new instance")

	// Assert that the value is the same
	assert.Equal(t, c1.Size, c2.Size, "Copied instance should have the same Size")

	// Modify the copy and check if the original is unchanged
	c2.Size = 100.0
	assert.NotEqual(t, c1.Size, c2.Size, "Modifying the copy should not affect the original")
	assert.Equal(t, 50.5, c1.Size, "Original size should remain unchanged")
}

// TestCapacity_String verifies the String method returns the expected format.
func TestCapacity_String(t *testing.T) {
	c := Capacity{Size: 123.45}
	expected := "Capacity: 123.45"
	assert.Equal(t, expected, c.String(), "Capacity.String() should return the correct format for floats")

	c = Capacity{Size: 100}
	expected = "Capacity: 100"
	assert.Equal(t, expected, c.String(), "Capacity.String() should handle whole numbers correctly")
}

// TestCapacity_CreateAction verifies the CreateAction method handles various scenarios correctly.
func TestCapacity_CreateAction(t *testing.T) {
	// Setup mock objects
	mockLocation := &core.Location{} // Assuming core.Location is a struct

	resLight := core.NewResource("Wood", core.WithResourceAttributes(&Weight{Amount: 5.0}))
	resHeavy := core.NewResource("Stone", core.WithResourceAttributes(&Weight{Amount: 50.0}))
	resNoWeight := core.NewResource("Feather")

	// Define test cases
	testCases := map[string]struct {
		capacity   Capacity
		holder     core.AttributeHolder
		params     core.CreateActionParams
		wantAction core.Action
		wantErr    error
	}{
		"Successful Deposit": {
			capacity: Capacity{Size: 10.0},
			holder:   mockLocation,
			params:   core.CreateActionParams{Resource: resLight},
			wantAction: &Deposit{
				DepResource:    resLight,
				Amount:         defaultDepositAmount,
				ActionLocation: mockLocation,
				ActionCost:     5.0, // Cost matches weight in the original code
			},
			wantErr: nil,
		},
		"Holder Not Location": {
			capacity:   Capacity{Size: 10.0},
			holder:     &core.Resource{}, // Pass a holder that is not a *core.Location
			params:     core.CreateActionParams{Resource: resLight},
			wantAction: nil,
			wantErr:    errors.New("capacity can only be applied to a location"),
		},
		"Nil Resource": {
			capacity:   Capacity{Size: 10.0},
			holder:     mockLocation,
			params:     core.CreateActionParams{Resource: nil},
			wantAction: nil,
			wantErr:    errors.New("CreateAction was called for a location with capacity, but no resource was provided"),
		},
		"Resource No Weight": {
			capacity:   Capacity{Size: 10.0},
			holder:     mockLocation,
			params:     core.CreateActionParams{Resource: resNoWeight},
			wantAction: nil,
			wantErr:    nil,
		},
		"Resource Too Heavy": {
			capacity:   Capacity{Size: 40.0}, // Capacity < resHeavy.Weight (50.0)
			holder:     mockLocation,
			params:     core.CreateActionParams{Resource: resHeavy},
			wantAction: nil,
			wantErr:    nil,
		},
		"Resource Exact Weight": {
			capacity: Capacity{Size: 50.0}, // Capacity == resHeavy.Weight (50.0)
			holder:   mockLocation,
			params:   core.CreateActionParams{Resource: resHeavy},
			wantAction: &Deposit{
				DepResource:    resHeavy,
				Amount:         defaultDepositAmount,
				ActionLocation: mockLocation,
				ActionCost:     50.0,
			},
			wantErr: nil,
		},
	}

	// Run test cases
	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			gotAction, gotErr := tc.capacity.CreateAction(tc.holder, tc.params)

			// Assert that the returned action is deeply equal to the expected action
			assert.Equal(t, tc.wantAction, gotAction, "Returned action does not match expected action")

			// Assert that the returned error matches the expected error
			assert.Equal(t, tc.wantErr, gotErr, "Returned error does not match expected error")
		})
	}
}
