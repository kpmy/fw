package rules

import (
	"fmt"
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"math/big"
	"reflect"
	"ypk/assert"
)

func int32Of(x interface{}) (a int32) {
	//fmt.Println(reflect.TypeOf(x))
	switch v := x.(type) {
	case *int32:
		z := *x.(*int32)
		a = z
	case int32:
		a = x.(int32)
	case *big.Int:
		a = int32(v.Int64())
	default:
		//panic(fmt.Sprintln("unsupported type", reflect.TypeOf(x)))
	}
	return a
}

func mopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	n := rt2.NodeOf(f).(node.MonadicNode)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.NOT:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Not(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		case operation.IS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Is(sc.Select(n.Left().Object().Adr()), n.Object().Complex())
			return frame.End()
		case operation.ABS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Abs(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		case operation.ODD:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Odd(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		case operation.CAP:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Cap(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		case operation.BITS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Bits(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		case operation.MINUS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Minus(rt2.ValueOf(f)[n.Left().Adr()])
			return frame.End()
		default:
			panic("no such op")
		}

	}

	switch n.Operation() {
	case operation.ALIEN_CONV:
		conv := func(x scope.Value) {
			if n.Type() != object.NOTYPE {
				rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Conv(x, n.Type())
			} else {
				rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Conv(x, n.Type(), n.(node.MonadicNode).Complex())
			}
		}
		switch n.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			x := sc.Select(n.Left().Object().Adr())
			assert.For(x != nil, 40)
			conv(scope.ValueFrom(x))
			return frame.End()
		case node.OperationNode, node.DerefNode, node.CallNode, node.IndexNode:
			return This(expectExpr(f, n.Left(), func(...IN) OUT {
				conv(rt2.ValueOf(f)[n.Left().Adr()])
				return End()
			}))
		default:
			panic(fmt.Sprintln("unsupported left", reflect.TypeOf(n.Left())))
		}
	case operation.NOT:
		switch n.Left().(type) {
		case node.ConstantNode:
			rt2.ValueOf(f)[n.Left().Adr()] = sc.Provide(n.Left())(nil)
			return op, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f)[n.Left().Adr()] = sc.Select(n.Left().Object().Adr())
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode:
			rt2.Push(rt2.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.CallNode, node.FieldNode:
			return This(expectExpr(f, n.Left(), Expose(op)))
		default:
			fmt.Println(reflect.TypeOf(n.Left()))
			panic("wrong left")
		}
	case operation.IS:
		switch n.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		default:
			fmt.Println(reflect.TypeOf(n.Left()))
			panic("wrong left")
		}
	case operation.ABS, operation.ODD, operation.CAP, operation.BITS, operation.MINUS:
		return This(expectExpr(f, n.Left(), Expose(op)))
	default:
		panic(fmt.Sprintln("no such operation", n.(node.MonadicNode).Operation()))
	}
}

func dopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f).(node.OperationNode)
		switch n.Operation() {
		case operation.PLUS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Sum(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.MINUS:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Sub(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.EQUAL:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Eq(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.LESSER:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Lss(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.LESS_EQUAL:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Leq(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.LEN:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Len(n.Left().Object(), rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.NOT_EQUAL:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Neq(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.GREATER:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Gtr(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.MAX:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Max(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.MIN:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Min(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.DIV:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Div(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.MOD:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Mod(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.TIMES:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Mult(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.SLASH:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Divide(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.IN:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.In(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.ASH:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Ash(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.AND:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.And(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.OR:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Or(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		case operation.GREAT_EQUAL:
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Geq(rt2.ValueOf(f)[n.Left().Adr()], rt2.ValueOf(f)[n.Right().Adr()])
			return frame.End()
		default:
			panic(fmt.Sprintln("unknown operation", n.(node.OperationNode).Operation()))
		}
	}

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f)
		switch r := n.Right().(type) {
		case node.ConstantNode:
			rt2.ValueOf(f)[n.Right().Adr()] = sc.Provide(n.Right())(nil)
			return op, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f)[n.Right().Adr()] = sc.Select(n.Right().Object().Adr())
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode:
			rt2.Push(rt2.New(n.Right()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[r.Adr()] != nil, 61
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.CallNode, node.FieldNode:
			return This(expectExpr(f, r, Expose(op)))
		/*case node.FieldNode:
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			//				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			//				rt2.DataOf(f)[n.Right()] = sc.Select(scope.Designator(n.Right()))
			panic(0)
			return op, frame.NOW
		}
		ret = frame.NOW
		return seq, ret*/
		default:
			fmt.Println(reflect.TypeOf(n.Right()))
			panic("wrong right")
		}
	}

	short := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		n := rt2.NodeOf(f).(node.OperationNode)
		switch n.Operation() {
		case operation.AND:
			val := scope.GoTypeFrom(rt2.ValueOf(f)[n.Left().Adr()]).(bool)
			if val {
				return right, frame.NOW
			} else {
				rt2.ValueOf(f.Parent())[n.Adr()] = scope.TypeFromGo(false)
				return frame.End()
			}
		case operation.OR:
			val := scope.GoTypeFrom(rt2.ValueOf(f)[n.Left().Adr()]).(bool)
			if !val {
				return right, frame.NOW
			} else {
				rt2.ValueOf(f.Parent())[n.Adr()] = scope.TypeFromGo(true)
				return frame.End()
			}
		default:
			return right, frame.NOW
		}
	}

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f)
		switch l := n.Left().(type) {
		case node.ConstantNode:
			rt2.ValueOf(f)[n.Left().Adr()] = sc.Provide(n.Left())(nil)
			return short, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f)[n.Left().Adr()] = sc.Select(n.Left().Object().Adr())
				return short, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode, node.RangeNode:
			rt2.Push(rt2.New(n.Left()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 60
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return short, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.CallNode:
			return This(expectExpr(f, l, Expose(short)))
		case node.FieldNode:
			rt2.Push(rt2.New(l), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 60
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//sc.Select(l.Left().Object().Adr(), func(v scope.Value) {
				//	rt2.ValueOf(f)[n.Left().Adr()] = v.(scope.Record).Get(l.Object().Adr())
				//})
				return short, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.IndexNode:
			return This(expectExpr(f, n.Left(), Expose(short)))
		default:
			fmt.Println(reflect.TypeOf(n.Left()))
			panic("wrong left")
		}
	}
	return left, frame.NOW
}
