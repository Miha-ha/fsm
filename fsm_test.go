package fsm

import (
	"log"
	"testing"
)

const (
	StateTypeFirst StateType = iota
	StateTypeSecond
)

type firstState struct {
	BaseState
}

type secondState struct {
	BaseState
}

var (
	first  = &firstState{}
	second = &secondState{}
)

// first state

func (f *firstState) Type() StateType {
	return StateTypeFirst
}

func (f *firstState) IsValidNextState(stateType StateType) bool {
	return stateType == StateTypeSecond
}

func (f *firstState) Process() {
	log.Println("[first state] Process...")
}

// second state

func (f *secondState) Type() StateType {
	return StateTypeSecond
}

func (f *secondState) IsValidNextState(stateType StateType) bool {
	return false
}

func (f *secondState) Process() {
	log.Println("[second state] Process...")
}

func TestFSM_CanEnterState(t *testing.T) {
	fsm := NewFSM(first, second)

	type args struct {
		stateType StateType
	}
	tests := []struct {
		name string
		fsm  *FSM
		args args
		want bool
	}{
		{
			name: "entering to the first state",
			fsm:  fsm,
			args: args{StateTypeFirst},
			want: true,
		},
		{
			name: "entering to the second state",
			fsm:  fsm,
			args: args{StateTypeSecond},
			want: true,
		},
		{
			name: "entering to the first state again",
			fsm:  fsm,
			args: args{StateTypeFirst},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fsm.CanEnterState(tt.args.stateType); got != tt.want {
				t.Errorf("FSM.CanEnterState() = %v, want %v", got, tt.want)
			}
			tt.fsm.Enter(tt.args.stateType)
		})
	}
}

func TestNewFSM(t *testing.T) {
	type args struct {
		states []State
	}
	tests := []struct {
		name   string
		states []State
		want   bool
	}{
		{
			name:   "inject FSM to states",
			states: []State{first, second},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := NewFSM(tt.states...)
			if fsm == nil {
				t.Errorf("NewFSM() = nil")
			}

			for _, state := range tt.states {
				if state.StateMachine() != fsm {
					t.Errorf("state %T does not contain FSM", state)
				}
			}
		})
	}
}

func TestFSM_Enter(t *testing.T) {

	tests := []struct {
		name string
		fsm  *FSM
		to   StateType
		want bool
	}{
		{
			name: "empty",
			fsm:  NewFSM(),
			to:   StateTypeFirst,
			want: false,
		},
		{
			name: "not in states",
			fsm:  NewFSM(first),
			to:   StateTypeSecond,
			want: false,
		},
		{
			name: "in states",
			fsm:  NewFSM(first),
			to:   StateTypeFirst,
			want: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fsm.Enter(tt.to); got != tt.want {
				t.Errorf("Enter(to %T) return %v, want %v", tt.to, got, tt.want)
			}
		})
	}
}
