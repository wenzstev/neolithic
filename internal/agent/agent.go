package agent

import (
	"Neolithic/internal/astar"
	"Neolithic/internal/grid"
	"Neolithic/internal/planner"
	world2 "Neolithic/internal/world"
	"errors"
)

type Plan struct {
	Actions     []planner.Action
	curLocation int
}

func (p *Plan) isComplete() bool {
	return p.curLocation >= len(p.Actions)
}

func (p *Plan) NextAction() planner.Action {
	return p.Actions[p.curLocation]
}

type Agent struct {
	Name     string
	Actions  *[]planner.Action
	curPlan  *Plan
	Goal     *planner.State
	curState State
	loc      world2.Coord
	grid     *grid.Grid
}

type State interface {
	Execute(world *planner.State) (*planner.State, error)
}

type Idle struct {
	iterationsPerCall int
	planner           *astar.SearchState
	agent             *Agent
}

func (i *Idle) Execute(world *planner.State) (*planner.State, error) {
	if i.agent.curPlan != nil {
		i.agent.curState = &Moving{}
		return nil, nil
	}

	if i.planner == nil {
		search, err := i.createSearchState(world)
		if err != nil {
			return nil, err
		}
		i.planner = search
	}

	if err := i.planner.RunIterations(i.iterationsPerCall); err != nil {
		return nil, err
	}

	if !i.planner.FoundBest {
		return nil, nil // we'll continue the search next Execute call
	}

	actionList, err := i.createActionListFromSearchState()
	if err != nil {
		return nil, err
	}

	i.agent.curPlan = &Plan{Actions: actionList}
	return nil, nil // will switch to move state when we see that curPlan is fulfilled
}

func (i *Idle) createSearchState(world *planner.State) (*astar.SearchState, error) {
	runInfo := &planner.GoapRunInfo{
		Agent:               i.agent,
		PossibleNextActions: i.agent.Actions,
	}
	start := &planner.GoapNode{
		State:       world,
		GoapRunInfo: runInfo,
	}

	goal := &planner.GoapNode{
		State:       i.agent.Goal,
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

type Moving struct {
	agent     *Agent
	target    *world2.Coord
	path      []*world2.Tile
	pathIndex int
}

func (m *Moving) Execute(world *planner.State) (*planner.State, error) {
	if m.agent.curPlan == nil || m.agent.curPlan.isComplete() {
		m.agent.curState = &Idle{
			agent:             m.agent,
			iterationsPerCall: 10,
		}
		return nil, nil
	}

	if m.target == nil {
		target := m.getTarget()
		if target == nil {
			m.agent.curState = &Performing{} // no location needed for next action
		}
		m.target = target
	}

	if m.agent.loc.IsWithin(m.target, 1) {
		m.agent.curState = &Performing{}
		return nil, nil
	}

	if m.path == nil {
		path, err := m.createPathToTarget()
		if err != nil {
			return nil, err
		}
		m.path = path
	}

	newState := world.Copy()
	newState.Agents[m.agent].Coord = m.path[m.pathIndex]

	m.pathIndex++

	return newState, nil
}

func (m *Moving) getTarget() *world2.Coord {
	nextAction := m.agent.curPlan.NextAction()
	loc, ok := nextAction.(planner.Locatable)
	if !ok { // no location needed for next action
		return nil
	}
	targetCoord := loc.Location()
	return &targetCoord
}

func (m *Moving) createPathToTarget() ([]*world2.Tile, error) {
	start := m.agent.grid.Tiles[m.agent.loc.X][m.agent.loc.Y]
	end := m.agent.grid.Tiles[m.target.X][m.target.Y]

	startTile := start.(*world2.Tile)
	endTile := end.(*world2.Tile)

	search, err := astar.NewSearch(startTile, endTile)
	if err != nil {
		return nil, err
	}

	if err = search.RunIterations(10000); err != nil {
		return nil, err
	}

	var path []*world2.Tile
	for _, node := range search.CurrentBestPath() {
		nodeTile := node.(*world2.Tile)
		path = append(path, nodeTile)
	}
	return path, nil
}

type Performing struct {
	action   planner.Action
	timeLeft float64
	agent    *Agent
}

func (a *Performing) Execute(world *planner.State) (*planner.State, error) {
	a.action = a.agent.curPlan.NextAction()
	actionDuration, ok := a.action.(planner.RequiresTime)
	if ok {
		a.timeLeft = actionDuration.TimeNeeded()
	}

	if a.timeLeft > 0 {
		a.timeLeft -= 1.0 / 60.0 // called every tick, update is called 60 times a second
		return nil, nil
	}

	newState := a.action.Perform(world, agent)
	if newState == nil { // action failed
		a.agent.curState = &Idle{}
		return nil, nil
	}
	a.agent.curPlan.curLocation++
	if a.agent.curPlan.isComplete() {
		a.agent.curState = &Idle{}
	} else {
		a.agent.curState = &Moving{}
	}
	
	return newState, nil
}
