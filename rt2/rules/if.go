package rules

import (
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/halt"
)

func ifExpr(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	switch l := n.Left().(type) {
	case node.OperationNode:
		rt2.Push(rt2.New(n.Left()), f)
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.RegOf(f.Parent())[n] = rt2.RegOf(f)[n.Left()]
			rt2.ValueOf(f.Parent())[n.Adr()] = rt2.ValueOf(f)[n.Left().Adr()]
			return frame.End()
		}
		ret = frame.LATER
	case node.ConstantNode, node.VariableNode, node.ParameterNode, node.FieldNode:
		return This(expectExpr(f, l, func(...IN) OUT {
			rt2.ValueOf(f.Parent())[n.Adr()] = rt2.ValueOf(f)[l.Adr()]
			return End()
		}))
	default:
		halt.As(100, reflect.TypeOf(l))
	}
	return seq, ret
}

func ifSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	rt2.RegOf(f)[0] = n.Left() //if
	rt2.Push(rt2.New(n.Left()), f)
	rt2.Assert(f, func(f frame.Frame) (bool, int) {
		return rt2.ValueOf(f)[n.Left().Adr()] != nil, 60
	})
	seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		last := rt2.RegOf(f)[0].(node.Node)
		done := scope.GoTypeFrom(rt2.ValueOf(f)[last.Adr()]).(bool)
		rt2.RegOf(f)[last] = nil
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
