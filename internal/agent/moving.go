package agent

import (
	"Neolithic/internal/astar"
	"Neolithic/internal/world"
)

type Moving struct {
	agent     *Agent
	target    *world.Coord
	path      []*world.Coord
	pathIndex int
}

func (m *Moving) Execute(world WorldState) (WorldState, error) {
	if m.agent.behavior.curPlan == nil || m.agent.behavior.curPlan.IsComplete() {
		m.agent.behavior.curState = &Idle{
			agent:             m.agent,
			iterationsPerCall: 10,
		}
		return nil, nil
	}

	if m.target == nil {
		target := m.getTarget()
		if target == nil {
			m.agent.behavior.curState = &Performing{} // no location needed for next action
		}
		m.target = target
	}

	if m.agent.loc.IsWithin(m.target, 1) {
		m.agent.behavior.curState = &Performing{}
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
	m.agent.SetLocation(m.path[m.pathIndex])

	m.pathIndex++

	return newState, nil
}

func (m *Moving) getTarget() *world.Coord {
	nextAction := m.agent.behavior.curPlan.PeekAction()
	loc, ok := nextAction.(world.Locatable)
	if !ok { // no location needed for next action
		return nil
	}
	targetCoord := loc.Location()
	return &targetCoord
}

func (m *Moving) createPathToTarget() ([]*world.Coord, error) {
	start := m.agent.grid.Tiles[m.agent.loc.X][m.agent.loc.Y]
	end := m.agent.grid.Tiles[m.target.X][m.target.Y]

	startTile := start.(*world.Tile)
	endTile := end.(*world.Tile)

	search, err := astar.NewSearch(startTile, endTile)
	if err != nil {
		return nil, err
	}

	if err = search.RunIterations(10000); err != nil {
		return nil, err
	}

	var path []*world.Coord
	for _, node := range search.CurrentBestPath() {
		nodeTile := node.(*world.Tile)
		path = append(path, nodeTile.Coord())
	}
	return path, nil
}
