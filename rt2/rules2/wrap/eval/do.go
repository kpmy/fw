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
	case node.ProcedureNode:
		out = Now(getProc)
	case node.DerefNode:
		out = Now(getDeref)
	case node.FieldNode:
		out = Now(getField)
	case node.IndexNode:
		out = Now(getIndex)
	default:
		halt.As(100, reflect.TypeOf(e))
	}
	return
}

func GetStrange(in IN, key interface{}, ss node.Node, next Do) (out OUT) {
	assert.For(ss != nil, 20)
	assert.For(key != nil, 21)
	switch ss.(type) {
	case node.IfNode:
		nf := rt2.New(ss)
		rt2.Push(nf, in.Frame)
		rt2.RegOf(in.Frame)[context.KEY] = key
		rt2.Assert(in.Frame, func(f frame.Frame) (bool, int) {
			v := rt2.RegOf(f)[key]
			return v != nil, 1980
		})
		return Later(func(IN) OUT {
			return Now(next)
		})
	default:
		halt.As(100, reflect.TypeOf(ss))
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
func BeginStrange(in IN) OUT {
	switch s := in.IR.(type) {
	case node.IfNode:
		return Now(doIf)
	default:
		halt.As(100, reflect.TypeOf(s))
	}
	panic(0)
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
	case node.ConditionalNode:
		out = Now(doCondition)
	case node.WhileNode:
		out = Now(doWhile)
	case node.RepeatNode:
		out = Now(doRepeat)
	case node.LoopNode:
		out = Now(doLoop)
	case node.ExitNode:
		out = Now(doExit)
	case node.InitNode:
		out = Later(Tail(STOP))
	case node.TrapNode:
		out = Now(doTrap)
	case node.WithNode:
		out = Now(doWith)
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
	case node.CallNode:
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
	case node.AssignNode, node.ConditionalNode, node.WhileNode, node.RepeatNode, node.ExitNode, node.InitNode, node.WithNode:
		out = Now(func(in IN) OUT {
			next := n.Link()
			if next != nil {
				nf := rt2.New(next)
				if nf != nil {
					in.Frame.Root().PushFor(nf, in.Parent)
				}
			}
			return End()
		})
	case node.ReturnNode, node.LoopNode: //do nothing
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	if out.Next != WRONG {
		return out.Do(in)
	}
	return End()
}
