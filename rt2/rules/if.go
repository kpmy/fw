package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"reflect"
)

func ifExpr(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	switch n.Left().(type) {
	case node.OperationNode:
		rt2.Push(rt2.New(n.Left()), f)
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.DataOf(f.Parent())[n] = rt2.DataOf(f)[n.Left()]
			return frame.End()
		}
		ret = frame.LATER
	default:
		panic(fmt.Sprintf("unknown condition expression", reflect.TypeOf(n.Left())))
	}
	return seq, ret
}

func ifSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	rt2.DataOf(f)[0] = n.Left() //if
	rt2.Push(rt2.New(n.Left()), f)
	seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		last := rt2.DataOf(f)[0].(node.Node)
		done := rt2.DataOf(f)[last].(bool)
		rt2.DataOf(f)[last] = nil
		if done && last.Right() != nil {
			rt2.Push(rt2.New(last.Right()), f)
			return frame.Tail(frame.STOP), frame.LATER
		} else if last.Right() == nil {
			return frame.End()
		} else if last.Link() != nil { //elsif
			rt2.DataOf(f)[0] = last.Link()
			rt2.Push(rt2.New(last.Link()), f)
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
