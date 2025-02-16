package planner

import (
	"Neolithic/internal/astar"
	"fmt"
	"math"
)

// GoapNode represents a point in a GOAP process, where the planner is choosing a plan
type GoapNode struct {
	action      Action
	state       *State
	goapRunInfo *GoapRunInfo
}

// GoapRunInfo represents the information that doesn't change across the GOAP planning call
type GoapRunInfo struct {
	Agent               *Agent
	PossibleNextActions *[]Action
}

var _ astar.Node = (*GoapNode)(nil)

// Heuristic implements astar.Node, and represents a best guess estimate of how far the
// given node is from the goal node.
func (g *GoapNode) Heuristic(goal astar.Node) (float64, error) {
	goapNode, ok := goal.(*GoapNode)
	if !ok {
		return 0, fmt.Errorf("GoapNode expected for goal")
	}
	return g.heuristic(g, goapNode)
}

// ID implements astar.Node and returns a unique string representing the node.
func (g *GoapNode) ID() (string, error) {
	return g.state.ID()
}

// Cost implements astar.Node and returns the cost of performing the acion associated with this node.
func (g *GoapNode) Cost(_ astar.Node) float64 {
	return g.action.Cost(g.goapRunInfo.Agent)
}

// GetSuccessors implements astar.Node and returns a list of successor astar.Node to this astar.Node.
func (g *GoapNode) GetSuccessors() ([]astar.Node, error) {
	successors := make([]astar.Node, 0)
	for _, action := range *g.goapRunInfo.PossibleNextActions {
		newState := action.Perform(g.state, g.goapRunInfo.Agent)
		if newState == nil {
			continue
		}
		successors = append(successors, &GoapNode{
			action:      action,
			state:       newState,
			goapRunInfo: g.goapRunInfo,
		})
	}
	return successors, nil
}

// heuristic is the function used to estimate how close to the goal a given action is. It does so by calculating the
// lowest "cost per unit" of all Action(s) that operates on a resource relevant to the goal. That value is then
// multiplied by the difference in amount of that resource between the current and the goal location.
// This heuristic is admissible because it always chooses the least "cost per unit" available, meaning it cannot
// overestimate the total cost of a given path.
func (g *GoapNode) heuristic(cur, goal *GoapNode) (float64, error) {
	var totalCost float64
	for loc, goalInventory := range goal.state.Locations {
		currentInventory, ok := cur.state.Locations[loc]
		for item, goalAmount := range goalInventory {
			currentAmount := 0
			if ok {
				currentAmount = currentInventory[item]
			}

			diff := goalAmount - currentAmount
			if diff == 0 {
				continue
			}

			requiredChange := math.Abs(float64(diff))
			var bestCostPerUnit = math.Inf(1)

			var relevantActions []Action
			var err error
			if diff > 0 {
				relevantActions, err = cur.getActionsThatAdd(item, loc)
				if err != nil {
					return math.Inf(1), err
				}
			} else {
				relevantActions, err = cur.getActionsThatRemove(item, loc)
				if err != nil {
					return math.Inf(1), err
				}
			}

			for _, action := range relevantActions {
				stateDiff := action.GetStateChange(cur.goapRunInfo.Agent)
				effectAmount := stateDiff.Locations[loc][item]

				costPerUnit := action.Cost(cur.goapRunInfo.Agent) / math.Abs(float64(effectAmount))
				if costPerUnit < bestCostPerUnit {
					bestCostPerUnit = costPerUnit
				}
			}
			totalCost += requiredChange * bestCostPerUnit
		}
	}
	return totalCost, nil
}

// getActionsThatAdd returns all actions that the Agent on the GoapNode can take that _add_ the given Resource to the given
// Location
func (g *GoapNode) getActionsThatAdd(res *Resource, loc *Location) ([]Action, error) {
	addActions := make([]Action, 0)

	successors, err := g.GetSuccessors()
	if err != nil {
		return nil, err
	}

	for _, successor := range successors {
		action := successor.(*GoapNode).action

		stateDiff := action.GetStateChange(g.goapRunInfo.Agent)
		if stateDiff.Locations[loc] == nil {
			continue
		}
		if stateDiff.Locations[loc][res] == 0 {
			continue
		}
		if stateDiff.Locations[loc][res] < 0 {
			continue
		}
		addActions = append(addActions, action)
	}
	return addActions, nil
}

// getActionsThatRemove returns all actions that the Agent on the GoapNode can take that _remove_ the given Resource
// from the given location.
func (g *GoapNode) getActionsThatRemove(res *Resource, loc *Location) ([]Action, error) {
	removeActions := make([]Action, 0)

	successors, err := g.GetSuccessors()
	if err != nil {
		return nil, err
	}

	for _, successor := range successors {
		action := successor.(*GoapNode).action

		stateDiff := action.GetStateChange(g.goapRunInfo.Agent)
		if stateDiff.Locations[loc] == nil {
			continue
		}
		if stateDiff.Locations[loc][res] == 0 {
			continue
		}
		if stateDiff.Locations[loc][res] > 0 {
			continue
		}
		removeActions = append(removeActions, action)
	}
	return removeActions, nil
}
