package rt

import "container/list"

type Wait int

const (
	WRONG Wait = iota
	SKIP
	STOP
	DO
)

type Stack interface {
	Push(frame Frame)
	Pop() Frame
	Top() Frame
}

type Frame interface {
	Do() (wait Wait)
	OnPush()
	OnPop()
}

func NewStack() Stack {
	return new(stdStack).Init()
}

type stdStack struct {
	inner list.List
}

func (s *stdStack) Init() *stdStack {
	s.inner = *list.New()
	return s
}

func (s *stdStack) Push(frame Frame) {
	s.inner.PushFront(frame)
	frame.OnPush()
}

func (s *stdStack) Pop() (frame Frame) {
	if s.inner.Front() != nil {
		elem := s.inner.Front()
		frame = elem.Value.(Frame)
		s.inner.Remove(elem)
		frame.OnPop()
	} else {
		panic("it's empty stack")
	}
	return frame
}

func (s *stdStack) Top() (frame Frame) {
	elem := s.inner.Front()
	if elem != nil {
		frame = elem.Value.(Frame)
	}
	return frame
}

func (w Wait) String() string {
	switch w {
	case DO:
		return "DO"
	case SKIP:
		return "SKIP"
	case STOP:
		return "STOP"
	case WRONG:
		return "WRONG"
	default:
		panic("wrong wait value")
	}
}
