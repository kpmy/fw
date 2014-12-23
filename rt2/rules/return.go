package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"reflect"
)

func returnSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	var fu nodeframe.FrameUtils
	a := fu.NodeOf(f)

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		switch a.Left().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				fu.DataOf(f.Parent())[a.Object()] = a.Left().(node.ConstantNode).Data()
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f.Parent())[a.Object()] = sc.Select(a.Left().Object())
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode:
			fu.Push(fu.New(a.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				fu.DataOf(f.Parent())[a.Object()] = fu.DataOf(f)[a.Left()]
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
