package fsm

type StateType int

type State interface {
	Type() StateType
	SetStateMachine(bssm *FSM)
	StateMachine() *FSM
	IsValidNextState(StateType) bool
	DidEnter(from State)
	Process()
	WillExit(to State)
}

type BaseState struct {
	fsm *FSM
}

func (bs *BaseState) SetStateMachine(fsm *FSM) {
	bs.fsm = fsm
}

func (bs *BaseState) StateMachine() *FSM {
	return bs.fsm
}

func (bs *BaseState) DidEnter(from State) {}

func (bs *BaseState) Process() {}

func (bs *BaseState) WillExit(to State) {}
