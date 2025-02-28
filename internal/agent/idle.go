package agent

import (
	"errors"

	"Neolithic/internal/astar"
	"Neolithic/internal/planner"
)

type Idle struct {
	iterationsPerCall int
	planner           *astar.SearchState
	agent             Agent
}

func (i *Idle) Execute(world WorldState, _ float64) (WorldState, error) {
	if i.planner == nil {
		search, err := i.createSearchState(world)
		if err != nil {
			return nil, err
		}
		i.planner = search
	}

	if err := i.planner.RunIterations(i.iterationsPerCall); err != nil {
		if errors.Is(err, astar.ErrNoPath) {
			// TODO: we need some way to communicate to the agent that the goal is unreachable. but need to implement goal first
		}
		return nil, err
	}

	if !i.planner.FoundBest {
		return nil, nil // we'll continue the search next Execute call
	}

	actionList, err := i.createActionListFromSearchState()
	if err != nil {
		return nil, err
	}

	i.agent.Behavior().curPlan = &plan{Actions: actionList}
	return nil, nil // will switch to move state when we see that curPlan is fulfilled
}

func (i *Idle) createSearchState(world WorldState) (*astar.SearchState, error) {
	behavior := i.agent.Behavior()
	worldState, ok := world.(*planner.State)
	if !ok {
		return nil, errors.New("world state is not a search state")
	}
	runInfo := &planner.GoapRunInfo{
		Agent:               i.agent,
		PossibleNextActions: behavior.PossibleActions,
	}
	start := &planner.GoapNode{
		State:       worldState,
		GoapRunInfo: runInfo,
	}

	goal := &planner.GoapNode{
		State:       behavior.Goal.(*planner.State), // temporary until state is moved out of planner
		GoapRunInfo: runInfo,
	}

	return astar.NewSearch(start, goal)
}

func (i *Idle) createActionListFromSearchState() ([]planner.Action, error) {
	if i.planner == nil {
		return nil, errors.New("no planner")
	}

	nodePlan := i.planner.CurrentBestPath()
	var actionList []planner.Action
	for _, node := range nodePlan {
		action := node.(*planner.GoapNode).Action
		actionList = append(actionList, action)
	}
	return actionList, nil
}
