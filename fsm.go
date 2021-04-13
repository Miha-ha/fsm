package fsm

import "sync"

type FSM struct {
	muCurrentState sync.RWMutex
	CurrentState   State
}

func NewFSM(states []State) *FSM {
	fsm := &FSM{}

	for _, state := range states {
		if state != nil {
			state.SetStateMachine(fsm)
		}
	}

	return fsm
}

func (f *FSM) CanEnterState(state State) bool {
	f.muCurrentState.RLock()
	defer f.muCurrentState.RUnlock()

	if state == nil {
		return false
	}

	if state.StateMachine() != f {
		return false
	}

	if f.CurrentState == nil {
		return true
	}

	return f.CurrentState.IsValidNextState(state)
}

func (f *FSM) Enter(state State) bool {
	if state == nil || !f.CanEnterState(state) {
		return false
	}

	f.muCurrentState.Lock()
	defer f.muCurrentState.Unlock()

	if f.CurrentState != nil {
		f.CurrentState.WillExit(state)
	}

	prevState := f.CurrentState
	f.CurrentState = state

	f.CurrentState.DidEnter(prevState)
	f.CurrentState.Process()

	return true
}
