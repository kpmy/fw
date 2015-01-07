package rules

import (
	"fmt"
	"fw/cp/constant"
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/statement"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"reflect"
)

func incSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.Utils.NodeOf(f)
	op := node.New(constant.DYADIC).(node.OperationNode)
	op.SetOperation(operation.PLUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Utils.Push(rt2.Utils.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(scope.Designator(n.Left()), func(interface{}) interface{} {
			return rt2.Utils.DataOf(f)[op]
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
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				sc.Update(leftId, func(interface{}) interface{} {
					return sc.Select(scope.Designator(a.Right()))
				})
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode:
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
			rt2.Utils.Push(rt2.Utils.New(a.Right()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rightId.Index = new(int64)
				*rightId.Index = int64(rt2.Utils.DataOf(f)[a.Right()].(int32))
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
				rt2.Utils.Push(rt2.Utils.New(l.Left()), f)
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					x := rt2.Utils.DataOf(f)[l.Left()].(node.Node)
					leftId = scope.Designator(a.Left(), x)
					fmt.Println(leftId)
					return right(f)
				}
				ret = frame.LATER
			default:
				leftId = scope.Designator(a.Left())
				seq, ret = right(f)
			}
		case node.IndexNode:
			leftId = scope.Designator(a.Left())
			rt2.Utils.Push(rt2.Utils.New(a.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				leftId.Index = new(int64)
				*leftId.Index = int64(rt2.Utils.DataOf(f)[a.Left()].(int32))
				return right(f)
			}
			ret = frame.LATER
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong left")
		}
	case statement.INC:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			seq, ret = incSeq(f)
		default:
			panic(fmt.Sprintln("wrong left", reflect.TypeOf(a.Left())))
		}
	default:
		panic(fmt.Sprintln("wrong statement", a.(node.AssignNode).Statement()))
	}
	return seq, ret
}
