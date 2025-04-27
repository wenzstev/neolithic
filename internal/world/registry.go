package world

import (
	"errors"

	"Neolithic/internal/core"
)

var (
	// ErrResourceAlreadyRegistered indicates that a resource has already been registered.
	ErrResourceAlreadyRegistered = errors.New("resource already registered")

	// ErrLocationAlreadyRegistered indicates that a location has already been registered.
	ErrLocationAlreadyRegistered = errors.New("location already registered")

	// ErrActionAlreadyRegistered indicates that an action has already been registered.
	ErrActionAlreadyRegistered = errors.New("action already registered")
)

// ActionCreatorParams contains parameters required for creating an action, including a Location and a Resource.
type ActionCreatorParams struct {
	Location *core.Location
	Resource *core.Resource
}

// ActionCreator defines a function type that generates a planner.Action based on the provided ActionCreatorParams.
type ActionCreator func(params ActionCreatorParams) core.Action

// Registry manages the registration of actions, locations, and resources in the world.
type Registry struct {
	// ActionRegistry holds all actions that are currently registered in the registry
	ActionRegistry []*ActionRegistryEntry
	// Actions holds all instantiated actions that can be performed.
	Actions []core.Action
	// Locations holds all locations in the registry
	Locations []*core.Location
	// Resources holds all resources in the registry
	Resources []*core.Resource
}

// ActionRegistry is a mapping of action names to their corresponding ActionRegistryEntry configurations.
type ActionRegistry map[string]*ActionRegistryEntry

// ActionRegistryEntry represents a registry entry for an action with associated metadata and creation logic.
type ActionRegistryEntry struct {
	// Name is the generic name of the action
	Name string
	// NeedsLocation indicates if the action requires a location to be performed
	NeedsLocation bool
	// NeedsResource indicates if the action requires a resource to be performed
	NeedsResource bool
	// Creator is the function that is used to generate an instance of the action
	Creator ActionCreator
}

// RegisterAction registers a new action with a name, checks for duplicates, and creates instances based on dependencies.
func (r *Registry) RegisterAction(name string, action core.Action, createFunc ActionCreator) error {
	for _, entry := range r.ActionRegistry {
		if entry.Name == name {
			return ErrActionAlreadyRegistered
		}
	}

	_, locatable := action.(core.Locatable)
	_, needsResource := action.(core.NeedsResource)

	entry := &ActionRegistryEntry{
		Name:          name,
		NeedsLocation: locatable,
		NeedsResource: needsResource,
		Creator:       createFunc,
	}

	if entry.NeedsLocation {
		for _, loc := range r.Locations {
			if entry.NeedsResource {
				for _, res := range r.Resources {
					newAction := entry.Creator(ActionCreatorParams{
						Location: loc,
						Resource: res,
					})
					r.Actions = append(r.Actions, newAction)
				}
			} else {
				newAction := entry.Creator(ActionCreatorParams{
					Location: loc,
				})
				r.Actions = append(r.Actions, newAction)
			}
		}
	} else if entry.NeedsResource {
		for _, res := range r.Resources {
			newAction := entry.Creator(ActionCreatorParams{
				Resource: res,
			})
			r.Actions = append(r.Actions, newAction)
		}
	}
	r.ActionRegistry = append(r.ActionRegistry, entry)
	return nil
}

// RegisterResource registers a new resource in the registry. Returns an error if the resource is already registered.
// It also creates actions for the resource, depending on the registered actions requiring resources and locations.
func (r *Registry) RegisterResource(resource *core.Resource) error {
	for _, res := range r.Resources {
		if res.Name == resource.Name {
			return ErrResourceAlreadyRegistered
		}
	}
	for _, action := range r.ActionRegistry {
		if action.NeedsResource {
			if action.NeedsLocation {
				for _, loc := range r.Locations {
					r.Actions = append(r.Actions, action.Creator(ActionCreatorParams{
						Location: loc,
						Resource: resource,
					}))
				}
				continue
			}
			r.Actions = append(r.Actions, action.Creator(ActionCreatorParams{
				Resource: resource,
			}))
		}
	}
	r.Resources = append(r.Resources, resource)
	return nil
}

// RegisterLocation registers a new location in the registry and associates it with applicable actions and resources.
// Returns ErrLocationAlreadyRegistered if the location is already registered.
func (r *Registry) RegisterLocation(location *core.Location) error {
	for _, loc := range r.Locations {
		if loc.Name == location.Name {
			return ErrLocationAlreadyRegistered
		}
	}

	for _, action := range r.ActionRegistry {
		if action.NeedsLocation {
			if action.NeedsResource {
				for _, res := range r.Resources {
					r.Actions = append(r.Actions, action.Creator(ActionCreatorParams{
						Location: location,
						Resource: res,
					}))
				}
				continue
			}
			r.Actions = append(r.Actions, action.Creator(ActionCreatorParams{
				Location: location,
			}))
		}
	}
	r.Locations = append(r.Locations, location)
	return nil
}
