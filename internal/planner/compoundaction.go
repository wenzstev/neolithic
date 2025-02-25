package planner

import (
	"fmt"
)

// CompoundAction represents a sequence of actions. It implements Action as well, allowing it to be used as a single
// node like any other Action.
type CompoundAction []Action

// Cost implements Action.Cost and returns the combined cost of all the actions.
func (seq *CompoundAction) Cost(agent Agent) float64 {
	totalCost := 0.0
	for _, action := range *seq {
		totalCost += action.Cost(agent)
	}
	return totalCost
}

// Perform implements Action.Perform and returns the new State after all actions have been performed, or nil
// if the Action sequence results in an invalid State.
func (seq *CompoundAction) Perform(state *State, agent Agent) *State {
	curState := state
	for _, action := range *seq {
		curState = action.Perform(curState, agent)
		if curState == nil {
			return nil
		}
	}
	return curState
}

// Description implements Action.Description and returns a list of all descriptions of individual actions.
func (seq *CompoundAction) Description() string {
	output := "Sequence:\n"
	for _, action := range *seq {
		output += fmt.Sprintf("  %s\n", action.Description())
	}
	return output
}
