package core

import (
	"fmt"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Helper to build an inventory for test expectations using the actual Inventory implementation.
func newTestInventory(entries ...InventoryEntry) Inventory {
	inv := NewInventory() // Assumes NewInventory() is available from the core package
	for _, entry := range entries {
		inv.AdjustAmount(entry.Resource, entry.Amount)
	}
	return inv
}

// Helper to build the expected string for an inventory based on common patterns.
// This assumes Inventory.String() produces a sorted, predictable output.
// For an empty inventory, it's assumed to be "{}".
// For non-empty, format like "  Resource1: Amount1\n  Resource2: Amount2" (sorted by name).
// This needs to match the actual Inventory.String() behavior.
func buildExpectedInventoryString(entries ...InventoryEntry) string {
	if len(entries) == 0 {
		return "{}"
	}
	// Sort entries by resource name for predictable string output
	sort.SliceStable(entries, func(i, j int) bool {
		if entries[i].Resource == nil || entries[j].Resource == nil {
			return entries[i].Resource != nil // nil resources last or first consistently
		}
		return entries[i].Resource.Name < entries[j].Resource.Name
	})

	var parts []string
	for _, entry := range entries {
		if entry.Resource != nil { // Guard against nil resource
			parts = append(parts, fmt.Sprintf("  %s: %d", entry.Resource.Name, entry.Amount))
		}
	}
	return strings.Join(parts, "\n")
}

func TestNewLocation(t *testing.T) {
	coord1 := Coord{X: 1, Y: 2}
	res1 := &Resource{Name: "Wood"}
	res2 := &Resource{Name: "Stone"}
	attr1 := &mockAttribute{attrType: "Terrain", value: "Forest"}
	attr2 := &mockAttribute{attrType: "Property", value: "Mine"}

	type testCase struct {
		name              string
		locName           string
		coord             Coord
		opts              []LocationOption
		expectedName      string
		expectedCoord     Coord
		expectedInventory Inventory
		expectedAttrs     AttributeList
	}

	tests := []testCase{
		{
			name:              "basic creation",
			locName:           "Testville",
			coord:             coord1,
			opts:              []LocationOption{},
			expectedName:      "Testville",
			expectedCoord:     coord1,
			expectedInventory: newTestInventory(), // Empty inventory
			expectedAttrs:     NewAttributeList(),
		},
		{
			name:    "with inventory",
			locName: "ResourceSpot",
			coord:   Coord{X: 10, Y: 10},
			opts: []LocationOption{
				WithInventory(
					InventoryEntry{Resource: res1, Amount: 100},
					InventoryEntry{Resource: res2, Amount: 50},
				),
			},
			expectedName:  "ResourceSpot",
			expectedCoord: Coord{X: 10, Y: 10},
			expectedInventory: newTestInventory(
				InventoryEntry{Resource: res1, Amount: 100},
				InventoryEntry{Resource: res2, Amount: 50},
			),
			expectedAttrs: NewAttributeList(),
		},
		{
			name:    "with attributes",
			locName: "SpecialPlace",
			coord:   Coord{X: 5, Y: 5},
			opts: []LocationOption{
				WithAttributes(attr1, attr2),
			},
			expectedName:      "SpecialPlace",
			expectedCoord:     Coord{X: 5, Y: 5},
			expectedInventory: newTestInventory(),
			expectedAttrs: func() AttributeList {
				al := NewAttributeList()
				al.UpsertAttribute(attr1)
				al.UpsertAttribute(attr2)
				return al
			}(),
		},
		{
			name:    "with inventory and attributes",
			locName: "FullFeature",
			coord:   Coord{X: 0, Y: 0},
			opts: []LocationOption{
				WithInventory(InventoryEntry{Resource: res1, Amount: 10}),
				WithAttributes(attr1),
			},
			expectedName:      "FullFeature",
			expectedCoord:     Coord{X: 0, Y: 0},
			expectedInventory: newTestInventory(InventoryEntry{Resource: res1, Amount: 10}),
			expectedAttrs: func() AttributeList {
				al := NewAttributeList()
				al.UpsertAttribute(attr1)
				return al
			}(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			loc := NewLocation(tc.locName, tc.coord, tc.opts...)

			assert.Equal(t, tc.expectedName, loc.Name)
			assert.Equal(t, tc.expectedCoord, loc.Coord)

			// Compare inventories using reflect.DeepEqual, assuming AdjustAmount sorts them for canonical order.
			// If DeepEqual fails due to pointer differences in Resource or complex internal state,
			// comparing String() or implementing a specific IsEqual method on Inventory would be alternatives.
			assert.True(t, reflect.DeepEqual(tc.expectedInventory, loc.Inventory),
				fmt.Sprintf("Inventory mismatch.\nExpected: %v (%s)\nGot: %v (%s)",
					tc.expectedInventory, tc.expectedInventory.String(), loc.Inventory, loc.Inventory.String()))

			assert.True(t, reflect.DeepEqual(tc.expectedAttrs, loc.Attributes()), // Use Attributes() getter
				fmt.Sprintf("Attributes mismatch.\nExpected: %v (%s)\nGot: %v (%s)",
					tc.expectedAttrs, tc.expectedAttrs.String(), loc.Attributes(), loc.Attributes().String()))
		})
	}
}

func TestLocation_String(t *testing.T) {
	// Assuming Coord{X:5, Y:10}.String() returns "(5,10)"
	// If your actual Coord.String() is different, adjust coordStr accordingly.
	mockCoord := Coord{X: 5, Y: 10}
	coordStr := mockCoord.String() // Use the actual String() output

	attr1 := &mockAttribute{attrType: "Type1", value: "Val1"}
	attrListWithAttr1 := NewAttributeList()
	attrListWithAttr1.UpsertAttribute(attr1)
	attrsStrWithAttr1 := attrListWithAttr1.String() // e.g., "(Type: Type1, Value: Val1)"
	attrsStrEmpty := NewAttributeList().String()    // e.g., "{}"

	// Define expected inventory strings based on buildExpectedInventoryString helper
	// or by knowing the exact output of your actual Inventory.String()
	invEntries1 := []InventoryEntry{{Resource: &Resource{Name: "ItemA"}, Amount: 10}}
	inventoryWithItemA := newTestInventory(invEntries1...)
	invStrItemA := inventoryWithItemA.String() // Get actual string output

	inventoryEmpty := newTestInventory()
	invStrEmpty := inventoryEmpty.String() // Should be "{}" if buildExpectedInventoryString is accurate

	type testCase struct {
		name        string
		location    *Location
		expectedStr string
	}

	tests := []testCase{
		{
			name: "basic location",
			location: NewLocation("Town", mockCoord,
				WithInventory(invEntries1...),
				WithAttributes(attr1),
			),
			expectedStr: fmt.Sprintf("Location: Town\nCoordinates: %s\nInventory: %s\nAttributes: %s", coordStr, invStrItemA, attrsStrWithAttr1),
		},
		{
			name: "empty inventory and attributes",
			location: NewLocation("EmptySpot", mockCoord,
				WithInventory(),  // No inventory entries
				WithAttributes(), // No attributes
			),
			expectedStr: fmt.Sprintf("Location: EmptySpot\nCoordinates: %s\nInventory: %s\nAttributes: %s", coordStr, invStrEmpty, attrsStrEmpty),
		},
		{
			name: "nil inventory and attributes in constructed Location (String method should handle)",
			// NewLocation initializes Inventory and attributes, so they won't be nil.
			// To test String's nil handling, we'd have to manually set them to nil post-creation.
			location: &Location{
				Name:       "Nilsville",
				Coord:      mockCoord,
				Inventory:  nil, // Manually set to nil
				attributes: nil, // Manually set to nil
			},
			expectedStr: fmt.Sprintf("Location: Nilsville\nCoordinates: %s\nInventory: {}\nAttributes: {}", coordStr),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedStr, tc.location.String())
		})
	}
}

func TestLocation_Attributes(t *testing.T) {
	attr1 := &mockAttribute{attrType: "A1", value: "V1"}
	attrListWithA1 := NewAttributeList()
	attrListWithA1.UpsertAttribute(attr1)

	type testCase struct {
		name              string
		location          *Location
		expectedAttrsList AttributeList
	}
	tests := []testCase{
		{
			name:              "location with attributes",
			location:          NewLocation("AttrLoc", Coord{}, WithAttributes(attr1)),
			expectedAttrsList: attrListWithA1,
		},
		{
			name:              "location with no attributes (NewLocation initializes to empty)",
			location:          NewLocation("NoAttrLoc", Coord{}), // No WithAttributes option
			expectedAttrsList: NewAttributeList(),
		},
		{
			name:              "location with explicitly empty attribute list",
			location:          NewLocation("EmptyAttrLoc", Coord{}, WithAttributes()), // Empty slice
			expectedAttrsList: NewAttributeList(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualAttrs := tc.location.Attributes()
			assert.True(t, reflect.DeepEqual(tc.expectedAttrsList, actualAttrs), "Expected %v, got %v", tc.expectedAttrsList, actualAttrs)
			// The current Attributes() returns the internal instance. If it were to return a copy,
			// assert.NotSame(t, tc.expectedAttrsList, actualAttrs) would be relevant for the list itself,
			// but DeepEqual checks content.
		})
	}
}

func TestLocation_DeepCopy(t *testing.T) {
	res1 := &Resource{Name: "testres"}
	attr1 := &mockAttribute{attrType: "TestAttr", value: "Val1"}

	initialInventoryEntries := []InventoryEntry{{Resource: res1, Amount: 4}}

	type testCase struct {
		name     string
		location *Location
	}

	tests := map[string]testCase{
		"basic copy": {
			name: "basic copy with all fields",
			location: NewLocation("testLocation", Coord{X: 5, Y: 2},
				WithInventory(initialInventoryEntries...),
				WithAttributes(attr1),
			),
		},
		"empty inventory and attributes": {
			name: "empty inventory and attributes",
			location: NewLocation("EmptyVille", Coord{X: 1, Y: 1},
				WithInventory(),
				WithAttributes(),
			),
		},
		"location with nil fields (DeepCopy should initialize them)": {
			name: "location with nil fields (DeepCopy should initialize them)",
			location: &Location{ // Manually create to test nil handling in DeepCopy
				Name:       "NilSource",
				Coord:      Coord{X: 0, Y: 0},
				Inventory:  nil,
				attributes: nil,
			},
		},
	}

	for testName, tc := range tests {
		t.Run(testName, func(t *testing.T) {
			originalLocation := tc.location
			copiedLocation := originalLocation.DeepCopy()

			// 1. Check for basic equality and different instances
			assert.Equal(t, originalLocation.Name, copiedLocation.Name, "Name should be equal")
			assert.Equal(t, originalLocation.Coord, copiedLocation.Coord, "Coord should be equal")
			assert.NotSame(t, originalLocation, copiedLocation, "Copied location should be a new instance")

			// 2. Check Inventory deep copy
			if originalLocation.Inventory != nil {
				assert.NotNil(t, copiedLocation.Inventory, "Copied inventory should not be nil if original was not")
				assert.NotSame(t, originalLocation.Inventory, copiedLocation.Inventory, "Inventory should be a new instance")
				assert.True(t, reflect.DeepEqual(originalLocation.Inventory, copiedLocation.Inventory), "Inventory content should match after copy")

				// Mutate original inventory, copy should not change
				originalLocation.Inventory.AdjustAmount(&Resource{Name: "NewResInOrig"}, 100)
				assert.False(t, reflect.DeepEqual(originalLocation.Inventory, copiedLocation.Inventory), "Copied inventory should not change when original inventory is modified")

				// Mutate copied inventory, original should not change (reset original for clean test)
				freshOriginalForInvTest := tc.location.DeepCopy() // Get a fresh original state
				freshCopyForInvTest := freshOriginalForInvTest.DeepCopy()
				if freshCopyForInvTest.Inventory != nil { // Should not be nil if original wasn't
					freshCopyForInvTest.Inventory.AdjustAmount(&Resource{Name: "NewResInCopy"}, 200)
					assert.False(t, reflect.DeepEqual(freshOriginalForInvTest.Inventory, freshCopyForInvTest.Inventory), "Original inventory should not change when copied inventory is modified")
				}
			} else { // Original inventory was nil
				assert.NotNil(t, copiedLocation.Inventory, "Copied inventory should be a new empty inventory if original was nil")
				assert.Equal(t, newTestInventory().String(), copiedLocation.Inventory.String(), "Copied inventory should be empty if original was nil")
			}

			// 3. Check Attributes (AttributeList) deep copy
			if originalLocation.attributes != nil {
				assert.NotNil(t, copiedLocation.attributes, "Copied attributes should not be nil if original was not")
				assert.NotSame(t, originalLocation.attributes, copiedLocation.attributes, "AttributeList should be a new instance")
				assert.True(t, reflect.DeepEqual(originalLocation.attributes, copiedLocation.attributes), "AttributeList content should match after copy")

				// Mutate original attributes, copy should not change
				originalLocation.attributes.UpsertAttribute(&mockAttribute{attrType: "OrigOnly", value: "VOrig"})
				assert.False(t, reflect.DeepEqual(originalLocation.attributes, copiedLocation.attributes), "Copied attributes should not change when original attributes are modified")

				// Mutate copied attributes, original should not change (reset original for clean test)
				freshOriginalForAttrTest := tc.location.DeepCopy() // Get a fresh original state
				freshCopyForAttrTest := freshOriginalForAttrTest.DeepCopy()
				if freshCopyForAttrTest.attributes != nil { // Should not be nil if original wasn't
					freshCopyForAttrTest.attributes.UpsertAttribute(&mockAttribute{attrType: "CopyOnly", value: "VCopy"})
					assert.False(t, reflect.DeepEqual(freshOriginalForAttrTest.attributes, freshCopyForAttrTest.attributes), "Original attributes should not change when copied attributes are modified")
				}
			} else { // Original attributes were nil
				assert.NotNil(t, copiedLocation.attributes, "Copied attributes should be a new empty list if original was nil")
				assert.Equal(t, NewAttributeList().String(), copiedLocation.attributes.String(), "Copied attributes should be empty if original was nil")
			}

			// Final deep equality check on a fresh copy (before any local mutations in this test run)
			cleanOriginal := tc.location
			cleanCopied := cleanOriginal.DeepCopy() // Perform a fresh copy for this check

			assert.Equal(t, cleanOriginal.Name, cleanCopied.Name)
			assert.Equal(t, cleanOriginal.Coord, cleanCopied.Coord)

			if cleanOriginal.Inventory != nil {
				assert.True(t, reflect.DeepEqual(cleanOriginal.Inventory, cleanCopied.Inventory), "Clean copy inventory deep equal check failed")
			} else {
				assert.NotNil(t, cleanCopied.Inventory, "Clean copy inventory should be non-nil if original was nil")
				assert.Equal(t, newTestInventory().String(), cleanCopied.Inventory.String(), "Clean copy inventory should be empty if original was nil")
			}

			if cleanOriginal.attributes != nil {
				assert.True(t, reflect.DeepEqual(cleanOriginal.attributes, cleanCopied.attributes), "Clean copy attributes deep equal check failed")
			} else {
				assert.NotNil(t, cleanCopied.attributes, "Clean copy attributes should be non-nil if original was nil")
				assert.Equal(t, NewAttributeList().String(), cleanCopied.attributes.String(), "Clean copy attributes should be empty if original was nil")
			}
		})
	}
}
