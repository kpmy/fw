package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/statement"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"reflect"
)

func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	a := fu.NodeOf(f)

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		switch a.Right().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(a.Left().Object(), func(interface{}) interface{} {
					return a.Right().(node.ConstantNode).Data()
				})
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(a.Left().Object(), func(interface{}) interface{} {
					return sc.Select(a.Right().Object())
				})
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode:
			fu.Push(fu.New(a.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(a.Left().Object(), func(interface{}) interface{} {
					return fu.DataOf(f)[a.Right()]
				})
				return frame.End()
			}
			ret = frame.LATER
		default:
			fmt.Println(reflect.TypeOf(a.Right()))
			panic("wrong right")
		}
		return seq, ret
	}
	switch a.(node.AssignNode).Statement() {
	case statement.ASSIGN:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			seq, ret = right(f)
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong left")
		}
	default:
		panic("wrong statement")
	}
	return seq, ret
}
