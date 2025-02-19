package planner

import (
	"Neolithic/internal/agent"
	"Neolithic/internal/astar"
)

//// agent.Agent represents an entity that can do things in the world (i.e., a villager)
//type Agent struct {
//	Name    string
//	Actions []Action
//}

// Resource represents a resource in the world
type Resource struct {
	Name string
}

// Location represents a location in the world
type Location struct {
	Name string
	X, Y int
}

func GetPlan(agent *agent.Agent, goal *State, world *State) ([]Action, error) {
	runInfo := &GoapRunInfo{
		Agent:               agent,
		PossibleNextActions: &agent.Actions,
	}

	start := &GoapNode{
		State:       world,
		GoapRunInfo: runInfo,
	}

	end := &GoapNode{
		State:       goal,
		GoapRunInfo: runInfo,
	}

	search, err := astar.NewSearch(start, end)
	if err != nil {
		return nil, err
	}
	if err = search.RunIterations(10000); err != nil { // todo assume better checks
		return nil, err
	}

	nodePlan := search.CurrentBestPath()
	actionList := make([]Action, 0)
	for _, node := range nodePlan {
		action := node.(*GoapNode).Action
		actionList = append(actionList, action)
	}

	return actionList, nil
}
