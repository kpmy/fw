package rules

import (
	"cp/constant/operation"
	"cp/node"
	"fmt"
	"reflect"
	"rt2/context"
	"rt2/frame"
	"rt2/nodeframe"
	"rt2/scope"
)

func intOf(x interface{}) (a int) {
	fmt.Println(reflect.TypeOf(x))
	switch x.(type) {
	case *int:
		z := *x.(*int)
		a = z
	case int:
		a = x.(int)
	default:
		panic("unsupported type")
	}
	return a
}

func sum(_a interface{}, _b interface{}) interface{} {
	var a int = intOf(_a)
	var b int = intOf(_b)
	return a + b
}

func opSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils

	m := new(frame.SetDataMsg)
	m.Data = make([]interface{}, 2)
	f.(context.ContextAware).Handle(m)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.PLUS:
			fu.DataOf(f.Parent())[0] = sum(fu.DataOf(f)[0], fu.DataOf(f)[1])
			return frame.End()
		default:
			panic("unknown operation")
		}
	}

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Right().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[1] = n.Right().(node.ConstantNode).Data()
			return op, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[1] = sc.Select(n.Right().Object())
				return op, frame.DO
			}
			ret = frame.DO
			return seq, ret
		default:
			panic("wrong right")
		}
	}

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Left().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[0] = n.Left().(node.ConstantNode).Data()
			return right, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[0] = sc.Select(n.Left().Object())
				return right, frame.DO
			}
			ret = frame.DO
			return seq, ret
		default:
			panic("wrong left")
		}
	}

	return left, frame.DO
}
