package rules

import (
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

type Do func(...IN) OUT

type IN struct {
	frame frame.Frame
}

type OUT struct {
	do   Do
	next WAIT
}

func (n WAIT) wait() frame.WAIT {
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

func waiting(n frame.WAIT) WAIT {
	switch n {
	case frame.WRONG:
		return WRONG
	case frame.STOP:
		return STOP
	case frame.LATER:
		return LATER
	case frame.NOW:
		return NOW
	case frame.BEGIN:
		return BEGIN
	case frame.END:
		return END
	default:
		panic(n)
	}
}

func End() OUT {
	return OUT{next: STOP}
}

func Tail(x WAIT) Do {
	return func(...IN) OUT { return OUT{next: x} }
}

func This(o OUT) (seq frame.Sequence, ret frame.WAIT) {
	ret = o.next.wait()
	if ret != frame.STOP {
		seq = Propose(o.do)
	}
	return seq, ret
}

func Propose(a Do) frame.Sequence {
	return func(fr frame.Frame) (frame.Sequence, frame.WAIT) {
		return This(a(IN{frame: fr}))
	}
}

func Expose(f frame.Sequence) Do {
	return func(in ...IN) (out OUT) {
		s, w := f(in[0].frame)
		return OUT{do: Expose(s), next: waiting(w)}
	}
}
