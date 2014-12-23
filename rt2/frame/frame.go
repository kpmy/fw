package frame

import (
	"fw/rt2/context"
)

type WAIT int

const (
	WRONG WAIT = iota
	STOP
	LATER
	NOW
)

// LIFO-стек, позволяет затолкнуть фрейм связанный с другим фреймом
type Stack interface {
	PushFor(f, parent Frame)
	Pop()
	Top() Frame
}

//фрейм
type Frame interface {
	Do() WAIT
	OnPush(root Stack, parent Frame)
	OnPop()
	Parent() Frame
	Root() Stack
	context.ContextAware
}

//пользовательская функция, которую выполнит фрейм, может поставить на очередь выполнения себя или другую функцию
type Sequence func(f Frame) (Sequence, WAIT)

func Tail(x WAIT) (seq Sequence) {
	return func(f Frame) (Sequence, WAIT) { return nil, x }
}

func End() (Sequence, WAIT) { return nil, STOP }
