package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"reflect"
)

func repeatSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)

	rt2.DataOf(f)[n.Right()] = false
	var cond func(f frame.Frame) (frame.Sequence, frame.WAIT)
	next := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		done := rt2.DataOf(f)[n.Right()].(bool)
		rt2.DataOf(f)[n.Right()] = nil
		if !done && n.Left() != nil {
			rt2.Push(rt2.New(n.Left()), f)
			return cond, frame.LATER
		} else if done {
			return frame.End()
		} else if n.Left() == nil {
			return frame.End()
		} else {
			panic("unexpected repeat seq")
		}
	}

	cond = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		switch n.Right().(type) {
		case node.OperationNode:
			rt2.Push(rt2.New(n.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.DataOf(f.Parent())[n] = rt2.DataOf(f)[n.Right()]
				return next, frame.LATER
			}
			ret = frame.LATER
			return seq, ret
		default:
			panic(fmt.Sprintf("unknown repeat expression", reflect.TypeOf(n.Left())))
		}
	}
	return next, frame.NOW
}
