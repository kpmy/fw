package rules

import (
	"fmt"
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
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

func sub(_a interface{}, _b interface{}) interface{} {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a - b
}

func eq(_a interface{}, _b interface{}) bool {
	assert.For(_a != nil, 20)
	assert.For(_b != nil, 21)
	var a int32 = int32Of(_a)
	var b int32 = int32Of(_b)
	return a == b
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
	var fu nodeframe.FrameUtils
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	n := fu.NodeOf(f).(node.MonadicNode)

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.NOT:
			fu.DataOf(f.Parent())[n] = not(fu.DataOf(f)[n.Left()])
			return frame.End()
		case operation.IS:
			x := sc.Select(scope.Designator(n.Left())).(object.Object)
			fu.DataOf(f.Parent())[n] = is(x, n.Object())
			return frame.End()
		case operation.ABS:
			fu.DataOf(f.Parent())[n] = abs(fu.DataOf(f)[n.Left()])
			return frame.End()
		case operation.ODD:
			fu.DataOf(f.Parent())[n] = odd(fu.DataOf(f)[n.Left()])
			return frame.End()
		case operation.CAP:
			fu.DataOf(f.Parent())[n] = cap_char(fu.DataOf(f)[n.Left()])
			return frame.End()
		case operation.BITS:
			fu.DataOf(f.Parent())[n] = bits(fu.DataOf(f)[n.Left()])
			return frame.End()
		default:
			panic("no such op")
		}

	}

	switch n.Operation() {
	case operation.CONVERT:
		switch n.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			x := sc.Select(scope.Designator(n.Left()))
			assert.For(x != nil, 40)
			switch n.Type() {
			case object.INTEGER:
				switch v := x.(type) {
				case int8:
					fu.DataOf(f.Parent())[n] = int32(x.(int8))
				case *big.Int:
					fu.DataOf(f.Parent())[n] = int32(v.Int64())
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			case object.SET:
				switch v := x.(type) {
				case int32:
					fu.DataOf(f.Parent())[n] = big.NewInt(int64(v))
				default:
					panic(fmt.Sprintln("ooops", reflect.TypeOf(x)))
				}
			default:
				panic(fmt.Sprintln("wrong type", n.Type()))
			}
			return frame.End()
		default:
			panic(fmt.Sprintln("unsupported left", reflect.TypeOf(n.Left())))
		}
	case operation.NOT:
		switch n.Left().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[n.Left()] = n.Left().(node.ConstantNode).Data()
			return op, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode:
			fu.Push(fu.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
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
		return expectExpr(f, n.Left(), op)
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
		case operation.MINUS:
			fu.DataOf(f.Parent())[n] = sub(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.EQUAL:
			fu.DataOf(f.Parent())[n] = eq(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.LESSER:
			fu.DataOf(f.Parent())[n] = lss(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.LESS_EQUAL:
			fu.DataOf(f.Parent())[n] = leq(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.LEN:
			fu.DataOf(f.Parent())[n] = length(n.Left().Object(), fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.NOT_EQUAL:
			fu.DataOf(f.Parent())[n] = neq(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
			return frame.End()
		case operation.GREATER:
			fu.DataOf(f.Parent())[n] = gtr(fu.DataOf(f)[n.Left()], fu.DataOf(f)[n.Right()])
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
			return op, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Right()] = sc.Select(scope.Designator(n.Right()))
				//fmt.Println(n.Right().Object(), reflect.TypeOf(n.Right().Object()))
				assert.For(fu.DataOf(f)[n.Right()] != nil, 60)
				return op, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode:
			fu.Push(fu.New(n.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return op, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Right()] = sc.Select(scope.Designator(n.Right()))
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
		n := fu.NodeOf(f)
		switch n.Left().(type) {
		case node.ConstantNode:
			fu.DataOf(f)[n.Left()] = n.Left().(node.ConstantNode).Data()
			return right, frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
				return right, frame.NOW
			}
			ret = frame.NOW
			return seq, ret
		case node.OperationNode, node.DerefNode:
			fu.Push(fu.New(n.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				return right, frame.NOW
			}
			ret = frame.LATER
			return seq, ret
		case node.FieldNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				fu.DataOf(f)[n.Left()] = sc.Select(scope.Designator(n.Left()))
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
