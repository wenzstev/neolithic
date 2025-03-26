package agent

import (
	"Neolithic/internal/core"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAgent_Name(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}
	tests := map[string]struct {
		fields fields
		want   string
	}{
		"can provide name": {
			fields: fields{
				name:     "test",
				behavior: &Behavior{},
			},
			want: "test",
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a := &Agent{
				name:     tt.fields.name,
				Behavior: tt.fields.behavior,
			}
			assert.Equalf(t, tt.want, a.Name(), "Name()")
		})
	}
}

func TestAgent_Inventory(t *testing.T) {
	t.Run("Can get inventory", func(t *testing.T) {
		testInventory := core.NewInventory()
		testInventory.AdjustAmount(testResource, 5)
		agent := &Agent{
			name:      "test",
			inventory: testInventory,
		}

		assert.Equal(t, agent.Inventory(), testInventory)
	})
}

func TestAgent_DeepCopy(t *testing.T) {
	testInventory := core.NewInventory()
	testInventory.AdjustAmount(testResource, 5)
	testAgent := &Agent{
		name:      "test",
		inventory: testInventory,
	}

	agentCopy := testAgent.DeepCopy()
	assert.True(t, reflect.DeepEqual(testAgent, agentCopy))
	assert.False(t, testAgent == agentCopy)

}

func TestAgent_String(t *testing.T) {
	testInventory := core.NewInventory()
	testInventory.AdjustAmount(testResource, 5)
	testAgent := &Agent{
		name:      "test",
		inventory: testInventory,
	}

	assert.Equal(t, "Agent: test \nInventory   testResource: 5\n\n Position {0 0}\n", testAgent.String())
}
