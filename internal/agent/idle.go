package agent

import (
	"Neolithic/internal/core"
	"errors"

	"Neolithic/internal/astar"
	"Neolithic/internal/planner"
)

// Idle is the state the agent enters in when it has no working plan. It attempts to create a plan and will proceed
// to a different state once successful.
type Idle struct {
	// iterationsPerCall is the number of iterations to run the goap planner for in a given call.
	iterationsPerCall int
	// planner is the GOAP planner that creates the agent's plan
	planner *astar.SearchState
	// agent is the agent executing the state.
	agent *Agent
}

// Execute implements State.Exeucte. Using a defined goal, it creates a plan using the GOAP planner. It runs
// the planner a given number of iterations per call. Once a plan is found, it is set on the agent and
// the agent proceeds to a Moving state.
func (i *Idle) Execute(world *core.WorldState, _ float64) (*core.WorldState, error) {
	if i.planner == nil {
		search, err := i.createSearchState(world)
		if err != nil {
			return nil, err
		}
		i.planner = search
	}

	if err := i.planner.RunIterations(i.iterationsPerCall); err != nil {
		//nolint:staticcheck
		if errors.Is(err, astar.ErrNoPath) {
			// TODO: we need some way to communicate to the agent that the goal is unreachable. but need to implement goal first
			return nil, err
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

	i.agent.Behavior.CurPlan = &plan{Actions: actionList}
	return nil, nil // will switch to move state when we see that CurPlan is fulfilled
}

// createSearchState creates the search state for the planner.
func (i *Idle) createSearchState(world *core.WorldState) (*astar.SearchState, error) {
	behavior := i.agent.Behavior
	runInfo := &planner.GoapRunInfo{
		Agent:               i.agent,
		PossibleNextActions: behavior.PossibleActions,
	}
	start := &planner.GoapNode{
		State:       world,
		GoapRunInfo: runInfo,
	}

	goal := &planner.GoapNode{
		State:       behavior.Goal, // temporary until state is moved out of planner
		GoapRunInfo: runInfo,
	}

	return astar.NewSearch(start, goal)
}

// createActionListFromSearchState creates a list of actions for the agent to follow.
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
