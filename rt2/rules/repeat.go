package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"reflect"
)

func repeatSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)

	fu.DataOf(f)[n.Right()] = false
	var cond func(f frame.Frame) (frame.Sequence, frame.WAIT)
	next := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		done := fu.DataOf(f)[n.Right()].(bool)
		fu.DataOf(f)[n.Right()] = nil
		if !done && n.Right() != nil {
			fu.Push(fu.New(n.Left()), f)
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
			fu.Push(fu.New(n.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				fu.DataOf(f.Parent())[n] = fu.DataOf(f)[n.Right()]
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
