package agent

import (
	"github.com/stretchr/testify/assert"
	"testing"
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

func TestAgent_SetCurState(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}
	type args struct {
		state State
	}
	tests := map[string]struct {
		fields fields
		args   args
	}{
		"can set current state performing": {
			fields: fields{
				name:     "test",
				behavior: &Behavior{},
			},
			args: args{
				state: &Performing{},
			},
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			a.SetCurState(tt.args.state)
			assert.Equal(t, tt.args.state, a.behavior.curState)
		})
	}
}

func TestAgent_Plan(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}

	tests := map[string]struct {
		fields fields
		plan   Plan
	}{
		"can get plan": {
			fields: fields{
				name: "test",
				behavior: &Behavior{
					curPlan: &mockPlan{},
				},
			},
			plan: &mockPlan{},
		},
	}

	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			p := a.Plan()
			assert.Equal(t, tt.plan, p)
		})
	}
}
