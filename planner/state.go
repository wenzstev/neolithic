package planner

import (
	"bytes"
	"crypto/sha256"
	"encoding/gob"
	"fmt"
)

// State represents the state of the world
type State struct {
	Locations map[*Location]Inventory
	Agents    map[*Agent]Inventory
}

// Copy performs a deep copy of the state
func (s *State) Copy() *State {
	end := &State{
		Locations: make(map[*Location]Inventory),
		Agents:    make(map[*Agent]Inventory),
	}

	for k, v := range s.Locations {
		end.Locations[k] = v.Copy()
	}

	for k, v := range s.Agents {
		end.Agents[k] = v.Copy()
	}

	return end
}

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
		output += fmt.Sprintf("   %s: \n", k.Name)
		output += v.String()
	}

	return output
}

func (s *State) ID() (string, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(s); err != nil {
		return "", err
	}
	hash := sha256.Sum256(buf.Bytes())
	return fmt.Sprintf("%x", hash), nil
}
