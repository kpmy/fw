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
	//спец коды для начальной и конечной стадии
	BEGIN
	END
)

func (w WAIT) String() string {
	switch w {
	case NOW:
		return "NOW"
	case LATER:
		return "LATER"
	case STOP:
		return "STOP"
	case WRONG:
		return "WRONG"
	case BEGIN:
		return "BEGIN"
	case END:
		return "END"
	default:
		panic("wrong wait value")
	}
}

// LIFO-стек, позволяет затолкнуть фрейм связанный с другим фреймом
// добавлен Queue, что-то типа механизма прерываний, позволяет передать внеочередной фрейм на выполнение
type Stack interface {
	PushFor(f, parent Frame)
	Pop()
	Top() Frame
	ForEach(run func(this Frame) bool)
	Queue(...Frame) Frame
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
type Assert func(f Frame) (bool, int)

func Tail(x WAIT) (seq Sequence) {
	return func(f Frame) (Sequence, WAIT) { return nil, x }
}

func End() (Sequence, WAIT) { return nil, STOP }
