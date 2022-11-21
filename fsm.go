package fsm

import (
	"fmt"
	"reflect"
	"sync"
)

type FSM struct {
	muCurrentState sync.RWMutex
	currentState   State
	nextState      State
	states         map[reflect.Type]State
}

func NewFSM(states ...State) *FSM {
	fsm := &FSM{
		states: map[reflect.Type]State{},
	}

	for _, state := range states {
		if state != nil {
			if fsm.nextState == nil {
				fsm.nextState = state
			}
			fsm.states[reflect.TypeOf(state)] = state
		}
	}

	return fsm
}

func (f *FSM) Run(enter ...State) error {
	var ok bool
	if len(enter) > 0 {
		if f.nextState, ok = f.State(enter[0]); !ok {
			return fmt.Errorf("state '%T' not found in FSM ", enter)
		}
	}

	for {

		from := f.CurrentState()

		if !f.CanEnterState(f.nextState) {
			return fmt.Errorf("can't enter state %T from %T", f.nextState, from)
		}

		current := f.nextState

		f.muCurrentState.Lock()
		f.currentState = current
		f.muCurrentState.Unlock()

		current.DidEnter(from)
		next := current.Process()

		switch next.(type) {
		case nil:
			f.nextState = nil
		default:
			if f.nextState, ok = f.State(next); !ok {
				return fmt.Errorf("state '%T' not found in FSM ", next)
			}
		}

		current.WillExit(f.nextState)

		if f.nextState == nil {
			return nil
		}

	}
}

func (f *FSM) State(stateType State) (State, bool) {
	state, ok := f.states[reflect.TypeOf(stateType)]
	return state, ok
}

func (f *FSM) CurrentState() State {
	f.muCurrentState.RLock()
	defer f.muCurrentState.RUnlock()
	return f.currentState
}

func (f *FSM) CanEnterState(stateType State) bool {
	currentState := f.CurrentState()

	_, ok := f.State(stateType)
	if !ok {
		return false
	}

	if currentState == nil {
		return true
	}

	return currentState.IsValidNextState(stateType)
}
