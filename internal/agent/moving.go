package agent

import (
	"errors"

	"Neolithic/internal/astar"
	"Neolithic/internal/core"
)

type Moving struct {
	agent  *Agent
	target *core.Coord
	path   Path
}

const (
	maxPathfindingIterations = 10000
	targetProximityThreshold = 1
)

var ErrNoPathFound = errors.New("no path found")

var _ State = (*Moving)(nil)

func (m *Moving) Execute(world *core.WorldState, _ float64) (*core.WorldState, error) {
	behavior := m.agent.Behavior

	if behavior.CurPlan == nil || behavior.CurPlan.IsComplete() {
		behavior.curState = &Idle{
			agent:             m.agent,
			iterationsPerCall: 10,
		}
		return nil, nil
	}

	if m.target == nil {
		target := m.getTarget()
		if target == nil {
			behavior.curState = &Performing{agent: m.agent} // no location needed for next action
			return nil, nil
		}
		m.target = target
	}

	if m.agent.Position.IsWithin(*m.target, targetProximityThreshold) {
		behavior.curState = &Performing{agent: m.agent}
		return nil, nil
	}

	if m.path == nil {
		path, err := m.createPathToTarget(world)
		if err != nil {
			return nil, err
		}
		m.path = path
	}

	newState := world.DeepCopy()
	newAgent := newState.Agents[m.agent.Name()]
	if m.path.IsComplete() {
		behavior.curState = &Idle{agent: m.agent}
		return nil, nil
	}
	newAgent.(*Agent).Position = m.path.NextCoord()

	return newState, nil
}

func (m *Moving) getTarget() *core.Coord {
	nextAction := m.agent.Behavior.CurPlan.PeekAction()
	loc, ok := nextAction.(core.Locatable)
	if !ok { // no location needed for next action
		return nil
	}
	targetCoord := loc.Location().Coord
	return &targetCoord
}

func (m *Moving) createPathToTarget(world *core.WorldState) (Path, error) {
	start := world.Grid.CellAt(m.agent.Position)
	end := world.Grid.CellAt(*m.target)

	search, err := astar.NewSearch(start, end)
	if err != nil {
		return nil, err
	}

	if err = search.RunIterations(maxPathfindingIterations); err != nil {
		return nil, err
	}

	if !search.FoundBest {
		return nil, ErrNoPathFound
	}

	var coords []core.Coord
	for _, node := range search.CurrentBestPath() {
		nodeCoord := node.(core.Cell).Coord()
		coords = append(coords, nodeCoord)
	}
	return NewCoordPath(coords), nil
}
