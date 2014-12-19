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

func int32Of(x interface{}) (a int32) {
	fmt.Println(reflect.TypeOf(x))
	switch x.(type) {
	case *int32:
		z := *x.(*int32)
		a = z
	case int32:
		a = x.(int32)
	default:
		panic("unsupported type")
	}
	return a
}

func sum(_a interface{}, _b interface{}) interface{} {
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a + b
}

func opSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils

	m := new(frame.SetDataMsg)
	m.Data = make(map[interface{}]interface{}, 2)
	f.(context.ContextAware).Handle(m)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.PLUS:
			fu.DataOf(f.Parent())[n] = sum(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		default:
			panic("unknown operation")
		}
	}

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Right().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[n.Right()] = n.Right().(node.ConstantNode).Data()
			return op, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Right()] = sc.Select(n.Right().Object())
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
			fu.DataOf(f)[n.Left()] = n.Left().(node.ConstantNode).Data()
			return right, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(n.Left().Object())
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
