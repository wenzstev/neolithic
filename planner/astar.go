package planner

import (
	"container/heap"
	"errors"
	"math"
)

// ErrNoPath is thrown when the AStar Planner is unable to find a path to the goal state
var ErrNoPath = errors.New("no path found to goal state")

// ErrNoAgent is thrown when the provided agent is not in the provided state
var ErrNoAgent = errors.New("agent not in state")

// Actions represents a slice of Action. It is used to represent all possible actions that an agent can take.
// The AStar path is calculated off of this list of actions
type Actions []Action

// AStar is the function used to calculate the optimal (or best found) from a given start State to a goal State, using
// the provided Agent. It will run maxDistance times before returning.
func (a *Actions) AStar(startState, goalState *State, agent *Agent, maxDistance int) (*AStarOutput, error) {
	if startState.Agents[agent] == nil {
		return nil, ErrNoAgent
	}

	openSet := &PriorityQueue{}
	heap.Init(openSet)

	startNode := &Node{
		state:  startState,
		gCost:  0,
		hCost:  heuristic(startState, goalState, agent),
		parent: nil,
		action: nil,
	}
	heap.Push(openSet, startNode)

	closedSet := make(map[string]bool)

	for openSet.Len() > 0 {
		currentNode := heap.Pop(openSet).(*Node)
		if isGoal(currentNode.state, goalState, agent) {
			return &AStarOutput{
				actions:       reconstructPath(currentNode),
				totalCost:     currentNode.gCost,
				expectedState: currentNode.state,
			}, nil
		}

		id, err := currentNode.state.ID()
		if err != nil {
			return nil, err
		}
		closedSet[id] = true

		for _, successor := range a.getSuccessors(currentNode.state, agent) {
			newState := successor.newState
			newStateId, err := newState.ID()
			if err != nil {
				return nil, err
			}
			if closedSet[newStateId] {
				continue
			}

			stepCost := successor.seq.Cost(agent)
			newGCost := currentNode.gCost + stepCost
			newHCost := heuristic(newState, goalState, agent)

			newNode := &Node{
				state:  newState,
				gCost:  newGCost,
				hCost:  newHCost,
				parent: currentNode,
				action: successor.seq,
			}

			heap.Push(openSet, newNode)

		}
		if maxDistance == 0 {
			return &AStarOutput{
				actions:       reconstructPath(currentNode),
				totalCost:     currentNode.gCost,
				expectedState: currentNode.state,
			}, nil
		}
		maxDistance--
	}

	return nil, ErrNoPath
}

// SuccessorState represents a world state that exists after a given action (or action sequence) has been called.
type SuccessorState struct {
	seq      Action
	newState *State
}

// getSuccessors returns all actions and the results of those actions.
func (a *Actions) getSuccessors(currentState *State, agent *Agent) []SuccessorState {
	successors := make([]SuccessorState, 0)
	for _, action := range *a {
		newState := action.Perform(currentState, agent)
		if newState == nil {
			continue
		}
		successors = append(successors, SuccessorState{
			seq:      action,
			newState: newState,
		})
	}
	return successors
}

// heuristic calculates a distance metric between the current and goal states
// the lower the return value, the closer current is to the goal
func heuristic(current, goal *State, agent *Agent) float64 {
	var cost float64
	for loc, goalInventory := range goal.Locations {
		currentInventory, ok := current.Locations[loc]
		for item, goalAmount := range goalInventory {
			currentAmount := 0
			if ok {
				currentAmount = currentInventory[item]
			}
			agentInventory := current.Agents[agent]
			amountInAgentInventory := agentInventory[item] // amount will be zero if agent has none

			cost += math.Abs(math.Abs(float64(currentAmount-goalAmount)) - float64(amountInAgentInventory)*0.90)
		}
	}
	return cost
}

// isGoal checks if we are at the goal state
func isGoal(current, goal *State, agent *Agent) bool {
	return heuristic(current, goal, agent) == 0
}

// reconstructPath reconstructs the path taken to get to the goal, traversing backwards through the nodes and
// returning a slice of actions that were taken to get to the node.
func reconstructPath(n *Node) []Action {
	if n == nil {
		return nil
	}

	var actions []Action
	for current := n; current != nil && current.parent != nil; current = current.parent {
		actions = append(actions, current.action)
	}
	return actions
}
