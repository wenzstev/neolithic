package agent

import (
	"Neolithic/internal/planner"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAgent_GetCurrentAction(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}
	tests := []struct {
		name   string
		fields fields
		want   planner.Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			assert.Equalf(t, tt.want, a.GetCurrentAction(), "GetCurrentAction()")
		})
	}
}

func TestAgent_Name(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			assert.Equalf(t, tt.want, a.Name(), "Name()")
		})
	}
}

func TestAgent_PopCurrentAction(t *testing.T) {
	type fields struct {
		name     string
		behavior *Behavior
	}
	tests := []struct {
		name   string
		fields fields
		want   planner.Action
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			assert.Equalf(t, tt.want, a.PopCurrentAction(), "PopCurrentAction()")
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
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			a := &agent{
				name:     tt.fields.name,
				behavior: tt.fields.behavior,
			}
			a.SetCurState(tt.args.state)
		})
	}
}
