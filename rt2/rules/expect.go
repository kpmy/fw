package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	//	"fw/rt2/scope"
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
		/*id := scope.Designator(expr)
		rt2.Push(rt2.New(expr), parent)
		wait := func(...IN) OUT {
			if rt2.DataOf(parent)[expr] == nil {
				panic("no result")
			} else {
				id.Index = new(int64)
				*id.Index = int64(rt2.DataOf(parent)[expr].(int32))
				rt2.DataOf(parent)[expr] = sm.Select(id)
			}
			return OUT{do: next, next: NOW}
		}
		return OUT{do: wait, next: LATER}*/
		panic(0)
		return End()
	case node.ProcedureNode:
		rt2.RegOf(parent)[expr] = e.Object()
		return OUT{do: next, next: NOW}
	default:
		panic(fmt.Sprintln("not an expression", reflect.TypeOf(expr)))
	}
}
