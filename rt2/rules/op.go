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
	"math"
	"math/big"
	"reflect"
	"strings"
	"unicode/utf8"
	"ypk/assert"
)

func boolOf(x interface{}) (a bool) {
	switch x.(type) {
	case *bool:
		z := *x.(*bool)
		a = z
	case bool:
		a = x.(bool)
	default:
		panic(fmt.Sprintln("unsupported type", reflect.TypeOf(x)))
	}
	return a
}

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

func float64Of(x interface{}) (a float64) {
	//fmt.Println(reflect.TypeOf(x))
	switch v := x.(type) {
	case *int32:
		z := *x.(*int32)
		a = float64(z)
	case int32:
		a = float64(x.(int32))
	case float64:
		a = v
	default:
		panic(fmt.Sprintln("unsupported type", reflect.TypeOf(x)))
	}
	return a
}

func min(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return int32(math.Min(float64(a), float64(b)))
}

func max(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return int32(math.Max(float64(a), float64(b)))
}

func div(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a / b
}

func mod(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a % b
}

func times(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a * b
}

func slash(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a float64 = float64Of(_a)
	var b float64 = float64Of(_b)
	return a / b
}

func ash(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a << uint(b)
}

func in(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	fmt.Println("операция IN все врет")
	return a == b
}

func sub(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a - b
}

func and(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a bool = boolOf(_a)
	var b bool = boolOf(_b)
	return a && b
}

func or(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a bool = boolOf(_a)
	var b bool = boolOf(_b)
	return a || b
}

func lss(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a < b
}

func gtr(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a > b
}

func geq(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a >= b
}

func leq(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a <= b
}

func neq(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a != b
}

func not(_a interface{}) bool {
	assert.For(_a != nil, 20)
	var a bool = boolOf(_a)
	return !a
}

func is(p object.Object, typ object.ComplexType) bool {
	var compare func(x, a object.RecordType) bool
	compare = func(x, a object.RecordType) bool {
		switch {
		case x.Name() == a.Name():
			//	fmt.Println("eq")
			return true //опасно сравнивать имена конеш
		case x.BaseType() != nil:
			//	fmt.Println("go base")
			return compare(x.BaseType(), a)
		default:
			return false
		}
	}
	x, a := p.Complex().(object.RecordType)
	y, b := typ.(object.RecordType)
	//fmt.Println("compare", p.Complex(), typ, a, b, a && b && compare(x, y))
	return a && b && compare(x, y)
}

func length(a object.Object, _a, _b interface{}) (ret int64) {
	//assert.For(a != nil, 20)
	assert.For(_b != nil, 21)
	var b int32 = int32Of(_b)
	assert.For(b == 0, 22)
	if a != nil {
		assert.For(a.Type() == object.COMPLEX, 23)
		switch typ := a.Complex().(type) {
		case object.ArrayType:
			ret = typ.Len()
		case object.DynArrayType:
			switch _a.(type) {
			case string:
				ret = int64(utf8.RuneCountInString(_a.(string)))
			default:
				ret = 0
				fmt.Sprintln("unsupported", reflect.TypeOf(_a))
			}
		default:
			panic(fmt.Sprintln("unsupported", reflect.TypeOf(a.Complex())))
		}
	} else {
		switch _a.(type) {
		case string:
			ret = int64(utf8.RuneCountInString(_a.(string)))
		case []interface{}:
			ret = int64(len(_a.([]interface{})))
		default:
			panic(fmt.Sprintln("unsupported", reflect.TypeOf(_a)))
		}
	}
	return ret
}

func abs(_a interface{}) interface{} {
	assert.For(_a != nil, 20)
	var a int32 = int32Of(_a)
	return int32(math.Abs(float64(a)))
}

func odd(_a interface{}) bool {
	assert.For(_a != nil, 20)
	var a int32 = int32Of(_a)
	return int32(math.Abs(float64(a)))%2 == 1
}

func cap_char(_a interface{}) interface{} {
	assert.For(_a != nil, 20)
	var a int32 = int32Of(_a)
	x := []rune{rune(a), rune(0)}
	return int32([]rune(strings.ToUpper(string(x)))[0])
}

func bits(_a interface{}) interface{} {
	assert.For(_a != nil, 20)
	return big.NewInt(int64(int32Of(_a)))
}

func mopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	n := rt2.NodeOf(f).(node.MonadicNode)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := rt2.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.NOT:
			rt2.DataOf(f.Parent())[n] = not(rt2.DataOf(f)[n.Left()])
			return frame.End()
		case operation.IS:
			/*	x := sc.Select(scope.Designator(n.Left())).(object.Object)
				rt2.DataOf(f.Parent())[n] = is(x, n.Object())*/
			return frame.End()
		case operation.ABS:
			rt2.DataOf(f.Parent())[n] = abs(rt2.DataOf(f)[n.Left()])
			return frame.End()
		case operation.ODD:
			rt2.DataOf(f.Parent())[n] = odd(rt2.DataOf(f)[n.Left()])
			return frame.End()
		case operation.CAP:
			rt2.DataOf(f.Parent())[n] = cap_char(rt2.DataOf(f)[n.Left()])
			return frame.End()
		case operation.BITS:
			rt2.DataOf(f.Parent())[n] = bits(rt2.DataOf(f)[n.Left()])
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
			rt2.DataOf(f.Parent())[n] = length(n.Left().Object(), rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.NOT_EQUAL:
			rt2.DataOf(f.Parent())[n] = neq(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.GREATER:
			rt2.DataOf(f.Parent())[n] = gtr(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.MAX:
			rt2.DataOf(f.Parent())[n] = max(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.MIN:
			rt2.DataOf(f.Parent())[n] = min(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.DIV:
			rt2.DataOf(f.Parent())[n] = div(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.MOD:
			rt2.DataOf(f.Parent())[n] = mod(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.TIMES:
			rt2.DataOf(f.Parent())[n] = times(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.SLASH:
			rt2.DataOf(f.Parent())[n] = slash(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.IN:
			rt2.DataOf(f.Parent())[n] = in(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.ASH:
			rt2.DataOf(f.Parent())[n] = ash(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.AND:
			rt2.DataOf(f.Parent())[n] = and(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.OR:
			rt2.DataOf(f.Parent())[n] = or(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
			return frame.End()
		case operation.GREAT_EQUAL:
			rt2.DataOf(f.Parent())[n] = geq(rt2.DataOf(f)[n.Left()], rt2.DataOf(f)[n.Right()])
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
		switch n.Left().(type) {
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
				//				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				//				rt2.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
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
