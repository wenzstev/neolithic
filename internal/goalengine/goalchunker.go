package goalengine

import (
	"Neolithic/internal/core"
)

const DefaultIncreaseAmount = 100

// GoalEngine is a struct that encapsulates the Goal object and provides mechanisms to manage and process it.
type GoalEngine struct {
	// Goal is the current goal of the GoalEngine
	Goal Goal
}

// Goal defines a specific objective and includes its name, logic, target location, and associated resource.
type Goal struct {
	// Name is the name of the goal
	Name string
	// Logic encapsulates the logic for managing the goal
	Logic GoalLogic
	// Location is the location that relates to the goal
	Location *core.Location
	// Resource is the resource that relates to the goal
	Resource *core.Resource
}

// GoalLogic represents the logic for managing a goal, including chunking, fallback, and termination conditions.
type GoalLogic struct {
	// ID represents a unique identifier for the GoalLogic
	ID string // ID for gob
	// Chunker is the function used to break the goal into chunks
	Chunker ChunkerFunc
	// Fallback is the function used to decrease the default chunk's ambition
	Fallback FallbackChunk
	// ShouldGiveUp is used to determine if the goal is no longer worth pursuing
	ShouldGiveUp ShouldGiveUp
}

// ChunkerFunc is the function used to create a defualt chunk of a goal
type ChunkerFunc func(*core.Location, *core.Resource) *core.WorldState

// AddToLocation is a ChunkerFunc that adds the default increase amount to the location's inventory
var AddToLocation ChunkerFunc = func(location *core.Location, resource *core.Resource) *core.WorldState {
	goalLocation := core.Location{
		Name:      location.Name,
		Inventory: core.NewInventory(),
	}
	goalLocation.Inventory.AdjustAmount(resource, DefaultIncreaseAmount)

	return &core.WorldState{
		Locations: map[string]core.Location{
			location.Name: goalLocation,
		},
	}
}

// FallbackChunk is the function used to create a fallback chunk of a goal
type FallbackChunk func(*core.WorldState) *core.WorldState

// FallbackChunkFunc is a FallbackChunk that halves the amount of each resource in the world state
var FallbackChunkFunc FallbackChunk = func(worldState *core.WorldState) *core.WorldState {
	worldStateCopy := worldState.DeepCopy()
	for _, location := range worldStateCopy.Locations {
		for _, entry := range location.Inventory.Entries() {
			location.Inventory.AdjustAmount(entry.Resource, -entry.Amount/2)
		}
	}

	return worldStateCopy
}

// ShouldGiveUp is the function used to determine if the goal is no longer worth pursuing
type ShouldGiveUp func(*core.WorldState) bool

// GiveUpIfLessThanFive is a ShouldGiveUp that returns true if the total amount of resources in the world state is less than five
var GiveUpIfLessThanFive ShouldGiveUp = func(worldState *core.WorldState) bool {
	totalResources := 0
	for _, location := range worldState.Locations {
		for _, entry := range location.Inventory.Entries() {
			totalResources += entry.Amount
		}
	}
	return totalResources < 5
}

// GiveUpIfNoChange is a ShouldGiveUp that returns true if there are no changes to the world state
var GiveUpIfNoChange ShouldGiveUp = func(worldState *core.WorldState) bool {
	totalResources := 0
	for _, location := range worldState.Locations {
		for _, entry := range location.Inventory.Entries() {
			totalResources += entry.Amount
		}
	}
	return totalResources < 1
}

// GetDelta returns the delta for the goal; that is, the change in amount. It does not return a full WorldState
func (g *Goal) GetDelta(numRetries int) *core.WorldState {
	chunk := g.Logic.Chunker(g.Location, g.Resource)

	for i := 0; i < numRetries; i++ {
		chunk = g.Logic.Fallback(chunk)
		if g.Logic.ShouldGiveUp(chunk) {
			return nil
		}
	}

	return chunk
}

// GetGoalChunk takes in the current state of the world and returns a chunked goal for that world, based on the Goal's
// overarching requirements.
func (g *Goal) GetGoalChunk(state *core.WorldState, numRetries int) *core.WorldState {
	// Get the delta based on number of retries
	delta := g.GetDelta(numRetries)
	if delta == nil {
		return nil
	}

	// Create a deep copy of the current state
	result := state.DeepCopy()

	// Apply the delta to the location in the result
	for _, deltaLoc := range delta.Locations {
		if resultLoc, exists := result.Locations[deltaLoc.Name]; exists {
			deltaLoc.Inventory.AdjustAmount(g.Resource, resultLoc.Inventory.GetAmount(g.Resource))
		} else {
			result.Locations[deltaLoc.Name] = deltaLoc
		}
	}

	return delta
}

func (g *GoalEngine) GetNextGoal(worldState *core.WorldState, retries int) *core.WorldState {
	return g.Goal.GetGoalChunk(worldState, retries)
}
