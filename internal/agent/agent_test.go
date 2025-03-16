package agent

import (
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
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			assert.Equalf(t, tt.want, a.Name(), "Name()")
		})
	}
}

func TestAgent_Behavior(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}

	tests := map[string]struct {
		fields   fields
		behavior *Behavior
	}{
		"can get plan": {
			fields: fields{
				name: "test",
				behavior: &Behavior{
					CurPlan: &mockPlan{},
				},
			},
			behavior: &Behavior{
				CurPlan: &mockPlan{},
			},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			b := a.Behavior()
			assert.Equal(t, tt.behavior, b)
		})
	}
}
