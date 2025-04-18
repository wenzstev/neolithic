package core

import (
	"hash/fnv"
	"strconv"
)

// WorldState represents the current state of the simulation world.
type WorldState struct {
	// Grid represents the world's grid.
	Grid Grid
	// Locations is a map of locations in the world.
	Locations []Location
	// Agents is a map of agents in the world.
	Agents   []Agent
	cachedID string
}

// ID returns a unique identifier for the WorldState. It uses SHA256 to generate a hash of the state.
func (w *WorldState) ID() (string, error) {
	if w.cachedID != "" {
		return w.cachedID, nil
	}
	h := fnv.New64a()
	for _, loc := range w.Locations {
		h.Write([]byte(loc.String()))
	}
	for _, agent := range w.Agents {
		h.Write([]byte(agent.String()))
	}

	//w.cachedID = strconv.FormatUint(h.Sum64(), 16)
	return strconv.FormatUint(h.Sum64(), 16), nil

}

// DeepCopy creates a deep copy of the WorldState.
// TODO: worldstate will probably need to become slices to cut down on copy time
func (w *WorldState) DeepCopy() *WorldState {
	end := &WorldState{
		Locations: make([]Location, len(w.Locations)),
		Agents:    make([]Agent, len(w.Agents)),
	}
	for i := 0; i < len(end.Locations); i++ {
		end.Locations[i] = *w.Locations[i].DeepCopy()
	}
	for i := 0; i < len(end.Agents); i++ {
		end.Agents[i] = w.Agents[i].DeepCopy()
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

func (w *WorldState) GetLocation(name string) (*Location, bool) {
	for _, loc := range w.Locations {
		if loc.Name == name {
			return &loc, true
		}
	}
	return nil, false
}

func (w *WorldState) GetAgent(name string) (Agent, bool) {
	for _, agent := range w.Agents {
		if agent.Name() == name {
			return agent, true
		}
	}
	return nil, false
}
