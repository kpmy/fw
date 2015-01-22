package rules

import (
	"fmt"
	"fw/cp/constant/operation"
	"fw/cp/node"
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
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Is(sc.Select(n.Left().Object().Adr()), n.Object())
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
		default:
			panic("no such op")
		}

	}

	switch n.Operation() {
	case operation.ALIEN_CONV:
		conv := func(x scope.Value) {
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Conv(x, n.Type())
			/*switch n.Type() {
			case object.INTEGER:
				switch v := x.(type) {
				case int8:
					rt2.DataOf(f.Parent())[n] = int32(x.(int8))
				case *big.Int:
					rt2.DataOf(f.Parent())[n] = int32(v.Int64())
				case int32:
					rt2.DataOf(f.Parent())[n] = v
				case
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			case object.SET:
				switch v := x.(type) {
				case int32:
					rt2.DataOf(f.Parent())[n] = big.NewInt(int64(v))
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			case object.REAL:
				switch v := x.(type) {
				case int32:
					rt2.DataOf(f.Parent())[n] = float64(v)
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			default:
				panic(fmt.Sprintln("wrong type", n.Type()))
			} */
		}
		switch n.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			x := sc.Select(n.Left().Object().Adr())
			assert.For(x != nil, 40)
			conv(scope.ValueFrom(x))
			return frame.End()
		case node.OperationNode:
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
			rt2.DataOf(f)[n.Left()] = n.Left().(node.ConstantNode).Data()
			return op, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//				rt2.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
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
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//				rt2.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
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
	case operation.ABS, operation.ODD, operation.CAP, operation.BITS:
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
		switch n.Right().(type) {
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
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//				rt2.DataOf(f)[n.Right()] = sc.Select(scope.Designator(n.Right()))
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		default:
			fmt.Println(reflect.TypeOf(n.Right()))
			panic("wrong right")
		}
	}

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f)
		switch l := n.Left().(type) {
		case node.ConstantNode:
			rt2.ValueOf(f)[n.Left().Adr()] = sc.Provide(n.Left())(nil)
			return right, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f)[n.Left().Adr()] = sc.Select(n.Left().Object().Adr())
				return right, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode, node.RangeNode:
			rt2.Push(rt2.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return right, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Select(l.Left().Object().Adr(), func(v scope.Value) {
					rt2.ValueOf(f)[n.Left().Adr()] = v.(scope.Record).Get(l.Object().Adr())
				})
				return right, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		default:
			fmt.Println(reflect.TypeOf(n.Left()))
			panic("wrong left")
		}
	}
	return left, frame.NOW
}
