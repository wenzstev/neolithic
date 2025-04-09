package goalengine

import (
	"testing"

	"Neolithic/internal/core"

	"github.com/stretchr/testify/require"
)

func TestGoal_GetDelta(t *testing.T) {
	type testCase struct {
		goal                Goal
		numRetries          int
		expectedInInventory int
		expectNil           bool
	}

	testCases := map[string]testCase{
		"can return default chunked goal": {
			goal: Goal{
				Logic: GoalLogic{
					Chunker:      AddToLocation,
					Fallback:     FallbackChunkFunc,
					ShouldGiveUp: GiveUpIfLessThanFive,
				},
				Location: &core.Location{
					Name: "test",
				},
				Resource: &core.Resource{
					Name: "test",
				},
			},
			numRetries:          0,
			expectedInInventory: 100,
			expectNil:           false,
		},
		"can return chunked goal with one fallback": {
			goal: Goal{
				Logic: GoalLogic{
					Chunker:      AddToLocation,
					Fallback:     FallbackChunkFunc,
					ShouldGiveUp: GiveUpIfLessThanFive,
				},
				Location: &core.Location{
					Name: "test",
				},
				Resource: &core.Resource{
					Name: "test",
				},
			},
			numRetries:          1,
			expectedInInventory: 50,
			expectNil:           false,
		},
		"can return nil if should give up": {
			goal: Goal{
				Logic: GoalLogic{
					Chunker:      AddToLocation,
					Fallback:     FallbackChunkFunc,
					ShouldGiveUp: GiveUpIfLessThanFive,
				},
				Location: &core.Location{
					Name: "test",
				},
				Resource: &core.Resource{
					Name: "test",
				},
			},
			numRetries:          10,
			expectedInInventory: 0,
			expectNil:           true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			chunk := tc.goal.GetDelta(tc.numRetries)
			if tc.expectNil {
				require.Nil(t, chunk)
				return
			}
			require.NotNil(t, chunk)
			require.Equal(t, tc.expectedInInventory, chunk.Locations["test"].Inventory.GetAmount(tc.goal.Resource))
		})
	}
}

func TestAddToLocation(t *testing.T) {
	location := &core.Location{
		Name:      "test-location",
		Inventory: core.NewInventory(),
	}
	resource := &core.Resource{
		Name: "test-resource",
	}

	// Test that AddToLocation adds DefaultIncreaseAmount to location's inventory
	result := AddToLocation(location, resource)

	// Verify location in result has correct resource amount
	require.Equal(t, DefaultIncreaseAmount, result.Locations["test-location"].Inventory.GetAmount(resource))

	// Verify original location is unchanged (deep copy worked)
	require.Equal(t, 0, location.Inventory.GetAmount(resource))
}

func TestFallbackChunkFunc(t *testing.T) {
	// Create a world state with resources
	location1 := core.Location{
		Name:      "location1",
		Inventory: core.NewInventory(),
	}
	location2 := core.Location{
		Name:      "location2",
		Inventory: core.NewInventory(),
	}

	resource1 := &core.Resource{Name: "resource1"}
	resource2 := &core.Resource{Name: "resource2"}

	location1.Inventory.AdjustAmount(resource1, 100)
	location1.Inventory.AdjustAmount(resource2, 50)
	location2.Inventory.AdjustAmount(resource1, 20)

	worldState := &core.WorldState{
		Locations: map[string]core.Location{
			"location1": location1,
			"location2": location2,
		},
	}

	// Test that FallbackChunkFunc halves the amount of each resource
	result := FallbackChunkFunc(worldState)

	// Verify amounts were halved
	require.Equal(t, 50, result.Locations["location1"].Inventory.GetAmount(resource1))
	require.Equal(t, 25, result.Locations["location1"].Inventory.GetAmount(resource2))
	require.Equal(t, 10, result.Locations["location2"].Inventory.GetAmount(resource1))

	// Verify original worldState is unchanged (deep copy worked)
	require.Equal(t, 100, worldState.Locations["location1"].Inventory.GetAmount(resource1))
	require.Equal(t, 50, worldState.Locations["location1"].Inventory.GetAmount(resource2))
	require.Equal(t, 20, worldState.Locations["location2"].Inventory.GetAmount(resource1))
}

func TestGoal_GetGoalChunk(t *testing.T) {
	type testCase struct {
		startingAmount int
		expectedAmount int
		numRetries     int
		expectNil      bool
	}

	resource := &core.Resource{Name: "test-resource"}

	testCases := map[string]testCase{
		"processes state with no retries": {
			startingAmount: 50,
			expectedAmount: 50 + DefaultIncreaseAmount,
			numRetries:     0,
			expectNil:      false,
		},
		"processes state with one retry": {
			startingAmount: 50,
			expectedAmount: 50 + DefaultIncreaseAmount/2,
			numRetries:     1,
			expectNil:      false,
		},
		"returns nil with multiple retries": {
			startingAmount: 0,
			expectedAmount: 0,
			numRetries:     5,
			expectNil:      true,
		},
	}

	for name, tc := range testCases {
		t.Run(name, func(t *testing.T) {
			// Create the starting location
			location := &core.Location{
				Name:      "test-location",
				Inventory: core.NewInventory(),
			}
			location.Inventory.AdjustAmount(resource, tc.startingAmount)

			// Create world state
			worldState := &core.WorldState{
				Locations: map[string]core.Location{
					"test-location": *location.DeepCopy(),
				},
			}

			// Create the goal
			goal := Goal{
				Name:     "test-goal",
				Resource: resource,
				Location: location.DeepCopy(),
				Logic: GoalLogic{
					Chunker:      AddToLocation,
					Fallback:     FallbackChunkFunc,
					ShouldGiveUp: GiveUpIfLessThanFive,
				},
			}

			// Get the result
			result := goal.GetGoalChunk(worldState, tc.numRetries)

			if tc.expectNil {
				require.Nil(t, result, "expected nil result")
				return
			}

			require.NotNil(t, result, "expected non-nil result")
			resultLoc, exists := result.Locations["test-location"]
			require.True(t, exists, "expected test location in result")

			actualAmount := resultLoc.Inventory.GetAmount(resource)
			require.Equal(t, tc.expectedAmount, actualAmount,
				"expected resource amount %d but got %d",
				tc.expectedAmount, actualAmount)
		})
	}
}
