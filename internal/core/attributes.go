package core

// CreateActionParams are the parameters used for the CreateAction function.
type CreateActionParams struct {
	// Location is the location to be passed to the CreateAction func
	Location *Location
	// Resource is the resource to be passed to the CreatAction func
	Resource *Resource
}

// Attribute is a representation of a unique aspect of a Resource or a Location
type Attribute interface {
	// CreateAction provides a way to create an action from an interface
	CreateAction(CreateActionParams) Action
	// NeedsLocation determines if an attribute needs a separate location to be relevant
	NeedsLocation() bool
	// NeedsResource determines if an attribute needs a separate resource to be relevant
	NeedsResource() bool
}
