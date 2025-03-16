package planner

import (
	"Neolithic/internal/core"
	"fmt"
	"math"

	"Neolithic/internal/astar"
)

// GoapNode represents a point in a GOAP process, where the planner is choosing a plan
type GoapNode struct {
	// Action is the Action taken to reach this node
	Action Action
	// State is the State of the world after running the Action
	State *core.WorldState
	// GoapRunInfo is a set of attributes that carry over throughout the goap planning process
	GoapRunInfo *GoapRunInfo
}

// GoapRunInfo represents the information that doesn't change across the GOAP planning call
type GoapRunInfo struct {
	// Agent is the agent running the planner
	Agent core.Agent
	// PossibleNextActions are all actions that the agent could take
	PossibleNextActions *[]Action
}

// Ensure GoapNode implements astar.Node
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
	return g.State.ID()
}

// Cost implements astar.Node and returns the cost of performing the acion associated with this node.
func (g *GoapNode) Cost(_ astar.Node) float64 {
	return g.Action.Cost(g.GoapRunInfo.Agent)
}

// GetSuccessors implements astar.Node and returns a list of successor astar.Node to this astar.Node.
func (g *GoapNode) GetSuccessors() ([]astar.Node, error) {
	successors := make([]astar.Node, 0)
	for _, action := range *g.GoapRunInfo.PossibleNextActions {
		newState := action.Perform(g.State, g.GoapRunInfo.Agent)
		if newState == nil {
			continue
		}
		successors = append(successors, &GoapNode{
			Action:      action,
			State:       newState,
			GoapRunInfo: g.GoapRunInfo,
		})
	}
	return successors, nil
}

// heuristic is the function used to estimate how close to the goal a given Action is. It does so by calculating the
// lowest "cost per unit" of all Action(s) that operates on a resource relevant to the goal. That value is then
// multiplied by the difference in amount of that resource between the current and the goal location.
// This heuristic is admissible because it always chooses the least "cost per unit" available, meaning it cannot
// overestimate the total cost of a given path.
func (g *GoapNode) heuristic(cur, goal *GoapNode) (float64, error) {
	var totalCost float64
	for _, goalLocation := range goal.State.Locations {
		currentLocation, ok := cur.State.Locations[goalLocation.Name]

		currentInventory, ok := cur.State.Locations[loc]
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
				stateDiff := action.GetStateChange(cur.GoapRunInfo.Agent)
				effectAmount := stateDiff.Locations[loc][item]

				costPerUnit := action.Cost(cur.GoapRunInfo.Agent) / math.Abs(float64(effectAmount))
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
		action := successor.(*GoapNode).Action

		stateDiff := action.GetStateChange(g.GoapRunInfo.Agent)
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
		action := successor.(*GoapNode).Action

		stateDiff := action.GetStateChange(g.GoapRunInfo.Agent)
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
