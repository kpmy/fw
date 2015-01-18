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
	//	sm := rt2.ScopeOf(parent)
	switch e := expr.(type) {
	case node.ConstantNode:
		rt2.DataOf(parent)[expr] = e.Data()
		return OUT{do: next, next: NOW}
	case node.VariableNode, node.ParameterNode:
		//		rt2.DataOf(parent)[expr] = sm.Select(scope.Designator(expr))
		return OUT{do: next, next: NOW}
	case node.OperationNode, node.CallNode:
		rt2.Push(rt2.New(expr), parent)
		wait := func(...IN) OUT {
			if rt2.DataOf(parent)[expr] == nil {
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
		return End()
	case node.ProcedureNode:
		rt2.DataOf(parent)[expr] = e.Object()
		return OUT{do: next, next: NOW}
	default:
		panic(fmt.Sprintln("not an expression", reflect.TypeOf(expr)))
	}
}
