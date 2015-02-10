package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

//функция вернет в данные родительского фрейма вычисленное значение expr
func expectExpr(parent frame.Frame, expr node.Node, next Do) OUT {
	assert.For(expr != nil, 20)
	switch e := expr.(type) {
	case node.ConstantNode:
		rt2.ValueOf(parent)[expr.Adr()] = rt2.ThisScope(parent).Provide(e)(nil)
		return OUT{do: next, next: NOW}
	case node.VariableNode, node.ParameterNode:
		m := rtm.ModuleOfObject(parent.Domain(), expr.Object())
		assert.For(m != nil, 40)
		imp := m.ImportOf(expr.Object())
		if imp != "" {
			//			md := rtm.ModuleDomain(parent.Domain(), imp)
			//sm = md.Discover(context.SCOPE).(scope.Manager)
			fm := rtm.Module(parent.Domain(), imp)
			ol := fm.ObjectByName(fm.Enter, expr.Object().Name())
			for _, obj := range ol {
				if _, ok := obj.(object.VariableObject); ok {
					sm := rt2.ScopeFor(parent, obj.Adr())
					rt2.ValueOf(parent)[expr.Adr()] = sm.Select(obj.Adr())
				}
			}
		} else {
			sm := rt2.ScopeFor(parent, expr.Object().Adr())
			rt2.ValueOf(parent)[expr.Adr()] = sm.Select(expr.Object().Adr())
		}
		return OUT{do: next, next: NOW}
	case node.OperationNode, node.CallNode, node.DerefNode, node.FieldNode:
		rt2.Push(rt2.New(expr), parent)
		wait := func(...IN) OUT {
			if rt2.RegOf(parent)[expr] == nil && rt2.ValueOf(parent)[expr.Adr()] == nil {
				raw := rt2.RegOf(parent)["RETURN"]
				if val, ok := raw.(scope.Value); !ok {
					halt.As(100, "no result from ", expr.Adr(), raw)
				} else {
					rt2.ValueOf(parent)[expr.Adr()] = val
				}
			}
			return OUT{do: next, next: NOW}
		}
		return OUT{do: wait, next: LATER}
	case node.IndexNode:
		rt2.Push(rt2.New(e), parent)
		rt2.Assert(parent, func(f frame.Frame) (bool, int) {
			return rt2.ValueOf(f)[e.Adr()] != nil, 64
		})
		wait := func(in ...IN) OUT {
			idx := rt2.ValueOf(parent)[e.Adr()]
			return expectExpr(in[0].frame, e.Left(), func(...IN) OUT {
				arr := rt2.ValueOf(in[0].frame)[e.Left().Adr()].(scope.Array)
				idx = arr.Get(idx)
				rt2.ValueOf(parent)[e.Adr()] = idx
				return OUT{do: next, next: NOW}
			})
		}
		return OUT{do: wait, next: LATER}
	case node.ProcedureNode:
		rt2.RegOf(parent)[expr] = e.Object()
		return OUT{do: next, next: NOW}
	default:
		panic(fmt.Sprintln("not an expression", reflect.TypeOf(expr)))
	}
}
