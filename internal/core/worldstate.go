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
	Locations map[string]*Location
	// Agents is a map of agents in the world.
	Agents   map[string]Agent
	cachedID string
}

// ID returns a unique identifier for the WorldState. It uses SHA256 to generate a hash of the state.
func (w *WorldState) ID() (string, error) {
	if w.cachedID != "" {
		return w.cachedID, nil
	}
	h := fnv.New64a()

	locKeys := getSortedLocationKeys(w.Locations)
	for _, k := range locKeys {
		h.Write([]byte(w.Locations[k].String()))
	}

	agentKeys := getSortedAgentKeys(w.Agents)
	for _, k := range agentKeys {
		h.Write([]byte(w.Agents[k].String()))
	}

	w.cachedID = strconv.FormatUint(h.Sum64(), 16)
	return w.cachedID, nil

}

// DeepCopy creates a deep copy of the WorldState.
func (w *WorldState) DeepCopy() *WorldState {
	end := &WorldState{
		Locations: make(map[string]*Location, len(w.Locations)),
		Agents:    make(map[string]Agent, len(w.Agents)),
	}
	for k, v := range w.Locations {
		end.Locations[k] = v.DeepCopy()
	}
	for k, v := range w.Agents {
		end.Agents[k] = v.DeepCopy()
	}

	return end
}

// ShallowCopy creates a shallow copy of the world state; all Locations, Agents and the grid are the same
func (w *WorldState) ShallowCopy() *WorldState {
	newState := &WorldState{
		Grid:      w.Grid,
		Locations: make(map[string]*Location, len(w.Locations)),
		Agents:    make(map[string]Agent, len(w.Agents)),
	}
	for k, v := range w.Locations {
		newState.Locations[k] = v
	}
	for k, v := range w.Agents {
		newState.Agents[k] = v
	}
	return newState
}

func getSortedLocationKeys(m map[string]*Location) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func getSortedAgentKeys(m map[string]Agent) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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

func (w *WorldState) GetLocation(name string) (*Location, bool) {
	loc, ok := w.Locations[name]
	return loc, ok
}
func (w *WorldState) GetAgent(name string) (Agent, bool) {
	agent, ok := w.Agents[name]
	return agent, ok
}
