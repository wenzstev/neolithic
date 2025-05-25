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

	for _, attr := range resource.Attributes().List() {
		if err := r.createActionsForAttribute(resource, attr); err != nil {
			return err
		}
	}

	for _, loc := range r.Locations {
		for _, locAttr := range loc.Attributes().List() {
			if locAttr.NeedsResource() {
				switch {
				case locAttr.NeedsLocation():
					for _, secondLoc := range r.Locations {
						if err := r.createAndAddAction(loc, locAttr, core.CreateActionParams{
							Location: secondLoc,
							Resource: resource,
						}); err != nil {
							return err
						}
					}
				default:
					if err := r.createAndAddAction(loc, locAttr, core.CreateActionParams{
						Resource: resource,
					}); err != nil {
						return err
					}
				}
			}
		}
	}

	for _, res := range r.Resources {
		for _, resAttr := range res.Attributes().List() {
			if resAttr.NeedsResource() {
				switch {
				case resAttr.NeedsLocation():
					for _, loc := range r.Locations {
						if err := r.createAndAddAction(res, resAttr, core.CreateActionParams{
							Location: loc,
							Resource: resource,
						}); err != nil {
							return err
						}
					}
				default:
					if err := r.createAndAddAction(res, resAttr, core.CreateActionParams{
						Resource: resource,
					}); err != nil {
						return err
					}
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

	for _, attr := range location.Attributes().List() {
		if err := r.createActionsForAttribute(location, attr); err != nil {
			return err
		}
	}

	for _, loc := range r.Locations {
		for _, locAttr := range loc.Attributes().List() {
			if locAttr.NeedsLocation() {
				switch {
				case locAttr.NeedsResource():
					for _, res := range r.Resources {
						if err := r.createAndAddAction(loc, locAttr, core.CreateActionParams{
							Location: location,
							Resource: res,
						}); err != nil {
							return err
						}
					}
				default:
					if err := r.createAndAddAction(loc, locAttr, core.CreateActionParams{
						Location: location,
					}); err != nil {
						return err
					}
				}
			}
		}
	}

	for _, res := range r.Resources {
		for _, resAttr := range res.Attributes().List() {
			if resAttr.NeedsLocation() {
				switch {
				case resAttr.NeedsResource():
					for _, res2 := range r.Resources {
						if err := r.createAndAddAction(res, resAttr, core.CreateActionParams{
							Resource: res2,
							Location: location,
						}); err != nil {
							return err
						}
					}
				default:
					if err := r.createAndAddAction(res, resAttr, core.CreateActionParams{
						Location: location,
					}); err != nil {
						return err
					}
				}

			}
		}
	}

	r.Locations = append(r.Locations, location)
	return nil
}

func (r *Registry) createActionsForAttribute(holder core.AttributeHolder, attr core.Attribute) error {
	switch {
	case attr.NeedsResource() && attr.NeedsLocation():
		for _, loc := range r.Locations {
			for _, res := range r.Resources {
				if err := r.createAndAddAction(holder, attr, core.CreateActionParams{
					Location: loc,
					Resource: res,
				}); err != nil {
					return err
				}
			}
		}
	case attr.NeedsResource() && !attr.NeedsLocation():
		for _, res := range r.Resources {
			if err := r.createAndAddAction(holder, attr, core.CreateActionParams{
				Resource: res,
			}); err != nil {
				return err
			}
		}
	case !attr.NeedsResource() && attr.NeedsLocation():
		for _, loc := range r.Locations {
			if err := r.createAndAddAction(holder, attr, core.CreateActionParams{
				Location: loc,
			}); err != nil {
				return err
			}
		}
	default:
		return r.createAndAddAction(holder, attr, core.CreateActionParams{})
	}
	return nil
}

func (r *Registry) createAndAddAction(holder core.AttributeHolder, attribute core.Attribute, params core.CreateActionParams) error {
	action, err := attribute.CreateAction(holder, params)
	if err != nil {
		return err
	}
	if action != nil {
		r.Actions = append(r.Actions, action)
	}
	return nil
}
