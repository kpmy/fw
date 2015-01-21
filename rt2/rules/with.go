package rules

import (
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
)

func withSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	rt2.RegOf(f)[0] = n.Left() //if
	rt2.Push(rt2.New(n.Left()), f)
	rt2.Assert(f, func(f frame.Frame) (bool, int) {
		return rt2.ValueOf(f)[n.Left().Adr()] != nil, 60
	})
	seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		last := rt2.RegOf(f)[0].(node.Node)
		done := scope.GoTypeFrom(rt2.ValueOf(f)[last.Adr()]).(bool)
		rt2.ValueOf(f)[last.Adr()] = nil
		if done && last.Right() != nil {
			rt2.Push(rt2.New(last.Right()), f)
			return frame.Tail(frame.STOP), frame.LATER
		} else if last.Right() == nil {
			return frame.End()
		} else if last.Link() != nil { //elsif
			rt2.RegOf(f)[0] = last.Link()
			rt2.Push(rt2.New(last.Link()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[last.Link().Adr()] != nil, 61
			})
			return seq, frame.LATER
		} else if n.Right() != nil { //else
			rt2.Push(rt2.New(n.Right()), f)
			return frame.Tail(frame.STOP), frame.LATER
		} else if n.Right() == nil {
			return frame.End()
		} else if last == n.Right() {
			return frame.End()
		} else {
			panic("conditional sequence wrong")
		}
	}
	return seq, frame.LATER
}
