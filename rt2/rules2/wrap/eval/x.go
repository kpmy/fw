package eval

import (
	"fw/cp/node"
	"fw/rt2/frame"
)

type WAIT int

const (
	WRONG WAIT = iota
	STOP
	LATER
	NOW
	BEGIN
	END
)

func (n WAIT) Wait() frame.WAIT {
	switch n {
	case WRONG:
		return frame.WRONG
	case STOP:
		return frame.STOP
	case LATER:
		return frame.LATER
	case NOW:
		return frame.NOW
	case BEGIN:
		return frame.BEGIN
	case END:
		return frame.END
	default:
		panic(n)
	}
}

type Do func(IN) OUT

type IN struct {
	IR     node.Node
	Frame  frame.Frame
	Parent frame.Frame
	Key    interface{}
}

type OUT struct {
	Do   Do
	Next WAIT
}

func End() OUT {
	return OUT{Next: STOP}
}

func Tail(x WAIT) Do {
	return func(IN) OUT { return OUT{Next: x} }
}

func Later(x Do) OUT {
	return OUT{Do: x, Next: LATER}
}

func Now(x Do) OUT {
	return OUT{Do: x, Next: NOW}
}
