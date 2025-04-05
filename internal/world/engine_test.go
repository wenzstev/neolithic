package world

import (
	"errors"
	"log/slog"
	"testing"

	"Neolithic/internal/agent"
	"Neolithic/internal/core"
	"Neolithic/internal/grid"
	"Neolithic/internal/logging"
	"github.com/stretchr/testify/assert"
)

func TestNewEngine(t *testing.T) {
	type testCase struct {
		width             int
		height            int
		expectEngineError bool
		expectedState     *core.WorldState
	}

	tests := map[string]testCase{
		"valid dimensions": {
			width:  10,
			height: 10,
			expectedState: &core.WorldState{
				Locations: map[string]core.Location{},
				Agents:    map[string]core.Agent{},
			},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {

			worldGrid, err := grid.New(tc.width, tc.height, cellSize)
			assert.NoError(t, err)
			assert.NoError(t, worldGrid.Initialize(testMakeTile))

			logger := slog.New(slog.NewTextHandler(nil, nil))
			engine, err := NewEngine(worldGrid, logger)

			if tc.expectEngineError {
				assert.Error(t, err)
				assert.Nil(t, engine)
				return
			}

			assert.NoError(t, err)
			assert.NotNil(t, engine)
			assert.NotNil(t, engine.World)
			assert.NotNil(t, engine.villagerImage)
			assert.NotNil(t, engine.locationImage)
			assert.NotNil(t, engine.logger)
			assert.Equal(t, tc.expectedState.Locations, engine.World.Locations)
			assert.Equal(t, tc.expectedState.Agents, engine.World.Agents)
		})
	}
}

func TestEngine_Tick(t *testing.T) {
	type testCase struct {
		agents         map[string]core.Agent
		expectedError  error
		expectedAgents map[string]core.Agent
		expectedGrid   *grid.Grid
	}

	logger := logging.NewLogger("info")

	// Create test grid
	testGrid, err := grid.New(10, 10, cellSize)
	assert.NoError(t, err)
	assert.NoError(t, testGrid.Initialize(testMakeTile))

	// Create test agents
	normalAgent := agent.NewAgent("normal", logger)
	normalAgent.Behavior.CurState = agent.NewIdle(normalAgent, logger)

	errorAgent := agent.NewAgent("error", logger)
	moveState := agent.NewMoving(errorAgent, logger)
	moveState.Target = &core.Coord{X: -5, Y: -5}
	errorAgent.Behavior.CurState = moveState
	errorAgent.Behavior.CurPlan = &agent.MockPlan{Complete: false}

	stateChangeAgent := agent.NewAgent("stateChange", logger)
	stateChangeAgent.Behavior.CurState = agent.NewPerforming(stateChangeAgent, logger)
	stateChangeAgent.Behavior.CurPlan = &agent.MockPlan{
		Complete:   false,
		NextAction: &mockAction{testVal: "test"},
	}

	stateChangeAfterAgent := stateChangeAgent.DeepCopy()
	stateChangeAfterAgent.(*agent.Agent).Behavior.CurState = &agent.Moving{}

	tests := map[string]testCase{
		"no agents": {
			agents:         map[string]core.Agent{},
			expectedAgents: map[string]core.Agent{},
			expectedGrid:   testGrid,
		},
		"single agent no state change": {
			agents: map[string]core.Agent{
				normalAgent.Name(): normalAgent,
			},
			expectedAgents: map[string]core.Agent{
				normalAgent.Name(): normalAgent,
			},
			expectedGrid: testGrid,
		},
		"single agent with state change": {
			agents: map[string]core.Agent{
				stateChangeAgent.Name(): stateChangeAgent,
			},
			expectedAgents: map[string]core.Agent{
				stateChangeAgent.Name(): stateChangeAfterAgent,
			},
			expectedGrid: testGrid,
		},
		"multiple agents": {
			agents: map[string]core.Agent{
				normalAgent.Name():      normalAgent,
				stateChangeAgent.Name(): stateChangeAgent,
			},
			expectedAgents: map[string]core.Agent{
				normalAgent.Name():      normalAgent,
				stateChangeAgent.Name(): stateChangeAfterAgent,
			},
			expectedGrid: testGrid,
		},
		"agent returns error": {
			agents: map[string]core.Agent{
				errorAgent.Name(): errorAgent,
			},
			expectedError: errors.New("heuristic called on non-Tile"),
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			engine, err := NewEngine(testGrid, logger)
			assert.NoError(t, err)

			engine.World.Agents = tc.agents

			err = engine.Tick(1.0 / 60.0)

			if tc.expectedError != nil {
				assert.ErrorContains(t, err, tc.expectedError.Error())
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tc.expectedGrid, engine.World.Grid)

			// Compare agents
			assert.Equal(t, len(tc.expectedAgents), len(engine.World.Agents))
			for name, expectedAgent := range tc.expectedAgents {
				actualAgent, exists := engine.World.Agents[name]
				assert.True(t, exists)
				expectedAgentStruct := expectedAgent.(*agent.Agent)
				actualAgentStruct := actualAgent.(*agent.Agent)
				assert.IsType(t, expectedAgentStruct.Behavior.CurState, actualAgentStruct.Behavior.CurState)
			}
		})
	}
}
