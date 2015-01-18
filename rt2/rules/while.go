package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
)

func whileSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)

	var cond func(f frame.Frame) (frame.Sequence, frame.WAIT)
	next := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		done := scope.GoTypeFrom(rt2.ValueOf(f)[n.Left().Adr()]).(bool)
		rt2.RegOf(f)[n.Left()] = nil
		if done && n.Right() != nil {
			rt2.Push(rt2.New(n.Right()), f)
			return cond, frame.LATER
		} else if !done {
			return frame.End()
		} else if n.Right() == nil {
			return frame.End()
		} else {
			panic("unexpected while seq")
		}
	}

	cond = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		switch n.Left().(type) {
		case node.OperationNode:
			rt2.Push(rt2.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.RegOf(f.Parent())[n] = rt2.RegOf(f)[n.Left()]
				rt2.ValueOf(f.Parent())[n.Adr()] = rt2.ValueOf(f)[n.Left().Adr()]
				return next, frame.LATER
			}
			ret = frame.LATER
			return seq, ret
		default:
			panic(fmt.Sprintf("unknown while expression", reflect.TypeOf(n.Left())))
		}
	}
	return cond, frame.NOW
}
