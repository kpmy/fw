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

func returnSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	a := rt2.NodeOf(f)

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		switch a.Left().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//rt2.DataOf(f.Parent())[a.Object()] = a.Left().(node.ConstantNode).Data()
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = sc.Provide(a.Left())(nil)
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = sc.Select(a.Left().Object().Adr())
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode:
			rt2.Push(rt2.New(a.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = rt2.ValueOf(f)[a.Left().Adr()]
				return frame.End()
			}
			ret = frame.LATER
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong right")
		}
		return seq, ret
	}
	return left, frame.NOW
}
