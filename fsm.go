package fsm

import "sync"

type FSM struct {
	muCurrentState sync.RWMutex
	currentState   State
	states         map[StateType]State
}

func NewFSM(states ...State) *FSM {
	fsm := &FSM{
		states: map[StateType]State{},
	}

	for _, state := range states {
		if state != nil {
			state.SetStateMachine(fsm)
			fsm.states[state.Type()] = state
		}
	}

	return fsm
}

func (f *FSM) CurrentState() State {
	f.muCurrentState.RLock()
	defer f.muCurrentState.RUnlock()
	return f.currentState
}

func (f *FSM) CanEnterState(stateType StateType) bool {
	currentState := f.CurrentState()

	state, ok := f.states[stateType]
	if !ok {
		return false
	}

	if state == nil {
		return false
	}

	if state.StateMachine() != f {
		return false
	}

	if currentState == nil {
		return true
	}

	return currentState.IsValidNextState(stateType)
}

func (f *FSM) Enter(stateType StateType) bool {
	if !f.CanEnterState(stateType) {
		return false
	}

	state := f.states[stateType]
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
