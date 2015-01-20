package rules

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
)

func incSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	op := node.New(constant.DYADIC, int(cp.SomeAdr())).(node.OperationNode)
	op.SetOperation(operation.PLUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Push(rt2.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(n.Left().Object().Adr(), scope.Simple(rt2.ValueOf(f)[op.Adr()]))
		return frame.End()
	}
	ret = frame.LATER
	return seq, ret
}

func decSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	op := node.New(constant.DYADIC, int(cp.SomeAdr())).(node.OperationNode)
	op.SetOperation(operation.MINUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Push(rt2.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(n.Left().Object().Adr(), scope.Simple(rt2.ValueOf(f)[op.Adr()]))
		return frame.End()
	}
	ret = frame.LATER
	return seq, ret
}

func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	a := rt2.NodeOf(f)

	var left scope.Value
	var rightId cp.ID

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		vleft := left.(scope.Variable)
		switch l := a.Right().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				vleft.Set(sc.Provide(a.Right())(nil)) //scope.Simple(a.Right().(node.ConstantNode).Data()))
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				vleft.Set(sc.Select(a.Right().Object().Adr()))
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode, node.DerefNode:
			rt2.Push(rt2.New(a.Right()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[a.Right().Adr()] != nil, 61
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				vleft.Set(rt2.ValueOf(f)[a.Right().Adr()])
				return frame.End()
			}
			ret = frame.LATER
		case node.IndexNode:
			rightId = l.Adr()
			rt2.Push(rt2.New(l), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[rightId] != nil, 62
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				right := rt2.ValueOf(f)[l.Adr()]
				arr := sc.Select(l.Left().Object().Adr()).(scope.Array)
				right = arr.Get(right)
				vleft.Set(right)
				return frame.End()
			}
			ret = frame.LATER
		case node.ProcedureNode:
			sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			vleft.Set(sc.Provide(a.Right().Object())(nil))
			return frame.End()
		default:
			fmt.Println(reflect.TypeOf(a.Right()))
			panic("wrong right")
		}
		return seq, ret
	}

	switch a.(node.AssignNode).Statement() {
	case statement.ASSIGN:
		switch l := a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			left = sc.Select(a.Left().Object().Adr())
			seq, ret = right(f)
		case node.FieldNode:
			rt2.Push(rt2.New(l), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 62
			})
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				left = rt2.ValueOf(f)[l.Adr()]
				return right(f)
			}
			ret = frame.LATER
		case node.IndexNode:
			rt2.Push(rt2.New(a.Left()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 62
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				left = rt2.ValueOf(f)[l.Adr()]
				arr := sc.Select(l.Left().Object().Adr()).(scope.Array)
				left = arr.Get(left)
				return right(f)
			}
			ret = frame.LATER
		case node.DerefNode:
			//			rt2.DataOf(f)[a.Left()] = scope.ID{}
			rt2.Push(rt2.New(a.Left()), f)
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				//				leftId = rt2.DataOf(f)[a.Left()].(scope.ID)
				return right(f)
			}
			ret = frame.LATER
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong left")
		}
	case statement.INC, statement.INCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			seq, ret = incSeq(f)
		default:
			panic(fmt.Sprintln("wrong left", reflect.TypeOf(a.Left())))
		}
	case statement.DEC, statement.EXCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			seq, ret = decSeq(f)
		default:
			panic(fmt.Sprintln("wrong left", reflect.TypeOf(a.Left())))
		}
	case statement.NEW:
		//		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		//		heap := scope.This(f.Domain().Discover(context.HEAP))
		if a.Right() != nil {
			seq, ret = This(expectExpr(f, a.Right(), func(...IN) OUT {
				fmt.Println("NEW", rt2.DataOf(f)[a.Right()], "here")
				//				sc.Update(scope.Designator(a.Left()), heap.Target().(scope.HeapAllocator).Allocate(a.Left(), rt2.DataOf(f)[a.Right()]))
				return End()
			}))
		} else {
			fmt.Println("NEW here")
			//			sc.Update(scope.Designator(a.Left()), heap.Target().(scope.HeapAllocator).Allocate(a.Left()))
			return frame.End()
		}
	default:
		panic(fmt.Sprintln("wrong statement", a.(node.AssignNode).Statement()))
	}
	return seq, ret
}
