package agent

import (
	"Neolithic/internal/core"
	"errors"
	"log/slog"

	"Neolithic/internal/astar"
	"Neolithic/internal/planner"
)

const defaultNumIterations = 100

// Idle is the state the Agent enters in when it has no working plan. It attempts to create a plan and will proceed
// to a different state once successful.
type Idle struct {
	// IterationsPerCall is the number of iterations to run the goap planner for in a given call.
	IterationsPerCall int
	// planner is the GOAP planner that creates the agent's plan
	planner *astar.SearchState
	// agent is the agent executing the state.
	agent *Agent
	// logger is used for logging state events
	logger *slog.Logger
}

// Execute implements State.Exeucte. Using a defined goal, it creates a plan using the GOAP planner. It runs
// the planner a given number of iterations per call. Once a plan is found, it is set on the Agent and
// the Agent proceeds to a Moving state.
func (i *Idle) Execute(world *core.WorldState, _ float64) (*core.WorldState, error) {
	i.logger.Debug("idle state execute", "agent", i.agent.Name())

	if i.IterationsPerCall == 0 {
		i.IterationsPerCall = defaultNumIterations
	}

	if i.agent.Behavior.Goal == nil {
		// no goal, do nothing
		i.logger.Debug("no goal set, staying idle", "agent", i.agent.Name())
		return nil, nil
	}

	if i.planner == nil {
		i.logger.Info("creating new search state", "agent", i.agent.Name())
		search, err := i.createSearchState(world)
		if err != nil {
			i.logger.Error("failed to create search state", "agent", i.agent.Name(), "error", err)
			return nil, err
		}
		i.planner = search
	}

	if err := i.planner.RunIterations(i.IterationsPerCall); err != nil {
		//nolint:staticcheck
		if errors.Is(err, astar.ErrNoPath) {
			// TODO: we need some way to communicate to the agent that the goal is unreachable. but need to implement goal first
			i.logger.Error("no Path found to goal", "agent", i.agent.Name(), "error", err)
			return nil, err
		}
		i.logger.Error("planner iteration error", "agent", i.agent.Name(), "error", err)
		return nil, err
	}

	if !i.planner.FoundBest {
		i.logger.Debug("plan not found yet, continuing search", "agent", i.agent.Name())
		return nil, nil // we'll continue the search next Execute call
	}

	i.logger.Info("plan found, creating action list", "agent", i.agent.Name())
	actionList, err := i.createActionListFromSearchState()
	if err != nil {
		i.logger.Error("failed to create action list", "agent", i.agent.Name(), "error", err)
		return nil, err
	}

	i.agent.Behavior.CurPlan = &plan{Actions: actionList, curLocation: 1} // 1 because index of zero is null
	i.agent.Behavior.CurState = &Moving{agent: i.agent, logger: i.logger}
	i.logger.Info("transitioning to moving state", "agent", i.agent.Name(), "planLength", len(actionList))
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

	return astar.NewSearch(start, goal, i.logger)
}

// createActionListFromSearchState creates a list of actions for the Agent to follow.
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

// NewIdle creates a new Idle state
func NewIdle(agent *Agent, logger *slog.Logger) *Idle {
	return &Idle{
		IterationsPerCall: defaultNumIterations,
		agent:             agent,
		logger:            logger,
	}
}
