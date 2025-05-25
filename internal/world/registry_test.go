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
			resource: core.NewResource("testResource"),
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				core.NewResource("testResource"),
			},
			expectedActions: []core.Action{},
		},
		"can register resource, attribute with no needs": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNoNeeds)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
					core.NewLocation("loc2", core.Coord{}),
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNoNeeds)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{}),
			},
		},
		"can register resource, attribute with location": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsLoc)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsLoc)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Location: core.NewLocation("loc1", core.Coord{}),
					},
				),
			},
		},
		"can register resource, attribute with resource": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsRes)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
					core.NewLocation("loc2", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1"),
				},
			},
			expectedResources: []*core.Resource{
				core.NewResource("res1"),
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsRes)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
					},
				),
			},
		},
		"can register resource, attribute with both": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
					core.NewLocation("loc2", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1"),
				},
			},
			expectedResources: []*core.Resource{
				core.NewResource("res1"),
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
						Location: core.NewLocation("loc1", core.Coord{}),
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
						Location: core.NewLocation("loc2", core.Coord{}),
					}),
			},
		},
		"can register resource, resource has attribute": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
					core.NewLocation("loc2", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1", core.WithResourceAttributes(mockAttributeNeedsRes)),
				},
			},
			expectedResources: []*core.Resource{
				core.NewResource("res1", core.WithResourceAttributes(mockAttributeNeedsRes)),
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1", core.WithResourceAttributes(mockAttributeNeedsRes)),
						Location: core.NewLocation("loc1", core.Coord{}),
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1", core.WithResourceAttributes(mockAttributeNeedsRes)),
						Location: core.NewLocation("loc2", core.Coord{}),
					}),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth)),
					}),
			},
		},
		"can register resource, two attributes": {
			resource: core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth, mockAttributeNeedsRes)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
					core.NewLocation("loc2", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1"),
				},
			},
			expectedResources: []*core.Resource{
				core.NewResource("res1"),
				core.NewResource("testResource", core.WithResourceAttributes(mockAttributeNeedsBoth, mockAttributeNeedsRes)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
						Location: core.NewLocation("loc1", core.Coord{}),
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
						Location: core.NewLocation("loc2", core.Coord{}),
					}),
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
					}),
			},
		},
		"fail to register resource": {
			resource: core.NewResource("testResource"),
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{
					core.NewResource("testResource"),
				},
			},
			expectedResources: []*core.Resource{
				core.NewResource("testResource"),
			},
			expectedError:   ErrResourceAlreadyRegistered,
			expectedActions: []core.Action{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			err := tc.registry.RegisterResource(tc.resource)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)

			require.ElementsMatch(t, tc.expectedResources, tc.registry.Resources, "Resources do not match")
			require.ElementsMatch(t, tc.expectedActions, tc.registry.Actions, "Actions do not match")
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
			location: core.NewLocation("testLocation", core.Coord{}),
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				core.NewLocation("testLocation", core.Coord{}),
			},
			expectedActions: []core.Action{},
		},
		"can register location, attribute with location": {
			location: core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsLoc)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
				},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				core.NewLocation("loc1", core.Coord{}),
				core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsLoc)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{Location: core.NewLocation("loc1", core.Coord{})}),
			},
		},
		"can register location, attribute needs resource": {
			location: core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsRes)),
			registry: &Registry{
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{
					core.NewResource("res1"),
				},
			},
			expectedLocations: []*core.Location{
				core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsRes)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					core.CreateActionParams{
						Resource: core.NewResource("res1"),
					},
				),
			},
		},
		"can register location, attribute needs both": {
			location: core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsBoth)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1"),
					core.NewResource("res2"),
				},
			},
			expectedLocations: []*core.Location{
				core.NewLocation("loc1", core.Coord{}),
				core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsBoth)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: core.NewLocation("loc1", core.Coord{}),
					Resource: core.NewResource("res1"),
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: core.NewLocation("loc1", core.Coord{}),
					Resource: core.NewResource("res2"),
				}),
			},
		},
		"can register location, two attributes": {
			location: core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsBoth, mockAttributeNeedsRes)),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("loc1", core.Coord{}),
				},
				Resources: []*core.Resource{
					core.NewResource("res1"),
					core.NewResource("res2"),
				},
			},
			expectedLocations: []*core.Location{
				core.NewLocation("loc1", core.Coord{}),
				core.NewLocation("testLocation", core.Coord{}, core.WithAttributes(mockAttributeNeedsBoth, mockAttributeNeedsRes)),
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: core.NewLocation("loc1", core.Coord{}),
					Resource: core.NewResource("res1"),
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Location: core.NewLocation("loc1", core.Coord{}),
					Resource: core.NewResource("res2"),
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Resource: core.NewResource("res1"),
				}),
				mockActionWithResourceAndLocationCreateFunc(core.CreateActionParams{
					Resource: core.NewResource("res2"),
				}),
			},
		},
		"fail to register location": {
			location: core.NewLocation("testLocation", core.Coord{}),
			registry: &Registry{
				Actions: []core.Action{},
				Locations: []*core.Location{
					core.NewLocation("testLocation", core.Coord{}),
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
