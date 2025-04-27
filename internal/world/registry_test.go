package world

import (
	"Neolithic/internal/core"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRegistry_RegisterAction(t *testing.T) {
	type testCase struct {
		actionName      string
		action          core.Action
		createFunc      ActionCreator
		registry        *Registry
		expectedActions []core.Action
		expectedEntry   *ActionRegistryEntry
		expectedError   error
	}

	tests := map[string]testCase{
		"can register action": {
			actionName: "mockAction",
			action:     &mockAction{},
			createFunc: mockActionCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedActions: []core.Action{},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockAction",
				NeedsLocation: false,
				NeedsResource: false,
				Creator:       mockActionCreateFunc,
			},
		},
		"can register location action": {
			actionName: "mockLocationAction",
			action:     &mockActionWithLocation{},
			createFunc: mockActionWithLocationCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedActions: []core.Action{},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockLocationAction",
				NeedsLocation: true,
				NeedsResource: false,
				Creator:       mockActionWithLocationCreateFunc,
			},
		},
		"can register resource action": {
			actionName: "mockResourceAction",
			action:     &mockActionWithLocation{},
			createFunc: mockActionWithResourceCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedActions: []core.Action{},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockResourceAction",
				NeedsLocation: true,
				NeedsResource: false,
				Creator:       mockActionWithResourceCreateFunc,
			},
		},
		"can register resource and location action": {
			actionName: "mockLocationResourceAction",
			action:     &mockActionWithLocationAndResource{},
			createFunc: mockActionWithResourceAndLocationCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedActions: []core.Action{},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockLocationResourceAction",
				NeedsLocation: true,
				NeedsResource: true,
				Creator:       mockActionWithResourceAndLocationCreateFunc,
			},
		},
		"fail to register action": {
			actionName: "mockAction",
			action:     &mockAction{},
			createFunc: mockActionCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name: "mockAction",
					},
				},
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedActions: []core.Action{},
			expectedEntry:   nil,
			expectedError:   ErrActionAlreadyRegistered,
		},
		"can register location action with locations": {
			actionName: "mockLocationAction",
			action:     &mockActionWithLocation{},
			createFunc: mockActionWithLocationCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations: []*core.Location{
					{Name: "location1"},
					{Name: "location2"},
				},
				Resources: []*core.Resource{},
			},
			expectedActions: []core.Action{
				mockActionWithLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location1"},
				}),
				mockActionWithLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location2"},
				}),
			},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockLocationAction",
				NeedsLocation: true,
				NeedsResource: false,
				Creator:       mockActionWithLocationCreateFunc,
			},
		},
		"can register resource action with resources": {
			actionName: "mockResourceAction",
			action:     &mockActionWithResource{},
			createFunc: mockActionWithResourceCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources: []*core.Resource{
					{Name: "resource1"},
					{Name: "resource2"},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceCreateFunc(ActionCreatorParams{
					Resource: &core.Resource{Name: "resource1"},
				}),
				mockActionWithResourceCreateFunc(ActionCreatorParams{
					Resource: &core.Resource{Name: "resource2"},
				}),
			},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockResourceAction",
				NeedsLocation: false,
				NeedsResource: true,
				Creator:       mockActionCreateFunc,
			},
		},
		"can register location/resource action with both": {
			actionName: "mockLocationResourceAction",
			action:     &mockActionWithLocationAndResource{},
			createFunc: mockActionWithResourceAndLocationCreateFunc,
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations: []*core.Location{
					{Name: "location1"},
					{Name: "location2"},
				},
				Resources: []*core.Resource{
					{Name: "resource1"},
					{Name: "resource2"},
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location1"},
					Resource: &core.Resource{Name: "resource1"},
				}),
				mockActionWithResourceAndLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location1"},
					Resource: &core.Resource{Name: "resource2"},
				}),
				mockActionWithResourceAndLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location2"},
					Resource: &core.Resource{Name: "resource1"},
				}),
				mockActionWithResourceAndLocationCreateFunc(ActionCreatorParams{
					Location: &core.Location{Name: "location2"},
					Resource: &core.Resource{Name: "resource2"},
				}),
			},
			expectedEntry: &ActionRegistryEntry{
				Name:          "mockLocationResourceAction",
				NeedsLocation: true,
				NeedsResource: true,
				Creator:       mockActionWithLocationCreateFunc,
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			err := tc.registry.RegisterAction(tc.actionName, tc.action, tc.createFunc)
			if tc.expectedError != nil {
				require.ErrorIs(t, err, tc.expectedError)
				return
			}
			require.NoError(t, err)
			require.Equal(t, tc.expectedActions, tc.registry.Actions)

			var entry *ActionRegistryEntry
			for _, e := range tc.registry.ActionRegistry {
				if e.Name == tc.actionName {
					entry = e
					break
				}
			}

			require.Equal(t, tc.expectedEntry.NeedsLocation, entry.NeedsLocation)
			require.Equal(t, tc.expectedEntry.NeedsResource, entry.NeedsResource)
			require.Equal(t, tc.expectedEntry.Name, entry.Name)

		})
	}
}

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
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
				},
			},
			expectedActions: []core.Action{},
		},
		"can register resource, actions with only locations": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "locationAction",
						NeedsLocation: true,
						NeedsResource: false,
						Creator:       mockActionWithLocationCreateFunc,
					},
				},
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
				},
			},
			expectedActions: []core.Action{},
		},
		"can register resource, actions with only resources": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "action2",
						NeedsLocation: false,
						NeedsResource: true,
						Creator:       mockActionWithResourceCreateFunc,
					},
				},
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceCreateFunc(ActionCreatorParams{Resource: &core.Resource{Name: "testResource"}}),
			},
		},
		"can register resource, actions with both": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "resourceLocationAction",
						NeedsLocation: true,
						NeedsResource: true,
						Creator:       mockActionWithResourceAndLocationCreateFunc,
					},
				},
				Actions: []core.Action{},
				Locations: []*core.Location{
					{Name: "loc1"},
					{Name: "loc2"},
				},
				Resources: []*core.Resource{},
			},
			expectedResources: []*core.Resource{
				{
					Name: "testResource",
				},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					ActionCreatorParams{
						Resource: &core.Resource{Name: "testResource"},
						Location: &core.Location{Name: "loc1"},
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					ActionCreatorParams{
						Resource: &core.Resource{Name: "testResource"},
						Location: &core.Location{Name: "loc2"},
					},
				),
			},
		},
		"fail to register resource": {
			resource: &core.Resource{
				Name: "testResource",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
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
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
				Locations:      []*core.Location{},
				Resources:      []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				{Name: "testLocation"},
			},
			expectedActions: []core.Action{},
		},
		"can register location, actions with only locations": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "locationAction",
						NeedsLocation: true,
						NeedsResource: false,
						Creator:       mockActionWithLocationCreateFunc,
					},
				},
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				{Name: "testLocation"},
			},
			expectedActions: []core.Action{
				mockActionWithLocationCreateFunc(ActionCreatorParams{Location: &core.Location{Name: "testLocation"}}),
			},
		},
		"can register location, actions with only resources": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "action2",
						NeedsLocation: false,
						NeedsResource: true,
						Creator:       mockActionWithResourceCreateFunc,
					},
				},
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{},
			},
			expectedLocations: []*core.Location{
				{Name: "testLocation"},
			},
			expectedActions: []core.Action{},
		},
		"can register location, actions with both": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{
					{
						Name:          "resourceLocationAction",
						NeedsLocation: true,
						NeedsResource: true,
						Creator:       mockActionWithResourceAndLocationCreateFunc,
					},
				},
				Actions:   []core.Action{},
				Locations: []*core.Location{},
				Resources: []*core.Resource{
					{Name: "res1"},
					{Name: "res2"},
				},
			},
			expectedLocations: []*core.Location{
				{Name: "testLocation"},
			},
			expectedActions: []core.Action{
				mockActionWithResourceAndLocationCreateFunc(
					ActionCreatorParams{
						Location: &core.Location{Name: "testLocation"},
						Resource: &core.Resource{Name: "res1"},
					},
				),
				mockActionWithResourceAndLocationCreateFunc(
					ActionCreatorParams{
						Location: &core.Location{Name: "testLocation"},
						Resource: &core.Resource{Name: "res2"},
					},
				),
			},
		},
		"fail to register location": {
			location: &core.Location{
				Name: "testLocation",
			},
			registry: &Registry{
				ActionRegistry: []*ActionRegistryEntry{},
				Actions:        []core.Action{},
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
