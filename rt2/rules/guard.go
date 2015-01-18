package rules

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	//	"fw/rt2/context"
	"fw/rt2/frame"
	//	"fw/rt2/scope"
	"reflect"
)

func guardSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	//	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	n := rt2.NodeOf(f).(node.GuardNode)
	var obj object.Object
	switch l := n.Left().(type) {
	case node.VariableNode:
		obj = l.Object()
	case node.ParameterNode:
		//		obj = sc.Select(scope.Designator(l)).(object.Object)
	default:
		panic(fmt.Sprintln("unsupported left", reflect.TypeOf(l)))
	}
	if is(obj, n.Type()) {
		rt2.DataOf(f.Parent())[n] = n.Left()
		return frame.End()
	} else {
		trap := node.New(constant.TRAP, int(cp.SomeAdr())).(node.TrapNode)
		code := node.New(constant.CONSTANT, int(cp.SomeAdr())).(node.ConstantNode)
		code.SetData(0)
		trap.SetLeft(code)
		rt2.Push(rt2.New(trap), f)
		return frame.Tail(frame.STOP), frame.LATER
	}
}
