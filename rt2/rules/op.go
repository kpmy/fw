package rules

import (
	"cp/constant/operation"
	"cp/node"
	"cp/object"
	"fmt"
	"reflect"
	"rt2/context"
	"rt2/frame"
	"rt2/nodeframe"
	"rt2/scope"
	"ypk/assert"
)

func int32Of(x interface{}) (a int32) {
	//fmt.Println(reflect.TypeOf(x))
	switch x.(type) {
	case *int32:
		z := *x.(*int32)
		a = z
	case int32:
		a = x.(int32)
	default:
		panic(fmt.Sprintln("unsupported type", reflect.TypeOf(x)))
	}
	return a
}

func sum(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a + b
}

func mopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f).(node.MonadicNode)
	switch n.Operation() {
	case operation.CONVERT:
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		switch n.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			x := sc.Select(n.Left().Object())
			assert.For(x != nil, 40)
			switch n.Type() {
			case object.INTEGER:
				switch x.(type) {
				case int8:
					fu.DataOf(f.Parent())[n] = int32(x.(int8))
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			default:
				panic("wrong type")
			}
			return frame.End()
		default:
			panic(fmt.Sprintln("unsupported left", reflect.TypeOf(n.Left())))

		}
	default:
		panic(fmt.Sprintln("no such operation", n.(node.MonadicNode).Operation()))
	}

}

func dopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils

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
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Right()] = sc.Select(n.Right().Object())
				//fmt.Println(n.Right().Object(), reflect.TypeOf(n.Right().Object()))
				assert.For(fu.DataOf(f)[n.Right()] != nil, 60)
				return op, frame.DO
			}
			ret = frame.DO
			return seq, ret
		case node.OperationNode:
			fu.Push(fu.New(n.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.DO
			}
			ret = frame.SKIP
			return seq, ret
		default:
			fmt.Println(reflect.TypeOf(n.Right()))
			panic("wrong right")
		}
	}

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Left().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[n.Left()] = n.Left().(node.ConstantNode).Data()
			return right, frame.DO
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(n.Left().Object())
				return right, frame.DO
			}
			ret = frame.DO
			return seq, ret
		case node.OperationNode:
			fu.Push(fu.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return right, frame.DO
			}
			ret = frame.SKIP
			return seq, ret
		default:
			fmt.Println(reflect.TypeOf(n.Left()))
			panic("wrong left")
		}
	}
	return left, frame.DO
}
