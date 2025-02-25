package planner

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
	"sort"
)

func init() {
	gob.Register(&mockAgent{})
}

// State represents the State of the world
type State struct {
	// Locations is a map of all Location in the world and their inventories
	Locations map[*Location]Inventory
	// Agents is a map of all Agent in the world and their inventory
	Agents map[Agent]Inventory
}

// Copy performs a deep copy of the State
func (s *State) Copy() *State {
	end := &State{
		Locations: make(map[*Location]Inventory),
		Agents:    make(map[Agent]Inventory),
	}

	for k, v := range s.Locations {
		end.Locations[k] = v.Copy()
	}

	for k, v := range s.Agents {
		end.Agents[k] = v.Copy()
	}

	return end
}

// String returns a string representation of the State
func (s *State) String() string {
	if s == nil {
		return "<nil>"
	}
	output := "State: \n"
	output += "  Locations: \n"
	for k, v := range s.Locations {
		output += fmt.Sprintf("   %s: \n", k.Name)
		output += v.String()
	}

	output += "  Agents:\n"
	for k, v := range s.Agents {
		output += fmt.Sprintf("   %s: \n", k.Name())
		output += v.String()
	}

	return output
}

// Add takes two states and adds them together. Note that these states may be ILLEGAL, in that they do not strictly
// preserve resources. Thus it is currently not possible to just add two states together to fulfill Action.Perform
func (s *State) Add(other *State, clamped bool) *State {
	newState := s.Copy()
	for loc, inv := range other.Locations {
		if _, ok := newState.Locations[loc]; !ok {
			newState.Locations[loc] = Inventory{}
		}
		for res, change := range inv {
			if _, ok := newState.Locations[loc][res]; !ok {
				newState.Locations[loc][res] = 0
			}
			newState.Locations[loc][res] += change
			if !clamped && newState.Locations[loc][res] < 0 {
				delete(newState.Locations[loc], res)
			}
		}
	}
	for agent, inv := range other.Agents {
		if _, ok := newState.Agents[agent]; !ok {
			newState.Agents[agent] = Inventory{}
		}
		for res, amount := range inv {
			if _, ok := newState.Agents[agent][res]; !ok {
				newState.Agents[agent][res] = 0
			}
			newState.Agents[agent][res] += amount
			if !clamped && newState.Agents[agent][res] < 0 {
				delete(newState.Agents[agent], res)
			}
		}
	}
	return newState
}

// ID returns a unique ID hash of the State. Turns all maps into slices to ensure consistent ordering.
func (s *State) ID() (string, error) {

	type invStruct struct {
		Resource *Resource
		Amount   int
	}

	type locStruct struct {
		Location  *Location
		Inventory []invStruct
	}

	type agentStruct struct {
		Agent Agent
		Inv   []invStruct
	}

	type sortedState struct {
		Locations []locStruct
		Agents    []agentStruct
	}

	sortInv := func(inv Inventory) []invStruct {
		items := make([]invStruct, 0, len(inv))
		for res, count := range inv {
			items = append(items, invStruct{res, count})
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Resource.Name < items[j].Resource.Name
		})
		return items
	}

	sortLoc := func(locs map[*Location]Inventory) []locStruct {
		items := make([]locStruct, 0, len(locs))
		for loc, inv := range locs {
			items = append(items, locStruct{loc, sortInv(inv)})
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Location.Name < items[j].Location.Name
		})
		return items
	}

	sortAgent := func(agents map[Agent]Inventory) []agentStruct {
		items := make([]agentStruct, 0, len(agents))
		for agent, inv := range agents {
			items = append(items, agentStruct{agent, sortInv(inv)})
		}
		sort.Slice(items, func(i, j int) bool {
			return items[i].Agent.Name() < items[j].Agent.Name()
		})
		return items
	}

	stateToEncode := sortedState{
		Locations: sortLoc(s.Locations),
		Agents:    sortAgent(s.Agents),
	}

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(stateToEncode); err != nil {
		return "", err
	}
	hash := sha256.Sum256(buf.Bytes())
	return fmt.Sprintf("%x", hash), nil
}
