package core

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewResource(t *testing.T) {
	attr1 := &mockAttribute{attrType: "Material", value: "Wood"}
	attr2 := &mockAttribute{attrType: "Quality", value: "High"}

	type testCase struct {
		name          string
		resName       string
		opts          []ResourceOption
		expectedName  string
		expectedAttrs AttributeList
	}

	tests := []testCase{
		{
			name:          "basic creation no options",
			resName:       "IronOre",
			opts:          []ResourceOption{},
			expectedName:  "IronOre",
			expectedAttrs: NewAttributeList(),
		},
		{
			name:    "creation with one attribute",
			resName: "MagicScroll",
			opts: []ResourceOption{
				WithResourceAttributes(attr1),
			},
			expectedName: "MagicScroll",
			expectedAttrs: func() AttributeList {
				al := NewAttributeList()
				al.UpsertAttribute(attr1)
				return al
			}(),
		},
		{
			name:    "creation with multiple attributes",
			resName: "EnchantedSword",
			opts: []ResourceOption{
				WithResourceAttributes(attr1, attr2),
			},
			expectedName: "EnchantedSword",
			expectedAttrs: func() AttributeList {
				al := NewAttributeList()
				al.UpsertAttribute(attr1)
				al.UpsertAttribute(attr2)
				return al
			}(),
		},
		{
			name:    "creation with WithResourceAttributes called multiple times",
			resName: "Potion",
			opts: []ResourceOption{
				WithResourceAttributes(attr1),
				WithResourceAttributes(attr2),
			},
			expectedName: "Potion",
			expectedAttrs: func() AttributeList {
				al := NewAttributeList()
				al.UpsertAttribute(attr1)
				al.UpsertAttribute(attr2)
				return al
			}(),
		},
		{
			name:    "creation with empty WithResourceAttributes",
			resName: "Rock",
			opts: []ResourceOption{
				WithResourceAttributes(),
			},
			expectedName:  "Rock",
			expectedAttrs: NewAttributeList(),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			resource := NewResource(tc.resName, tc.opts...)

			assert.Equal(t, tc.expectedName, resource.Name)
			assert.NotNil(t, resource.attributes, "Resource attributes should be initialized by NewResource")
			assert.True(t, reflect.DeepEqual(tc.expectedAttrs, resource.attributes),
				fmt.Sprintf("Attributes mismatch.\nExpected: %v (%s)\nGot: %v (%s)",
					tc.expectedAttrs, tc.expectedAttrs.String(), resource.attributes, resource.attributes.String()))
		})
	}
}

func TestResource_String(t *testing.T) {
	attr1 := &mockAttribute{attrType: "Flammable", value: "Yes"}
	attrListWithAttr1 := NewAttributeList()
	attrListWithAttr1.UpsertAttribute(attr1)
	attrsStrWithAttr1 := attrListWithAttr1.String()

	attr2 := &mockAttribute{attrType: "Heavy", value: "True"}
	attrListWithAttr1AndAttr2 := NewAttributeList()
	attrListWithAttr1AndAttr2.UpsertAttribute(attr1)
	attrListWithAttr1AndAttr2.UpsertAttribute(attr2)
	attrsStrWithAttr1AndAttr2 := attrListWithAttr1AndAttr2.String()

	emptyAttrsStr := NewAttributeList().String()

	type testCase struct {
		name        string
		resource    *Resource
		expectedStr string
	}

	tests := []testCase{
		{
			name:        "resource with no attributes",
			resource:    NewResource("Stone"),
			expectedStr: fmt.Sprintf("Resource: Stone\nAttributes: %s", emptyAttrsStr),
		},
		{
			name: "resource with one attribute",
			resource: NewResource("Coal",
				WithResourceAttributes(attr1),
			),
			expectedStr: fmt.Sprintf("Resource: Coal\nAttributes: %s", attrsStrWithAttr1),
		},
		{
			name: "resource with multiple attributes",
			resource: NewResource("GoldBar",
				WithResourceAttributes(attr1, attr2),
			),
			expectedStr: fmt.Sprintf("Resource: GoldBar\nAttributes: %s", attrsStrWithAttr1AndAttr2),
		},
		{
			name: "resource with manually set nil attributes (to test String's nil check)",
			resource: &Resource{
				Name:       "GhostRock",
				attributes: nil, // Manually set to nil
			},
			expectedStr: "Resource: GhostRock\nAttributes: {}",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expectedStr, tc.resource.String())
		})
	}
}

func TestResource_Attributes(t *testing.T) {
	attr1 := &mockAttribute{attrType: "Edible", value: "No"}
	attrListWithA1 := NewAttributeList()
	attrListWithA1.UpsertAttribute(attr1)

	emptyAttrList := NewAttributeList()

	type testCase struct {
		name              string
		resource          *Resource
		expectedAttrsList AttributeList
	}

	tests := []testCase{
		{
			name:              "resource with attributes",
			resource:          NewResource("Apple", WithResourceAttributes(attr1)),
			expectedAttrsList: attrListWithA1,
		},
		{
			name:              "resource with no attributes (created via NewResource)",
			resource:          NewResource("Water"),
			expectedAttrsList: emptyAttrList,
		},
		{
			name: "resource with attributes, checking instance returned",
			// Attributes() returns the internal list, not a copy.
			resource: func() *Resource {
				r := NewResource("TestRes")
				r.attributes = attrListWithA1
				return r
			}(),
			expectedAttrsList: attrListWithA1,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			actualAttrs := tc.resource.Attributes()
			assert.True(t, reflect.DeepEqual(tc.expectedAttrsList, actualAttrs), "Expected %v, got %v", tc.expectedAttrsList, actualAttrs)
			if tc.resource.attributes == tc.expectedAttrsList { // Only if we expect the exact same instance
				assert.Same(t, tc.expectedAttrsList, actualAttrs, "Attributes() should return the internal instance if it matches expected")
			}
		})
	}
}
