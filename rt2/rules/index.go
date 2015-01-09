package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
)

func indexSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	i := rt2.NodeOf(f)

	switch i.Right().(type) {
	case node.ConstantNode:
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.DataOf(f.Parent())[i] = i.Right().(node.ConstantNode).Data()
			return frame.End()
		}
		ret = frame.NOW
	case node.VariableNode:
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			rt2.DataOf(f.Parent())[i] = sc.Select(scope.Designator(i.Right()))
			return frame.End()
		}
		ret = frame.NOW
	case node.OperationNode, node.CallNode:
		rt2.Push(rt2.New(i.Right()), f)
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.DataOf(f.Parent())[i] = rt2.DataOf(f)[i.Right()]
			return frame.End()
		}
		ret = frame.LATER
	default:
		panic(fmt.Sprintln("unsupported type", reflect.TypeOf(i.Right())))
	}
	return seq, ret
}
