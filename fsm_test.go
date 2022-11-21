package fsm

import (
	"testing"
)

type firstState struct {
	BaseState
	logf func(format string, args ...any)
}

type secondState struct {
	BaseState
	logf func(format string, args ...any)
}

type thirdState struct {
	BaseState
	logf func(format string, args ...any)
}

// first state

func (f *firstState) IsValidNextState(state State) bool {
	switch state.(type) {
	case *secondState:
		return true
	}

	return false
}

func (f *firstState) DidEnter(from State) {
	f.logf("[first state] DidEnter from %v", from)
}

func (f *firstState) Process() State {
	f.logf("[first state] Process...")
	return (*secondState)(nil)
}

func (f *firstState) WillExit(to State) {
	f.logf("[first state] WillExit to %T", to)
}

// second state

func (f *secondState) IsValidNextState(stateType State) bool {
	return false
}

func (f *secondState) DidEnter(from State) {
	f.logf("[second state] DidEnter from %T", from)
}

func (f *secondState) Process() State {
	f.logf("[second state] Process...")
	return nil
}

func (f *secondState) WillExit(to State) {
	f.logf("[second state] WillExit to %v", to)
}

// third state

func (f *thirdState) IsValidNextState(stateType State) bool {
	return false
}

func (f *thirdState) DidEnter(from State) {
	f.logf("[third state] DidEnter from %T", from)
}

func (f *thirdState) Process() State {
	f.logf("[third state] Process...")
	return (*secondState)(nil)
}

func (f *thirdState) WillExit(to State) {
	f.logf("[third state] WillExit to %T", to)
}

func TestFSM_CanEnterState(t *testing.T) {

	type args struct {
		current State
		next    State
	}
	tests := []struct {
		name string
		fsm  *FSM
		args args
		want bool
	}{
		{
			name: "not found",
			fsm:  NewFSM(&secondState{}),
			args: args{
				current: nil,
				next:    (*firstState)(nil),
			},
			want: false,
		},
		{
			name: "entering to the first state",
			fsm:  NewFSM(&firstState{}, &secondState{}),
			args: args{
				current: nil,
				next:    (*firstState)(nil),
			},
			want: true,
		},
		{
			name: "entering to the second state",
			fsm:  NewFSM(&firstState{}, &secondState{}),
			args: args{
				current: nil,
				next:    (*secondState)(nil),
			},
			want: true,
		},
		{
			name: "entering to the first state from second",
			fsm:  NewFSM(&firstState{}, &secondState{}),
			args: args{
				current: (*secondState)(nil),
				next:    (*firstState)(nil),
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.current != nil {
				tt.fsm.currentState, _ = tt.fsm.State(tt.args.current)
			}

			if got := tt.fsm.CanEnterState(tt.args.next); got != tt.want {
				t.Errorf("FSM.CanEnterState() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFSM_Run(t *testing.T) {

	type args struct {
		states []State
		enter  State
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "empty",
			args: args{
				states: []State{},
				enter:  (*firstState)(nil),
			},
			wantErr: true,
		},
		{
			name: "not found enter state",
			args: args{
				states: []State{&firstState{logf: t.Logf}},
				enter:  (*secondState)(nil),
			},

			wantErr: true,
		},
		{
			name: "not found",
			args: args{
				states: []State{&firstState{logf: t.Logf}},
				enter:  (*firstState)(nil),
			},
			wantErr: true,
		},
		{
			name: "enter in first state",
			args: args{
				states: []State{&firstState{logf: t.Logf}, &secondState{logf: t.Logf}},
				enter:  nil,
			},
			wantErr: false,
		},
		{
			name: "chain",
			args: args{
				states: []State{&firstState{logf: t.Logf}, &secondState{logf: t.Logf}},
				enter:  (*firstState)(nil),
			},
			wantErr: false,
		},
		{
			name: "wrong transition",
			args: args{
				states: []State{&thirdState{logf: t.Logf}, &secondState{logf: t.Logf}},
				enter:  (*thirdState)(nil),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fsm := NewFSM(tt.args.states...)
			var err error
			switch tt.args.enter.(type) {
			case nil:
				err = fsm.Run()
			default:
				err = fsm.Run(tt.args.enter)
			}

			if err != nil {
				t.Logf("run error: %v", err)
			}

			if (err != nil) != tt.wantErr {
				t.Errorf("FSM.Run() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
