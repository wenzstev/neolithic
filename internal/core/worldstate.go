package core

import (
	"hash/fnv"
	"sort"
	"strconv"
)

// WorldState represents the current state of the simulation world.
type WorldState struct {
	// Grid represents the world's grid.
	Grid Grid
	// Locations is a map of locations in the world.
	Locations map[string]Location
	// Agents is a map of agents in the world.
	Agents   map[string]Agent
	cachedID string
}

// ID returns a unique identifier for the WorldState. It uses SHA256 to generate a hash of the state.
func (w *WorldState) ID() (string, error) {
	if w.cachedID != "" {
		return w.cachedID, nil
	}
	type sortedState struct {
		Locations []Location
		Agents    []Agent
	}

	sortLoc := func(locs map[string]Location) []Location {
		items := make([]Location, 0, len(locs))
		for _, loc := range locs {
			items = append(items, loc)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name < items[j].Name
		})
		return items
	}

	sortAgent := func(agents map[string]Agent) []Agent {
		items := make([]Agent, 0, len(agents))
		for _, agent := range agents {
			items = append(items, agent)
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Name() < items[j].Name()
		})
		return items
	}

	stateToEncode := sortedState{
		Locations: sortLoc(w.Locations),
		Agents:    sortAgent(w.Agents),
	}

	h := fnv.New64a()
	for _, loc := range stateToEncode.Locations {
		h.Write([]byte(loc.String()))
	}
	for _, agent := range stateToEncode.Agents {
		h.Write([]byte(agent.String()))
	}

	return strconv.FormatUint(h.Sum64(), 16), nil

}

// DeepCopy creates a deep copy of the WorldState.
// TODO: worldstate will probably need to become slices to cut down on copy time
func (w *WorldState) DeepCopy() *WorldState {
	end := &WorldState{
		Locations: make(map[string]Location),
		Agents:    make(map[string]Agent),
	}

	for k, v := range w.Locations {
		end.Locations[k] = *v.DeepCopy()
	}

	for k, v := range w.Agents {
		end.Agents[k] = v.DeepCopy()
	}

	return end
}

// String returns a string representation of the WorldState.
func (w *WorldState) String() string {
	if w == nil {
		return "<nil>"
	}
	output := "WorldState :\n Locations:\n"
	for _, loc := range w.Locations {
		output += loc.String()
	}
	output += "\n Agents:\n"
	for _, agent := range w.Agents {
		output += agent.String()
	}
	return output
}
