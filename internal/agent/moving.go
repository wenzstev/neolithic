package agent

import (
	"errors"
	"log/slog"

	"Neolithic/internal/astar"
	"Neolithic/internal/core"
)

// Moving represents the state of an agent as it navigates along a Path toward a Target location.
type Moving struct {
	// agent is the agent moving
	agent *Agent
	// Target is where the agent is moving to
	Target *core.Coord
	// Path is the sequence of coordinates to get to Target
	Path Path
	// logger is the logger
	logger *slog.Logger
}

const (
	maxPathfindingIterations = 10000
	targetProximityThreshold = 1
)

var ErrNoPathFound = errors.New("no Path found")

var _ State = (*Moving)(nil)

// Execute progresses the Moving state, handling path creation and movement, and updates the agent's state as needed.
func (m *Moving) Execute(world *core.WorldState, _ float64) (*core.WorldState, error) {
	m.logger.Debug("moving state execute", "agent", m.agent.Name())

	behavior := m.agent.Behavior

	if behavior.CurPlan == nil || behavior.CurPlan.IsComplete() {
		m.logger.Info("plan complete or nil, transitioning to idle", "agent", m.agent.Name())
		behavior.CurState = &Idle{agent: m.agent, logger: m.logger}
		return nil, nil
	}

	if m.Target == nil {
		target := m.getTarget()
		if target == nil {
			m.logger.Info("no Target location needed, transitioning to performing", "agent", m.agent.Name())
			behavior.CurState = &Performing{agent: m.agent, logger: m.logger} // no location needed for next action
			return nil, nil
		}
		m.Target = target
		m.logger.Debug("Target set", "agent", m.agent.Name(), "Target", m.Target)
	}

	if m.agent.Position.IsWithin(*m.Target, targetProximityThreshold) {
		m.logger.Info("reached Target, transitioning to performing", "agent", m.agent.Name(), "position", m.agent.Position, "Target", m.Target)
		behavior.CurState = &Performing{agent: m.agent, logger: m.logger}
		return nil, nil
	}

	if m.Path == nil {
		m.logger.Debug("creating Path to Target", "agent", m.agent.Name(), "start", m.agent.Position, "Target", m.Target)
		path, err := m.createPathToTarget(world)
		if err != nil {
			m.logger.Error("failed to create Path", "agent", m.agent.Name(), "error", err)
			return nil, err
		}
		m.Path = path
	}

	newState := world.DeepCopy()
	newAgent := newState.Agents[m.agent.Name()]
	if m.Path.IsComplete() {
		m.logger.Info("Path complete, transitioning to performing", "agent", m.agent.Name())
		newAgent.(*Agent).Behavior.CurState = &Performing{agent: m.agent, logger: m.logger}
		return newState, nil
	}

	nextCoord := m.Path.NextCoord()
	m.logger.Debug("moving to next coordinate", "agent", m.agent.Name(), "from", newAgent.(*Agent).Position, "to", nextCoord)
	newAgent.(*Agent).Position = nextCoord

	return newState, nil
}

// getTarget determines the target coordinate for the agent's next action and returns it, or nil if no location is needed.
func (m *Moving) getTarget() *core.Coord {
	nextAction := m.agent.Behavior.CurPlan.PeekAction()
	loc, ok := nextAction.(core.Locatable)
	if !ok { // no location needed for next action
		return nil
	}
	targetCoord := loc.Location().Coord
	return &targetCoord
}

// createPathToTarget generates a path from the agent's current position to the target using the A* algorithm.
// Returns the computed path or an error if no valid path is found or an issue occurs during pathfinding.
func (m *Moving) createPathToTarget(world *core.WorldState) (Path, error) {
	start := world.Grid.CellAt(m.agent.Position)
	end := world.Grid.CellAt(*m.Target)

	search, err := astar.NewSearch(start, end, m.logger)
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

// NewMoving creates a new Moving state
func NewMoving(agent *Agent, logger *slog.Logger) *Moving {
	return &Moving{
		agent:  agent,
		logger: logger,
	}
}
