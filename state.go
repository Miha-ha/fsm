package fsm

type State interface {
	IsValidNextState(State) bool
	DidEnter(from State)
	Process() State
	WillExit(to State)
}

type BaseState struct{}

func (bs *BaseState) DidEnter(from State) {}

func (bs *BaseState) Process() State {
	return nil
}

func (bs *BaseState) WillExit(to State) {}
