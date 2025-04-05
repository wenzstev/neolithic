package agent

import (
	"testing"

	"Neolithic/internal/core"
	"Neolithic/internal/logging"
	"github.com/stretchr/testify/require"
)

type mockPath struct {
	nextCoord  core.Coord
	isComplete bool
}

func (p *mockPath) NextCoord() core.Coord {
	return p.nextCoord
}

func (p *mockPath) IsComplete() bool {
	return p.isComplete
}

func TestMoving_Execute(t *testing.T) {
	type testCase struct {
		agentPosition   core.Coord
		hasPath         bool
		nextCoord       core.Coord
		isComplete      bool
		plan            Plan
		target          *core.Coord
		expectedErr     error
		newAgentPositon core.Coord
		expectedState   State
		expectedPath    *CoordPath
	}

	tests := map[string]testCase{
		"transitions to idle when plan is nil": {
			agentPosition: core.Coord{X: 0, Y: 0},
			plan:          nil,
			expectedState: &Idle{},
		},
		"transitions to idle when plan is complete": {
			agentPosition: core.Coord{X: 0, Y: 0},
			plan: &MockPlan{
				Complete: true,
			},
			expectedState: &Idle{},
		},
		"transitions to performing when no Target needed": {
			agentPosition: core.Coord{X: 0, Y: 0},
			plan: &MockPlan{
				NextAction: &mockAction{},
			},
			expectedState: &Performing{},
		},
		"transitions to performing when at Target": {
			agentPosition: core.Coord{X: 1, Y: 1},
			target:        &core.Coord{X: 1, Y: 1},
			plan: &MockPlan{
				NextAction: &mockLocationAction{
					location: &core.Location{
						Coord: core.Coord{X: 1, Y: 1},
					},
				},
			},
			expectedState: &Performing{},
		},
		"creates Path when no Path exists": {
			agentPosition: core.Coord{X: 0, Y: 0},
			hasPath:       false,
			target:        &core.Coord{X: 5, Y: 5},
			plan: &MockPlan{
				NextAction: &mockLocationAction{
					location: &core.Location{
						Coord: core.Coord{X: 1, Y: 1},
					},
				},
			},
			expectedState: &Moving{},
			expectedPath: &CoordPath{
				coords: []core.Coord{{X: 0, Y: 0}, {X: 1, Y: 1}, {X: 2, Y: 2}, {X: 3, Y: 3}, {X: 4, Y: 4}, {X: 5, Y: 5}},
				index:  2,
			},
		},
		"updates agent position to next Path coordinate": {
			agentPosition: core.Coord{X: 0, Y: 0},
			hasPath:       true,
			nextCoord:     core.Coord{X: 1, Y: 1},
			isComplete:    false,
			plan: &MockPlan{
				NextAction: &mockLocationAction{
					location: &core.Location{
						Coord: core.Coord{X: 2, Y: 2},
					},
				},
			},
			newAgentPositon: core.Coord{X: 1, Y: 1},
		},
		"maintains Path when not complete": {
			agentPosition: core.Coord{X: 0, Y: 0},
			hasPath:       true,
			nextCoord:     core.Coord{X: 1, Y: 1},
			isComplete:    false,
			plan: &MockPlan{
				NextAction: &mockLocationAction{
					location: &core.Location{
						Coord: core.Coord{X: 2, Y: 2},
					},
				},
			},
			expectedState: &Moving{},
		},
		"transitions to idle when Path is complete": {
			agentPosition: core.Coord{X: 0, Y: 0},
			hasPath:       true,
			isComplete:    true,
			expectedState: &Idle{},
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			testAgent := &Agent{
				name: "testAgent",
				Behavior: &Behavior{
					CurPlan:  tc.plan,
					CurState: &Moving{},
				},
				Position: tc.agentPosition,
			}
			var path Path
			if tc.hasPath {
				path = &mockPath{
					nextCoord:  tc.nextCoord,
					isComplete: tc.isComplete,
				}
			}

			testMoving := &Moving{
				agent:  testAgent,
				Path:   path,
				Target: tc.target,
				logger: logging.NewLogger("info"),
			}

			startWorld := &core.WorldState{
				Grid: &mockGrid{},
				Agents: map[string]core.Agent{
					"testAgent": testAgent,
				},
			}

			newState, err := testMoving.Execute(startWorld, 0)
			if expectedErr := tc.expectedErr; expectedErr != nil {
				require.ErrorIs(t, err, expectedErr)
			} else {
				require.NoError(t, err)
			}
			if tc.newAgentPositon != (core.Coord{}) {
				require.Equal(t, tc.newAgentPositon, newState.Agents["testAgent"].(*Agent).Position)
			}
			if tc.expectedState != nil {
				require.IsType(t, tc.expectedState, testAgent.Behavior.CurState)
			}
			if tc.expectedPath != nil {
				require.Equal(t, tc.expectedPath, testMoving.Path.(*CoordPath))
			}
		})
	}
}
