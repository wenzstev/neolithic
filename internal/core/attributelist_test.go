package core

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// mockAttribute is a simple implementation of the Attribute interface for testing.
type mockAttribute struct {
	attrType AttributeType
	value    string // Used for String() and to differentiate instances
}

func (m *mockAttribute) String() string {
	return fmt.Sprintf("(Type: %s, Value: %s)", m.attrType, m.value)
}

func (m *mockAttribute) Type() AttributeType {
	return m.attrType
}

// CreateAction is not relevant for attributeList tests, returns nil.
func (m *mockAttribute) CreateAction(holder AttributeHolder, params CreateActionParams) (Action, error) {
	return nil, nil
}

// NeedsLocation is not relevant for attributeList tests, returns false.
func (m *mockAttribute) NeedsLocation() bool {
	return false
}

// NeedsResource is not relevant for attributeList tests, returns false.
func (m *mockAttribute) NeedsResource() bool {
	return false
}

// Copy returns a copy of the mockAttribute.
func (m *mockAttribute) Copy() Attribute {
	return &mockAttribute{attrType: m.attrType, value: m.value}
}

// Helper to create an AttributeList populated with given attributes for expected states.
// Attributes will be upserted, so order of input doesn't matter for final sorted list.
func newTestAttributeList(attrs ...Attribute) AttributeList {
	al := NewAttributeList()
	for _, attr := range attrs {
		al.UpsertAttribute(attr)
	}
	return al
}

func TestAttributeList_UpsertAttribute(t *testing.T) {
	type testCase struct {
		name         string
		initialList  AttributeList
		attrToUpsert Attribute
		expectedList AttributeList
	}

	tests := []testCase{
		{
			name:         "add to empty list",
			initialList:  NewAttributeList(),
			attrToUpsert: &mockAttribute{attrType: "typeA", value: "valA"},
			expectedList: newTestAttributeList(&mockAttribute{attrType: "typeA", value: "valA"}),
		},
		{
			name:         "add new attribute (maintaining order - beginning)",
			initialList:  newTestAttributeList(&mockAttribute{attrType: "typeB", value: "valB"}),
			attrToUpsert: &mockAttribute{attrType: "typeA", value: "valA"},
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeB", value: "valB"}, // Original
				&mockAttribute{attrType: "typeA", value: "valA"}, // Added
			),
		},
		{
			name: "add new attribute (maintaining order - middle)",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
			attrToUpsert: &mockAttribute{attrType: "typeB", value: "valB"},
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeC", value: "valC"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
		},
		{
			name:         "add new attribute (maintaining order - end)",
			initialList:  newTestAttributeList(&mockAttribute{attrType: "typeA", value: "valA"}),
			attrToUpsert: &mockAttribute{attrType: "typeB", value: "valB"},
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
		},
		{
			name: "update existing attribute",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA_old"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
			attrToUpsert: &mockAttribute{attrType: "typeA", value: "valA_new"},
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA_new"}, // Updated
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
		},
		{
			name:         "add multiple attributes",
			initialList:  NewAttributeList(),
			attrToUpsert: &mockAttribute{attrType: "typeC", value: "valC"}, // This will be the first upsert
			// expectedList will be constructed step-by-step in the test logic
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			// Make a copy of initialList to avoid modifying the slice in the test case definition directly
			// as attributeList methods operate on pointers.
			currentList := tc.initialList.Copy()
			currentList.UpsertAttribute(tc.attrToUpsert)

			if tc.name == "add multiple attributes" {
				// Special handling for this case to demonstrate multiple upserts
				currentList.UpsertAttribute(&mockAttribute{attrType: "typeA", value: "valA"})
				currentList.UpsertAttribute(&mockAttribute{attrType: "typeB", value: "valB"})
				expected := newTestAttributeList(
					&mockAttribute{attrType: "typeC", value: "valC"},
					&mockAttribute{attrType: "typeA", value: "valA"},
					&mockAttribute{attrType: "typeB", value: "valB"},
				)
				assert.True(t, reflect.DeepEqual(currentList, expected), "Expected list to be %v, but got %v", expected, currentList)
			} else {
				assert.True(t, reflect.DeepEqual(currentList, tc.expectedList), "Expected list to be %v, but got %v", tc.expectedList, currentList)
			}
		})
	}
}

func TestAttributeList_RemoveAttribute(t *testing.T) {
	type testCase struct {
		name         string
		initialList  AttributeList
		attrToRemove AttributeType
		expectedList AttributeList
	}

	tests := []testCase{
		{
			name:         "remove from single-item list",
			initialList:  newTestAttributeList(&mockAttribute{attrType: "typeA", value: "valA"}),
			attrToRemove: "typeA",
			expectedList: NewAttributeList(),
		},
		{
			name: "remove from beginning of list",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
			attrToRemove: "typeA",
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeB", value: "valB"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
		},
		{
			name: "remove from middle of list",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
			attrToRemove: "typeB",
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
		},
		{
			name: "remove from end of list",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
				&mockAttribute{attrType: "typeC", value: "valC"},
			),
			attrToRemove: "typeC",
			expectedList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
		},
		{
			name: "remove non-existent attribute",
			initialList: newTestAttributeList(
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
			attrToRemove: "typeX", // Non-existent
			expectedList: newTestAttributeList( // Should remain unchanged
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
		},
		{
			name:         "remove from empty list",
			initialList:  NewAttributeList(),
			attrToRemove: "typeA",
			expectedList: NewAttributeList(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			currentList := tc.initialList.Copy() // Work on a copy
			currentList.RemoveAttribute(tc.attrToRemove)
			assert.True(t, reflect.DeepEqual(currentList, tc.expectedList), "Expected list to be %v, but got %v", tc.expectedList, currentList)
		})
	}
}

func TestAttributeList_AttributeByType(t *testing.T) {
	type testCase struct {
		name          string
		initialList   AttributeList
		attrTypeToGet AttributeType
		expectedAttr  Attribute
		expectFound   bool
	}

	attrA := &mockAttribute{attrType: "typeA", value: "valA"}
	attrB := &mockAttribute{attrType: "typeB", value: "valB"}

	tests := []testCase{
		{
			name:          "get existing attribute",
			initialList:   newTestAttributeList(attrA, attrB),
			attrTypeToGet: "typeA",
			expectedAttr:  attrA,
			expectFound:   true,
		},
		{
			name:          "get another existing attribute",
			initialList:   newTestAttributeList(attrA, attrB),
			attrTypeToGet: "typeB",
			expectedAttr:  attrB,
			expectFound:   true,
		},
		{
			name:          "get non-existent attribute from non-empty list",
			initialList:   newTestAttributeList(attrA, attrB),
			attrTypeToGet: "typeX",
			expectedAttr:  nil,
			expectFound:   false,
		},
		{
			name:          "get from empty list",
			initialList:   NewAttributeList(),
			attrTypeToGet: "typeA",
			expectedAttr:  nil,
			expectFound:   false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			foundAttr := tc.initialList.AttributeByType(tc.attrTypeToGet)
			if tc.expectFound {
				assert.NotNil(t, foundAttr)
				assert.True(t, reflect.DeepEqual(foundAttr, tc.expectedAttr), "Expected attribute %v, but got %v", tc.expectedAttr, foundAttr)
			} else {
				assert.Nil(t, foundAttr)
			}
		})
	}
}

func TestAttributeList_String(t *testing.T) {
	type testCase struct {
		name        string
		initialList AttributeList
		expectedStr string
	}

	tests := []testCase{
		{
			name:        "empty list",
			initialList: NewAttributeList(),
			expectedStr: "{}",
		},
		{
			name:        "single attribute",
			initialList: newTestAttributeList(&mockAttribute{attrType: "typeA", value: "valA"}),
			expectedStr: "(Type: typeA, Value: valA)",
		},
		{
			name: "multiple attributes (sorted by type)",
			initialList: newTestAttributeList( // Inserted in non-sorted order
				&mockAttribute{attrType: "typeC", value: "valC"},
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeB", value: "valB"},
			),
			// Expected string should have attributes sorted by type
			expectedStr: "(Type: typeA, Value: valA)(Type: typeB, Value: valB)(Type: typeC, Value: valC)",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedStr, tc.initialList.String())
		})
	}
}

func TestAttributeList_Copy(t *testing.T) {
	type testCase struct {
		name    string
		initial AttributeList
	}

	tests := []testCase{
		{
			name:    "empty list",
			initial: NewAttributeList(),
		},
		{
			name:    "list with one attribute",
			initial: newTestAttributeList(&mockAttribute{attrType: "typeA", value: "valA"}),
		},
		{
			name: "list with multiple attributes",
			initial: newTestAttributeList(
				&mockAttribute{attrType: "typeB", value: "valB"}, // Intentionally not in order of type
				&mockAttribute{attrType: "typeA", value: "valA"},
				&mockAttribute{attrType: "typeC", value: "valC_original"},
			),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			originalList := tc.initial.Copy() // Work with a copy of the test case's initial state

			// 1. Basic copy properties
			copiedList := originalList.Copy()
			assert.True(t, reflect.DeepEqual(originalList, copiedList), "Copied list should be deep equal to its source")
			assert.NotSame(t, originalList, copiedList, "Copied list should not be the same instance as its source")

			// If not empty, check that attributes are also new instances
			originalConcrete := originalList.(*attributeList)
			copiedConcrete := copiedList.(*attributeList)
			if len(*originalConcrete) > 0 {
				assert.NotSame(t, (*originalConcrete)[0], (*copiedConcrete)[0], "First attribute in copied list should be a new instance")
			}

			// 2. Mutate structure of originalList, copiedList should be unaffected
			originalList.UpsertAttribute(&mockAttribute{attrType: "typeZ_added", value: "valZ"})
			assert.False(t, reflect.DeepEqual(originalList, copiedList), "After originalList structure mutated, it should differ from copiedList")

			// 3. Mutate attribute *within* originalList, copiedList should be unaffected
			//    Reset originalList for this sub-test to avoid interference from previous mutation.
			originalListForAttrMutation := tc.initial.Copy() // Fresh copy of initial state
			copiedListForAttrTest := originalListForAttrMutation.Copy()

			if oConcrete, ok := originalListForAttrMutation.(*attributeList); ok && len(*oConcrete) > 0 {
				// Mutate the first attribute if it's a mockAttribute
				if ma, okMA := (*oConcrete)[0].(*mockAttribute); okMA {
					originalValue := ma.value
					ma.value = "MUTATED_IN_ORIGINAL"

					assert.False(t, reflect.DeepEqual(originalListForAttrMutation, copiedListForAttrTest), "After attribute in originalListForAttrMutation mutated, it should differ from copiedListForAttrTest")

					// Verify the attribute in copiedListForAttrTest was not affected
					if cConcrete, okC := copiedListForAttrTest.(*attributeList); okC && len(*cConcrete) > 0 {
						if maCopied, okMAC := (*cConcrete)[0].(*mockAttribute); okMAC {
							// Ensure it's the same attribute type we are comparing
							if maCopied.Type() == ma.Type() {
								assert.Equal(t, originalValue, maCopied.value, "Attribute in copiedListForAttrTest should retain its original copied value")
							}
						}
					}
				}
			}

			// 4. Mutate structure of copiedList, originalList should be unaffected
			//    Reset originalList again for clarity.
			originalListForCopyMutation := tc.initial.Copy()
			copiedListToMutate := originalListForCopyMutation.Copy()

			copiedListToMutate.UpsertAttribute(&mockAttribute{attrType: "typeY_added_to_copy", value: "valY"})
			assert.False(t, reflect.DeepEqual(originalListForCopyMutation, copiedListToMutate), "After copiedListToMutate structure mutated, it should differ from originalListForCopyMutation")

			// 5. Mutate attribute *within* copiedList, originalList should be unaffected
			originalListForCopyAttrMutation := tc.initial.Copy()
			copiedListToMutateAttr := originalListForCopyAttrMutation.Copy()

			if cConcrete, ok := copiedListToMutateAttr.(*attributeList); ok && len(*cConcrete) > 0 {
				if ma, okMA := (*cConcrete)[0].(*mockAttribute); okMA {
					originalValueInOriginal := ""
					// Find corresponding attribute in original to get its pre-mutation value
					if oConcreteOrig, okO := originalListForCopyAttrMutation.(*attributeList); okO && len(*oConcreteOrig) > 0 {
						if maOrig, okMAO := (*oConcreteOrig)[0].(*mockAttribute); okMAO && maOrig.Type() == ma.Type() {
							originalValueInOriginal = maOrig.value
						}
					}

					ma.value = "MUTATED_IN_COPY"
					assert.False(t, reflect.DeepEqual(originalListForCopyAttrMutation, copiedListToMutateAttr), "After attribute in copiedListToMutateAttr mutated, it should differ from originalListForCopyAttrMutation")

					// Verify the attribute in originalListForCopyAttrMutation was not affected
					if oConcreteOrig, okO := originalListForCopyAttrMutation.(*attributeList); okO && len(*oConcreteOrig) > 0 {
						if maOrig, okMAO := (*oConcreteOrig)[0].(*mockAttribute); okMAO && maOrig.Type() == ma.Type() {
							assert.Equal(t, originalValueInOriginal, maOrig.value, "Attribute in originalListForCopyAttrMutation should retain its original value")
						}
					}
				}
			}
		})
	}
}
