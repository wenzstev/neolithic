package agent

import (
	"Neolithic/internal/core"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestPlan_IsComplete tests the IsComplete method of plan.
func TestPlan_IsComplete(t *testing.T) {
	type testCase struct {
		actions     []core.Action
		curLocation int
		expected    bool
	}

	tests := map[string]testCase{
		"empty plan is complete": {
			actions:     []core.Action{},
			curLocation: 0,
			expected:    true,
		},
		"non-empty plan not complete": {
			actions:     []core.Action{&mockAction{}},
			curLocation: 0,
			expected:    false,
		},
		"plan complete after consuming all actions": {
			actions:     []core.Action{&mockAction{}},
			curLocation: 1,
			expected:    true,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &plan{
				Actions:     tc.actions,
				curLocation: tc.curLocation,
			}
			assert.Equal(t, tc.expected, p.IsComplete())
		})
	}
}

// TestPlan_PeekAction tests that PeekAction returns the next action without advancing curLocation.
func TestPlan_PeekAction(t *testing.T) {
	type testCase struct {
		actions     []core.Action
		curLocation int
		expected    core.Action
	}

	tests := map[string]testCase{
		"peek first action": {
			actions:     []core.Action{&mockAction{}, &mockNullAction{}},
			curLocation: 0,
			expected:    &mockAction{},
		},
		"peek returns nil when plan complete": {
			actions:     []core.Action{&mockAction{}},
			curLocation: 1,
			expected:    nil,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &plan{
				Actions:     tc.actions,
				curLocation: tc.curLocation,
			}
			// Use DeepEqual for interface comparison.
			assert.True(t, reflect.DeepEqual(tc.expected, p.PeekAction()))
		})
	}
}

// TestPlan_PopAction tests that PopAction returns actions in order and advances curLocation.
func TestPlan_PopAction(t *testing.T) {
	type testCase struct {
		actions          []core.Action
		expectedSequence []core.Action
		finalLocation    int
	}

	tests := map[string]testCase{
		"pop returns actions sequentially and then nil": {
			actions:          []core.Action{&mockAction{}, &mockNullAction{}},
			expectedSequence: []core.Action{&mockAction{}, &mockNullAction{}, nil},
			finalLocation:    2,
		},
		"pop on empty plan returns nil": {
			actions:          []core.Action{},
			expectedSequence: []core.Action{nil},
			finalLocation:    0,
		},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			p := &plan{
				Actions:     tc.actions,
				curLocation: 0,
			}
			var results []core.Action
			// Call PopAction len(expectedSequence) times.
			for i := 0; i < len(tc.expectedSequence); i++ {
				results = append(results, p.PopAction())
			}
			assert.Equal(t, tc.expectedSequence, results)
			assert.Equal(t, tc.finalLocation, p.curLocation)
		})
	}
}
