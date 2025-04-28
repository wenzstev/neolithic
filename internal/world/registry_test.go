package world

import (
	"Neolithic/internal/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_RegisterResource(t *testing.T) {
	type testCase struct {
		resource          *core.Resource
		registry          *Registry
		expectedResources []*core.Resource
		expectedError     error
		expectedActions   []core.Action
	}

	tests := map[string]testCase{
		"can register resource, no actions": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
				},
			},
			expectedActions: []core.Action{},
		},
		"can register resource, attribute with no needs": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNoNeeds,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{ // extra stuff for test, can be ignored
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNoNeeds,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{}),
			},
		},
		"can register resource, attribute with location": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNeedsLoc,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNeedsLoc,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Location: &core.Location{Name: "loc1"},
					},
				),
			},
		},
		"can register resource, attribute with resource": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNeedsRes,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{
					{Name: "res1"},
				},
			},
			expectedResources: []*core.Resource{
				{
					Name: "res1",
				},
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNeedsRes,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
					},
				),
			},
		},
		"can register resource, attribute with both": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNeedsBoth,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{
					{Name: "res1"},
				},
			},
			expectedResources: []*core.Resource{
				{
					Name: "res1",
				},
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNeedsBoth,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
						Location: &core.Location{Name: "loc1"},
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
						Location: &core.Location{Name: "loc2"},
					}),
			},
		},
		"can register resource, resource has attribute": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNeedsBoth,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{
					{
						Name: "res1",
						Attributes: []core.Attribute{
							mockAttributeNeedsRes,
						},
					},
				},
			},
			expectedResources: []*core.Resource{
				{
					Name: "res1",
					Attributes: []core.Attribute{
						mockAttributeNeedsRes,
					},
				},
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNeedsBoth,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{
							Name: "res1",
							Attributes: []core.Attribute{
								mockAttributeNeedsRes,
							},
						},
						Location: &core.Location{Name: "loc1"},
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{
							Name: "res1",
							Attributes: []core.Attribute{
								mockAttributeNeedsRes,
							},
						},
						Location: &core.Location{Name: "loc2"},
					}),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{
							Name: "testResource",
							Attributes: []core.Attribute{
								mockAttributeNeedsBoth,
							},
						},
					}),
			},
		},
		"can register resource, two attributes": {
			resource: &core.Resource{
				Name: "testResource",
				Attributes: []core.Attribute{
					mockAttributeNeedsBoth,
					mockAttributeNeedsRes,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{
					{Name: "res1"},
				},
			},
			expectedResources: []*core.Resource{
				{
					Name: "res1",
				},
				{
					Name: "testResource",
					Attributes: []core.Attribute{
						mockAttributeNeedsBoth,
						mockAttributeNeedsRes,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
						Location: &core.Location{Name: "loc1"},
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
						Location: &core.Location{Name: "loc2"},
					}),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
					}),
			},
		},
		"location attribute needs resource": { // The original selected test for reference
			resource: &core.Resource{Name: "newResource"}, // Resource being registered
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1", Attributes: []core.Attribute{mockLocationAttributeNeedsRes}}, // Location with attribute needing resource
					loc2, // Location without relevant attribute
				},
				Resources: []*core.Resource{res1}, // Existing resource doesn't affect this path
			},
			expectedResources: []*core.Resource{
				res1,                  // Existing resource
				{Name: "newResource"}, // New resource added
			},
			// Action created by step 3 (default case) for loc1's attribute needing the new resource
			expectedActions: []core.Action{
				mockLocationAttributeNeedsRes.CreateAction(core.CreateActionParams{Resource: &core.Resource{Name: "newResource"}}),
			},
			expectedError: nil,
		},
		"fail to register resource": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{
					{Name: "testResource"},
				},
			},
			expectedResources: []*core.Resource{
				{Name: "testResource"},
			},
			expectedError:   ErrResourceAlreadyRegistered,
			expectedActions: []core.Action{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.registry.RegisterResource(tc.resource)
			if tc.expectedError != nil {
				require.ErrorIs(t, tc.expectedError, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedResources, tc.registry.Resources)
			require.Equal(t, tc.expectedActions, tc.registry.Actions)
		})
	}
}

func TestRegistry_RegisterLocation(t *testing.T) {
	type testCase struct {
		location          *core.Location
		registry          *Registry
		expectedLocations []*core.Location
		expectedError     error
		expectedActions   []core.Action
	}

	tests := map[string]testCase{
		"can register location, no actions": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				{Name: "testLocation"},
			},
			expectedActions: []core.Action{},
		},
		"can register location, attribute with location": {
			location: &core.Location{
				Name: "testLocation",
				Attributes: []core.Attribute{
					mockAttributeNeedsLoc,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
				},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				{Name: "loc1"},
				{
					Name: "testLocation",
					Attributes: []core.Attribute{
						mockAttributeNeedsLoc,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{Location: &core.Location{Name: "loc1"}}),
			},
		},
		"can register location, attribute needs resource": {
			location: &core.Location{
				Name: "testLocation",
				Attributes: []core.Attribute{
					mockAttributeNeedsRes,
				},
			},
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{
					{Name: "res1"},
				},
			},
			expectedLocations: []*core.Location{
				{
					Name: "testLocation",
					Attributes: []core.Attribute{
						mockAttributeNeedsRes,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: &core.Resource{Name: "res1"},
					},
				),
			},
		},
		"can register location, attribute needs both": {
			location: &core.Location{
				Name: "testLocation",
				Attributes: []core.Attribute{
					mockAttributeNeedsBoth,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
				},
				Resources: []*core.Resource{
					{Name: "res1"},
					{Name: "res2"},
				},
			},
			expectedLocations: []*core.Location{
				{Name: "loc1"},
				{
					Name: "testLocation",
					Attributes: []core.Attribute{
						mockAttributeNeedsBoth,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: &core.Location{Name: "loc1"},
					Resource: &core.Resource{Name: "res1"},
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: &core.Location{Name: "loc1"},
					Resource: &core.Resource{Name: "res2"},
				}),
			},
		},
		"can register location, two attributes": {
			location: &core.Location{
				Name: "testLocation",
				Attributes: []core.Attribute{
					mockAttributeNeedsBoth,
					mockAttributeNeedsRes,
				},
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
				},
				Resources: []*core.Resource{
					{Name: "res1"},
					{Name: "res2"},
				},
			},
			expectedLocations: []*core.Location{
				{Name: "loc1"},
				{
					Name: "testLocation",
					Attributes: []core.Attribute{
						mockAttributeNeedsBoth,
						mockAttributeNeedsRes,
					},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: &core.Location{Name: "loc1"},
					Resource: &core.Resource{Name: "res1"},
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: &core.Location{Name: "loc1"},
					Resource: &core.Resource{Name: "res2"},
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Resource: &core.Resource{Name: "res1"},
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Resource: &core.Resource{Name: "res2"},
				}),
			},
		},
		"fail to register location": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "testLocation"},
				},
			},
			expectedError:   ErrLocationAlreadyRegistered,
			expectedActions: []core.Action{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.registry.RegisterLocation(tc.location)
			if tc.expectedError != nil {
				require.ErrorIs(t, tc.expectedError, err)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedLocations, tc.registry.Locations)
			require.Equal(t, tc.expectedActions, tc.registry.Actions)
		})
	}
}
