package planner

import "fmt"

// AStarOutput represents the output of the A* function, as in the set of actions determined to bring the starting state
// to the goal state
type AStarOutput struct {
	actions       []Action
	totalCost     float64
	expectedState *State
}

// String returns a string representation of the AStarOutput
func (a *AStarOutput) String() string {
	output := ""
	output += fmt.Sprintf("Total Cost: %f\n", a.totalCost)
	output += "Actions: \n"
	end := len(a.actions) - 1
	for i := range a.actions {
		output += fmt.Sprintf("%s\n", a.actions[end-i].Description())
	}
	return output
}
