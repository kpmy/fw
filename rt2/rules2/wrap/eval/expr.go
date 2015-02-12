package eval

import (
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/scope"
	"ypk/assert"
	"ypk/halt"
)

func getConst(in IN) OUT {
	c := in.IR.(node.ConstantNode)
	sc := rt2.ThisScope(in.Frame)
	fn := sc.Provide(c)
	assert.For(fn != nil, 40)
	rt2.ValueOf(in.Parent)[c.Adr()] = fn(nil)
	rt2.RegOf(in.Parent)[in.Key] = c.Adr()
	return End()
}

func getVar(in IN) OUT {
	v := in.IR.(node.VariableNode)
	rt2.ScopeFor(in.Frame, v.Object().Adr(), func(val scope.Value) {
		rt2.ValueOf(in.Parent)[v.Adr()] = val
		rt2.RegOf(in.Parent)[in.Key] = v.Adr()
	})
	return End()
}

func getVarPar(in IN) OUT {
	v := in.IR.(node.ParameterNode)
	rt2.ScopeFor(in.Frame, v.Object().Adr(), func(val scope.Value) {
		rt2.ValueOf(in.Parent)[v.Adr()] = val
		rt2.RegOf(in.Parent)[in.Key] = v.Adr()
	})
	return End()
}

func getDop(in IN) OUT {
	const (
		left  = "dop:left"
		right = "dop:right"
	)
	op := in.IR.(node.DyadicNode)

	do := func(in IN) OUT {
		var (
			res scope.Value
		)
		l := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		r := rt2.ValueOf(in.Frame)[KeyOf(in, right)]
		switch op.Operation() {
		case operation.PLUS:
			res = scope.Ops.Sum(l, r)
		case operation.MINUS:
			res = scope.Ops.Sub(l, r)
		case operation.EQUAL:
			res = scope.Ops.Eq(l, r)
		case operation.LESSER:
			res = scope.Ops.Lss(l, r)
		case operation.LESS_EQUAL:
			res = scope.Ops.Leq(l, r)
		case operation.LEN:
			res = scope.Ops.Len(op.Left().Object(), l, r)
		case operation.NOT_EQUAL:
			res = scope.Ops.Neq(l, r)
		case operation.GREATER:
			res = scope.Ops.Gtr(l, r)
		case operation.MAX:
			res = scope.Ops.Max(l, r)
		case operation.MIN:
			res = scope.Ops.Min(l, r)
		case operation.DIV:
			res = scope.Ops.Div(l, r)
		case operation.MOD:
			res = scope.Ops.Mod(l, r)
		case operation.ALIEN_MSK:
			res = scope.Ops.Msk(l, r)
		case operation.TIMES:
			res = scope.Ops.Mult(l, r)
		case operation.SLASH:
			res = scope.Ops.Divide(l, r)
		case operation.IN:
			res = scope.Ops.In(l, r)
		case operation.ASH:
			res = scope.Ops.Ash(l, r)
		case operation.AND:
			res = scope.Ops.And(l, r)
		case operation.OR:
			res = scope.Ops.Or(l, r)
		case operation.GREAT_EQUAL:
			res = scope.Ops.Geq(l, r)
		default:
			halt.As(100, op.Operation())
		}
		assert.For(res != nil, 40)
		rt2.ValueOf(in.Parent)[op.Adr()] = res
		rt2.RegOf(in.Parent)[in.Key] = op.Adr()
		return End()
	}

	next := func(IN) OUT {
		return GetExpression(in, right, op.Right(), do)
	}

	short := func(IN) OUT {
		id := KeyOf(in, left)
		lv := rt2.ValueOf(in.Frame)[id]
		switch op.Operation() {
		case operation.AND:
			val := scope.GoTypeFrom(lv).(bool)
			if val {
				return Now(next)
			} else {
				rt2.ValueOf(in.Parent)[op.Adr()] = scope.TypeFromGo(false)
				rt2.RegOf(in.Parent)[in.Key] = op.Adr()
				return End()
			}
		case operation.OR:
			val := scope.GoTypeFrom(lv).(bool)
			if !val {
				return Now(next)
			} else {
				rt2.ValueOf(in.Parent)[op.Adr()] = scope.TypeFromGo(true)
				rt2.RegOf(in.Parent)[in.Key] = op.Adr()
				return End()
			}
		default:
			return Now(next)
		}
	}
	return GetExpression(in, left, op.Left(), short)
}

func getMop(in IN) OUT {
	const left = "mop:left"
	op := in.IR.(node.MonadicNode)

	do := func(in IN) OUT {
		lv := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		var res scope.Value
		switch op.Operation() {
		case operation.ALIEN_CONV:
			if op.Type() != object.NOTYPE {
				res = scope.Ops.Conv(lv, op.Type())
			} else {
				res = scope.Ops.Conv(lv, op.Type(), op.Complex())
			}

		default:
			halt.As(100, "unknown op", op.Operation())
		}
		assert.For(res != nil, 60)
		rt2.ValueOf(in.Parent)[op.Adr()] = res
		rt2.RegOf(in.Parent)[in.Key] = op.Adr()
		return End()
	}

	return GetExpression(in, left, op.Left(), do)
}
