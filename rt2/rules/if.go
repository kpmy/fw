package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"reflect"
)

func ifExpr(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	switch n.Left().(type) {
	case node.OperationNode:
		fu.Push(fu.New(n.Left()), f)
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			fu.DataOf(f.Parent())[n] = fu.DataOf(f)[n.Left()]
			return frame.End()
		}
		ret = frame.LATER
	default:
		panic(fmt.Sprintf("unknown condition expression", reflect.TypeOf(n.Left())))
	}
	return seq, ret
}

func ifSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	fu.DataOf(f)[0] = n.Left() //if
	fu.Push(fu.New(n.Left()), f)
	seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		last := fu.DataOf(f)[0].(node.Node)
		done := fu.DataOf(f)[last].(bool)
		fu.DataOf(f)[last] = nil
		if done {
			fu.Push(fu.New(last.Right()), f)
			return frame.Tail(frame.STOP), frame.LATER
		} else if last.Link() != nil { //elsif
			fu.DataOf(f)[0] = last.Link()
			fu.Push(fu.New(last.Link()), f)
			return seq, frame.LATER
		} else if n.Right() != nil { //else
			fu.Push(fu.New(n.Right()), f)
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
