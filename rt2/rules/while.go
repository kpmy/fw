package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"reflect"
)

func whileSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)

	var cond func(f frame.Frame) (frame.Sequence, frame.WAIT)
	next := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		done := fu.DataOf(f)[n.Left()].(bool)
		fu.DataOf(f)[n.Left()] = nil
		if done && n.Right() != nil {
			fu.Push(fu.New(n.Right()), f)
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
			fu.Push(fu.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				fu.DataOf(f.Parent())[n] = fu.DataOf(f)[n.Left()]
				return next, frame.LATER
			}
			ret = frame.LATER
			return seq, ret
		default:
			panic(fmt.Sprintf("unknown condition expression", reflect.TypeOf(n.Left())))
		}
	}
	return cond, frame.NOW
}
