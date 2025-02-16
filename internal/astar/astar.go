package astar

import (
	"Neolithic/internal/logging"
	"container/heap"
	"errors"
	"fmt"
	"log/slog"
	"math"
)

// ErrNoPath is thrown when the Run Planner is unable to find a path to the goal state
var ErrNoPath = errors.New("no path found to goal state")

// SearchState represents a single AStar search
type SearchState struct {
	// Start is where the AStar search will begin
	Start Node
	// Goal is where the search is trying to go
	Goal Node
	// BestCost is the current best cost found
	BestCost float64
	// Iterations is the number of times the AStar algorithm has run through its main looop
	Iterations int
	// FoundBest indicates if the algorithm has found the optimal path to the end.
	FoundBest bool
	// openSet is a heap used to store all open nodes
	openSet *PriorityQueue
	// openSetMap is a map also used to store open nodes
	openSetMap map[string]*searchNode
	// closedSet is a map used to store all visited nodes
	closedSet map[string]bool
	// bestSolution is the head node of the current best solution. Not meant to be accessed directly,
	// instead use CurrentBestPath
	bestSolution *searchNode
	// logger is used to log information about the search
	logger *slog.Logger
}

// Node represents an intermediate state in the algorithm. It's expected to be different for each implementation
type Node interface {
	// Heuristic is the function that estimates the distance from the node to the goal node.
	Heuristic(goal Node) (float64, error)
	// ID returns a unique string identifier for the node.
	ID() (string, error)
	// Cost returns the cost to get from prev Node to this Node.
	Cost(prev Node) float64
	// GetSuccessors returns a slice of all Node this node is connected to.
	GetSuccessors() ([]Node, error)
}

// searchNode is the internal node struct used by Run to track its progress.
type searchNode struct {
	// gCost is the total cost from the start to this nodeState
	gCost float64
	// hCost is the heuristic cost from the nodeState to the goal
	hCost float64
	// parent is the nodeState that led to this nodeState
	parent *searchNode
	// nodeState is the current state of the world at this point in the algorithm
	nodeState Node
	// index is the index of the node in the PriorityQueue
	index int
}

// fCost calculates the sum of the gCost and the hCost for the searchNode.
func (n *searchNode) fCost() float64 {
	return n.gCost + n.hCost
}

// NewSearch initializes a new SearchState with a start and finish Node
func NewSearch(start, goal Node) (*SearchState, error) {
	search := &SearchState{
		Start:  start,
		Goal:   goal,
		logger: logging.NewLogger("info"),
	}
	if err := search.init(start, goal); err != nil {
		return nil, err
	}

	return search, nil
}

// init initializes a new Search State
func (s *SearchState) init(start, goal Node) error {
	s.openSet = &PriorityQueue{}
	heap.Init(s.openSet)

	s.closedSet = make(map[string]bool)
	s.openSetMap = make(map[string]*searchNode)

	s.BestCost = math.Inf(1)

	hCost, err := start.Heuristic(goal)
	if err != nil {
		return err
	}

	startNode := &searchNode{
		gCost:     0,
		hCost:     hCost,
		parent:    nil,
		nodeState: start,
	}

	heap.Push(s.openSet, startNode)
	startID, err := start.ID()
	if err != nil {
		return err
	}

	s.openSetMap[startID] = startNode

	return nil
}

// RunIterations runs the SearchState using the A* algorithm for the given number of iterations,
// or until an optimal path is found.
func (s *SearchState) RunIterations(numIterations int) error {
	curIterations := 0
	for s.openSet.Len() > 0 && curIterations < numIterations {
		curIterations++
		s.Iterations++

		currentNode := heap.Pop(s.openSet).(*searchNode)

		currentID, err := currentNode.nodeState.ID()
		if err != nil {
			return err
		}

		delete(s.openSetMap, currentID)

		if currentNode.fCost() > s.BestCost {
			continue // node path not better than what we already have
		}

		isGoal, err := isGoal(currentNode.nodeState, s.Goal)
		if err != nil {
			return err
		}

		if isGoal {
			s.logger.Debug(fmt.Sprintf("Found solution: %s", currentID))
			if currentNode.gCost < s.BestCost {
				s.BestCost = currentNode.gCost
				s.bestSolution = currentNode
				s.logger.Debug(fmt.Sprintf("New solution is %s", s.CurrentBestPath()))
				continue // don't need to check successors of goal state
			}
		}

		successors, err := currentNode.nodeState.GetSuccessors()
		if err != nil {
			return err
		}

		for _, successor := range successors {
			sucId, err := successor.ID()
			if err != nil {
				return err
			}

			s.logger.Debug(fmt.Sprintf("Checking successor %s", sucId))

			stepCost := successor.Cost(currentNode.nodeState)
			newGCost := currentNode.gCost + stepCost
			newHCost, err := successor.Heuristic(s.Goal)
			if err != nil {
				return err
			}

			s.logger.Debug(fmt.Sprintf("stepCost: %v, gCost: %v, hCost: %v", stepCost, newGCost, newHCost))

			newFCost := newGCost + newHCost
			if newFCost >= s.BestCost { // this is fine because we cannot have negative h values
				continue
			}
			if s.closedSet[sucId] {
				continue // already looked at this node
			}

			if existing, ok := s.openSetMap[sucId]; ok {
				if newGCost < existing.gCost {
					existing.gCost = newGCost
					existing.hCost = newHCost
					existing.parent = currentNode
					heap.Fix(s.openSet, existing.index)
				}
			} else {
				newNode := &searchNode{
					gCost:     newGCost,
					hCost:     newHCost,
					parent:    currentNode,
					nodeState: successor,
				}
				heap.Push(s.openSet, newNode)
				s.openSetMap[sucId] = newNode
			}
		}
		s.closedSet[currentID] = true
	}
	if len(s.openSetMap) == 0 {
		if s.bestSolution == nil {
			return ErrNoPath // no path was found
		}
		s.FoundBest = true
	}

	return nil
}

// CurrentBestPath returns an array of nodes as the current best path to the goal.
func (s *SearchState) CurrentBestPath() []Node {
	return reconstructPath(s.bestSolution)
}

// isGoal checks if the current node is the goal node or not.
func isGoal(currentNode Node, goal Node) (bool, error) {
	curId, err := currentNode.ID()
	if err != nil {
		return false, err
	}

	goalId, err := goal.ID()
	if err != nil {
		return false, err
	}

	// try to compare ID first
	if curId == goalId {
		return true, nil
	}

	// because heuristic might ignore details, compare heuristic as backup
	hVal, err := currentNode.Heuristic(goal)
	if err != nil {
		return false, err
	}
	return hVal == 0, nil
}

// reconstructPath reconstructs the given path, returning an array of Node that lists the steps
// needed to reach the goal.
func reconstructPath(n *searchNode) []Node {
	if n == nil {
		return nil
	}

	var nodes []Node
	for current := n; current != nil; current = current.parent {
		nodes = append(nodes, current.nodeState)
	}

	// reverse the slice so the first move is first
	for i, j := 0, len(nodes)-1; i < j; i, j = i+1, j-1 {
		nodes[i], nodes[j] = nodes[j], nodes[i]
	}

	return nodes
}
