package eval

import (
	"fw/cp/constant/enter"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
)

import (
	"ypk/halt"
)

func BeginExpression(in IN) (out OUT) {
	switch e := in.IR.(type) {
	case node.ConstantNode:
		out = Now(getConst)
	case node.DyadicNode:
		out = Now(getDop)
	case node.MonadicNode:
		out = Now(getMop)
	case node.Designator:
		out = Now(BeginDesignator)
	default:
		halt.As(100, reflect.TypeOf(e))
	}
	return
}

func BeginDesignator(in IN) (out OUT) {
	switch e := in.IR.(type) {
	case node.VariableNode:
		out = Now(getVar)
	case node.ParameterNode:
		out = Now(getVarPar)
	default:
		halt.As(100, reflect.TypeOf(e))
	}
	return
}

func GetExpression(in IN, key interface{}, expr node.Node, next Do) OUT {
	assert.For(expr != nil, 20)
	_, e_ok := expr.(node.Expression)
	_, d_ok := expr.(node.Designator)
	assert.For(e_ok || d_ok, 21, reflect.TypeOf(expr))
	assert.For(key != nil, 22)
	nf := rt2.New(expr)
	rt2.Push(nf, in.Frame)
	rt2.RegOf(in.Frame)[context.KEY] = key
	rt2.Assert(in.Frame, func(f frame.Frame) (bool, int) {
		v := rt2.RegOf(f)[key]
		return v != nil, 1961
	})
	return Later(func(IN) OUT {
		return Now(next)
	})
}

func GetDesignator(in IN, key interface{}, design node.Node, next Do) OUT {
	assert.For(design != nil, 20)
	_, ok := design.(node.Designator)
	assert.For(ok, 21)
	assert.For(key != nil, 22)
	nf := rt2.New(design)
	rt2.Push(nf, in.Frame)
	rt2.RegOf(in.Frame)[context.KEY] = key
	rt2.Assert(in.Frame, func(f frame.Frame) (bool, int) {
		v := rt2.RegOf(f)[key]
		return v != nil, 1957
	})
	return Later(func(IN) OUT {
		return Now(next)
	})
}

func BeginStatement(in IN) (out OUT) {
	switch n := in.IR.(type) {
	case node.EnterNode:
		var tail Do
		tail = func(in IN) OUT {
			q := in.Frame.Root().Queue()
			if q != nil {
				in.Frame.Root().PushFor(q, nil)
				return Later(tail)
			} else {
				return Now(doEnter)
			}
		}
		out = Now(tail)
	case node.AssignNode:
		out = Now(doAssign)
	case node.CallNode:
		out = Now(doCall)
	case node.ReturnNode:
		out = Now(doReturn)
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	return
}

func EndStatement(in IN) (out OUT) {
	switch n := in.IR.(type) {
	case node.EnterNode:
		out = Now(func(in IN) OUT {
			if n.Enter() == enter.PROCEDURE {
				rt2.ThisScope(in.Frame).Target().(scope.ScopeAllocator).Dispose(n)
			}
			if in.Parent != nil {
				par := rt2.RegOf(in.Parent)
				for k, v := range rt2.RegOf(in.Frame) {
					par[k] = v
				}
				val := rt2.ValueOf(in.Parent)
				for k, v := range rt2.ValueOf(in.Frame) {
					val[k] = v
				}
			}
			return End()
		})
	case node.AssignNode, node.CallNode:
		out = Now(func(in IN) OUT {
			next := n.Link()
			if next != nil {
				nf := rt2.New(next)
				if nf != nil {
					in.Frame.Root().PushFor(nf, in.Parent)
				}
			}
			if _, ok := n.(node.CallNode); ok {
				if in.Parent != nil {
					par := rt2.RegOf(in.Parent)
					for k, v := range rt2.RegOf(in.Frame) {
						par[k] = v
					}
					val := rt2.ValueOf(in.Parent)
					for k, v := range rt2.ValueOf(in.Frame) {
						val[k] = v
					}
				}
			}
			return End()
		})
	case node.ReturnNode: //do nothing
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	if out.Next != WRONG {
		return out.Do(in)
	}
	return End()
}
