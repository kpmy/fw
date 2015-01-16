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
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"reflect"
)

func incSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	op := node.New(constant.DYADIC, cp.SomeAdr()).(node.OperationNode)
	op.SetOperation(operation.PLUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Push(rt2.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(scope.Designator(n.Left()), func(interface{}) interface{} {
			return rt2.DataOf(f)[op]
		})
		return frame.End()
	}
	ret = frame.LATER
	return seq, ret
}

func decSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	op := node.New(constant.DYADIC, cp.SomeAdr()).(node.OperationNode)
	op.SetOperation(operation.MINUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Push(rt2.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(scope.Designator(n.Left()), func(interface{}) interface{} {
			return rt2.DataOf(f)[op]
		})
		return frame.End()
	}
	ret = frame.LATER
	return seq, ret
}

func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	a := fu.NodeOf(f)

	var leftId scope.ID
	var rightId scope.ID

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		switch a.Right().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(leftId, func(interface{}) interface{} {
					return a.Right().(node.ConstantNode).Data()
				})
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(leftId, func(interface{}) interface{} {
					return sc.Select(scope.Designator(a.Right()))
				})
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode, node.DerefNode:
			fu.Push(fu.New(a.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(leftId, func(interface{}) interface{} {
					return fu.DataOf(f)[a.Right()]
				})
				return frame.End()
			}
			ret = frame.LATER
		case node.IndexNode:
			rightId = scope.Designator(a.Right())
			rt2.Push(rt2.New(a.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rightId.Index = new(int64)
				*rightId.Index = int64(rt2.DataOf(f)[a.Right()].(int32))
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(leftId, func(interface{}) interface{} {
					return sc.Select(rightId)
				})
				return frame.End()
			}
			ret = frame.LATER
		case node.ProcedureNode:
			sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			sc.Update(leftId, func(interface{}) interface{} {
				return a.Right().Object()
			})
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
			leftId = scope.Designator(a.Left())
			seq, ret = right(f)
		case node.FieldNode:
			switch l.Left().(type) {
			case node.GuardNode:
				rt2.Push(rt2.New(l.Left()), f)
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					x := rt2.DataOf(f)[l.Left()].(node.Node)
					leftId = scope.Designator(a.Left(), x)
					fmt.Println(leftId)
					return right(f)
				}
				ret = frame.LATER
			case node.DerefNode:
				rt2.DataOf(f)[l.Left()] = scope.ID{}
				rt2.Push(rt2.New(l.Left()), f)
				seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					leftId = rt2.DataOf(f)[l.Left()].(scope.ID)
					leftId.Path = a.Left().Object().Name()
					return right(f)
				}
				ret = frame.LATER
			default:
				leftId = scope.Designator(a.Left())
				seq, ret = right(f)
			}
		case node.IndexNode:
			leftId = scope.Designator(a.Left())
			rt2.Push(rt2.New(a.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				leftId.Index = new(int64)
				*leftId.Index = int64(rt2.DataOf(f)[a.Left()].(int32))
				return right(f)
			}
			ret = frame.LATER
		case node.DerefNode:
			rt2.DataOf(f)[a.Left()] = scope.ID{}
			rt2.Push(rt2.New(a.Left()), f)
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				leftId = rt2.DataOf(f)[a.Left()].(scope.ID)
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
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		heap := scope.This(f.Domain().Discover(context.HEAP))
		if a.Right() != nil {
			seq, ret = This(expectExpr(f, a.Right(), func(...IN) OUT {
				fmt.Println("NEW", rt2.DataOf(f)[a.Right()], "here")
				sc.Update(scope.Designator(a.Left()), heap.Target().(scope.HeapAllocator).Allocate(a.Left(), rt2.DataOf(f)[a.Right()]))
				return End()
			}))
		} else {
			fmt.Println("NEW here")
			sc.Update(scope.Designator(a.Left()), heap.Target().(scope.HeapAllocator).Allocate(a.Left()))
			return frame.End()
		}
	default:
		panic(fmt.Sprintln("wrong statement", a.(node.AssignNode).Statement()))
	}
	return seq, ret
}
