package fsm

import "sync"

type FSM struct {
	muCurrentState sync.RWMutex
	currentState   State
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

func (f *FSM) CurrentState() State {
	f.muCurrentState.RLock()
	defer f.muCurrentState.RUnlock()
	return f.currentState
}

func (f *FSM) CanEnterState(state State) bool {
	currentState := f.CurrentState()

	if state == nil {
		return false
	}

	if state.StateMachine() != f {
		return false
	}

	if currentState == nil {
		return true
	}

	return currentState.IsValidNextState(state)
}

func (f *FSM) Enter(state State) bool {
	if state == nil || !f.CanEnterState(state) {
		return false
	}

	currentState := f.CurrentState()

	if currentState != nil {
		currentState.WillExit(state)
	}

	f.muCurrentState.Lock()
	f.currentState = state
	f.muCurrentState.Unlock()

	state.DidEnter(currentState)
	state.Process()

	return true
}
