package world

import (
	"Neolithic/internal/core"
	"errors"
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

// RegisterResource registers a new resource in the registry. Returns an error if the resource is already registered.
// It also creates actions for the resource, depending on the registered actions requiring resources and locations.
func (r *Registry) RegisterResource(resource *core.Resource) error {
	for _, res := range r.Resources {
		if res.Name == resource.Name {
			return ErrResourceAlreadyRegistered
		}
	}

	for _, attr := range resource.Attributes {
		r.createActionsForAttribute(attr)
	}

	for _, loc := range r.Locations {
		for _, locAttr := range loc.Attributes {
			if locAttr.NeedsResource() {
				switch {
				case locAttr.NeedsLocation():
					for _, secondLoc := range r.Locations {
						r.Actions = append(r.Actions, locAttr.CreateAction(core.CreateActionParams{
							Location: secondLoc,
							Resource: resource,
						}))
					}
				default:
					r.Actions = append(r.Actions, locAttr.CreateAction(core.CreateActionParams{
						Resource: resource,
					}))
				}
			}
		}
	}

	for _, res := range r.Resources {
		for _, resAttr := range res.Attributes {
			if resAttr.NeedsResource() {
				switch {
				case resAttr.NeedsLocation():
					for _, loc := range r.Locations {
						r.Actions = append(r.Actions, resAttr.CreateAction(core.CreateActionParams{
							Location: loc,
							Resource: resource,
						}))
					}
				default:
					r.Actions = append(r.Actions, resAttr.CreateAction(core.CreateActionParams{
						Resource: resource,
					}))
				}

			}
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

	for _, attr := range location.Attributes {
		r.createActionsForAttribute(attr)
	}

	for _, loc := range r.Locations {
		for _, locAttr := range loc.Attributes {
			if locAttr.NeedsLocation() {
				switch {
				case locAttr.NeedsResource():
					for _, res := range r.Resources {
						r.Actions = append(r.Actions, locAttr.CreateAction(core.CreateActionParams{
							Location: location,
							Resource: res,
						}))
					}
				default:
					r.Actions = append(r.Actions, locAttr.CreateAction(core.CreateActionParams{
						Location: location,
					}))
				}
			}
		}
	}

	for _, res := range r.Resources {
		for _, resAttr := range res.Attributes {
			if resAttr.NeedsLocation() {
				switch {
				case resAttr.NeedsResource():
					for _, res2 := range r.Resources {
						r.Actions = append(r.Actions, resAttr.CreateAction(core.CreateActionParams{
							Location: location,
							Resource: res2,
						}))
					}
				default:
					r.Actions = append(r.Actions, resAttr.CreateAction(core.CreateActionParams{
						Location: location,
					}))
				}

			}
		}
	}

	r.Locations = append(r.Locations, location)
	return nil
}

func (r *Registry) createActionsForAttribute(attr core.Attribute) {
	switch {
	case attr.NeedsResource() && attr.NeedsLocation():
		for _, loc := range r.Locations {
			for _, res := range r.Resources {
				r.Actions = append(r.Actions, attr.CreateAction(core.CreateActionParams{
					Location: loc,
					Resource: res,
				}))
			}
		}
	case attr.NeedsResource() && !attr.NeedsLocation():
		for _, res := range r.Resources {
			r.Actions = append(r.Actions, attr.CreateAction(core.CreateActionParams{
				Resource: res,
			}))
		}
	case !attr.NeedsResource() && attr.NeedsLocation():
		for _, loc := range r.Locations {
			r.Actions = append(r.Actions, attr.CreateAction(core.CreateActionParams{
				Location: loc,
			}))
		}
	default:
		r.Actions = append(r.Actions, attr.CreateAction(core.CreateActionParams{}))
	}
}
