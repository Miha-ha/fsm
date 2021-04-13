package fsm

import (
	"log"
	"testing"
)

type firstState struct {
	fsm *FSM
}

type secondState struct {
	fsm *FSM
}

var (
	first  = &firstState{}
	second = &secondState{}
)

func (f *firstState) SetStateMachine(fsm *FSM) {
	f.fsm = fsm
}

func (f *firstState) StateMachine() *FSM {
	return f.fsm
}

func (f *firstState) IsValidNextState(state State) bool {
	_, ok := state.(*secondState)
	return ok
}

func (f *firstState) DidEnter(from State) {
	log.Printf("[first state] Did enter from %T", from)
}

func (f *firstState) Process() {
	log.Println("[first state] Process...")
}

func (f *firstState) WillExit(to State) {
	log.Printf("[first state] WillExit to %T", to)
}

func (f *secondState) SetStateMachine(fsm *FSM) {
	f.fsm = fsm
}

func (f *secondState) StateMachine() *FSM {
	return f.fsm
}

func (f *secondState) IsValidNextState(state State) bool {
	return false
}

func (f *secondState) DidEnter(from State) {
	log.Printf("[second state] Did enter from %T", from)
}

func (f *secondState) Process() {
	log.Println("[second state] Process...")
}

func (f *secondState) WillExit(to State) {
	log.Printf("[second state] WillExit to %T", to)
}

func TestFSM_CanEnterState(t *testing.T) {
	fsm := NewFSM([]State{first, second})

	type args struct {
		state State
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
			args: args{first},
			want: true,
		},
		{
			name: "entering to the second state",
			fsm:  fsm,
			args: args{second},
			want: true,
		},
		{
			name: "entering to the first state again",
			fsm:  fsm,
			args: args{first},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			if got := tt.fsm.CanEnterState(tt.args.state); got != tt.want {
				t.Errorf("FSM.CanEnterState() = %v, want %v", got, tt.want)
			}
			tt.fsm.Enter(tt.args.state)
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
			fsm := NewFSM(tt.states)
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
		to   State
		want bool
	}{
		{
			name: "empty",
			fsm:  NewFSM([]State{}),
			to:   first,
			want: false,
		},
		{
			name: "not in states",
			fsm:  NewFSM([]State{first}),
			to:   second,
			want: false,
		},
		{
			name: "in states",
			fsm:  NewFSM([]State{first}),
			to:   first,
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
