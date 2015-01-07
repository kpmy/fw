package rules

import (
	"fmt"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
)

func guardSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	sc := scope.This(f.Domain().Discover(context.SCOPE))
	n := rt2.Utils.NodeOf(f).(node.GuardNode)
	var obj object.Object
	switch l := n.Left().(type) {
	case node.VariableNode:
		obj = l.Object()
	case node.ParameterNode:
		obj = sc.Select(scope.Designator(l)).(object.Object)
	default:
		panic(fmt.Sprintln("unsupported left", reflect.TypeOf(l)))
	}
	if is(obj, n.Type()) {
		rt2.Utils.DataOf(f.Parent())[n] = n.Left()
		return frame.End()
	} else {
		trap := node.New(constant.TRAP).(node.TrapNode)
		code := node.New(constant.CONSTANT).(node.ConstantNode)
		code.SetData(0)
		trap.SetLeft(code)
		rt2.Utils.Push(rt2.Utils.New(trap), f)
		return frame.Tail(frame.STOP), frame.LATER
	}
}
