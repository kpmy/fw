package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
)

//функция вернет в данные родительского фрейма вычисленное значение expr
func expectExpr(parent frame.Frame, expr node.Node, next Do) OUT {
	assert.For(expr != nil, 20)
	sm := rt2.ScopeOf(parent)
	switch e := expr.(type) {
	case node.ConstantNode:
		rt2.ValueOf(parent)[expr.Adr()] = sm.Provide(e)(nil)
		return OUT{do: next, next: NOW}
	case node.VariableNode, node.ParameterNode:
		rt2.ValueOf(parent)[expr.Adr()] = sm.Select(expr.Object().Adr())
		return OUT{do: next, next: NOW}
	case node.OperationNode, node.CallNode, node.DerefNode:
		rt2.Push(rt2.New(expr), parent)
		wait := func(...IN) OUT {
			if rt2.RegOf(parent)[expr] == nil && rt2.ValueOf(parent)[expr.Adr()] == nil {
				panic("no result")
			}
			return OUT{do: next, next: NOW}
		}
		return OUT{do: wait, next: LATER}
	case node.IndexNode:
		rt2.Push(rt2.New(e), parent)
		rt2.Assert(parent, func(f frame.Frame) (bool, int) {
			return rt2.ValueOf(f)[e.Adr()] != nil, 64
		})
		wait := func(...IN) OUT {
			idx := rt2.ValueOf(parent)[e.Adr()]
			arr := sm.Select(e.Left().Object().Adr()).(scope.Array)
			idx = arr.Get(idx)
			rt2.ValueOf(parent)[e.Adr()] = idx
			return OUT{do: next, next: NOW}
		}
		return OUT{do: wait, next: LATER}
	case node.ProcedureNode:
		rt2.RegOf(parent)[expr] = e.Object()
		return OUT{do: next, next: NOW}
	default:
		panic(fmt.Sprintln("not an expression", reflect.TypeOf(expr)))
	}
}
