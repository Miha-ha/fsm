package fsm

type State interface {
	SetStateMachine(fsm *FSM)
	StateMachine() *FSM
	IsValidNextState(State) bool
	DidEnter(from State)
	Process()
	WillExit(to State)
}
