package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
)

func guardSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f).(node.GuardNode)
	var obj scope.Value
	switch l := n.Left().(type) {
	case node.VariableNode, node.ParameterNode:
		sc := rt2.ScopeFor(f, l.Object().Adr())
		obj = sc.Select(l.Object().Adr())
	default:
		panic(fmt.Sprintln("unsupported left", reflect.TypeOf(l)))
	}
	if scope.GoTypeFrom(scope.Ops.Is(obj, n.Type())).(bool) {
		rt2.RegOf(f.Parent())[n] = n.Left()
		return frame.End()
	} else {
		return doTrap(f, 0)
	}
}
