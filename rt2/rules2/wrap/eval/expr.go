package eval

import (
	"fmt"
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/scope"
	"reflect"
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

func getField(in IN) OUT {
	const left = "field:left"
	f := in.IR.(node.FieldNode)
	return GetDesignator(in, left, f.Left(), func(in IN) OUT {
		_v := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		switch v := _v.(type) {
		case scope.Record:
			fld := v.Get(f.Object().Adr())
			rt2.ValueOf(in.Parent)[f.Adr()] = fld
			rt2.RegOf(in.Parent)[in.Key] = f.Adr()
			return End()
		default:
			halt.As(100, reflect.TypeOf(v))
		}
		panic(0)
	})
}

func getIndex(in IN) OUT {
	const (
		left  = "index:left"
		right = "index:right"
	)
	i := in.IR.(node.IndexNode)
	return GetExpression(in, right, i.Right(), func(IN) OUT {
		idx := rt2.ValueOf(in.Frame)[KeyOf(in, right)]
		assert.For(idx != nil, 40)
		return GetDesignator(in, left, i.Left(), func(IN) OUT {
			arr := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
			assert.For(arr != nil, 41)
			switch a := arr.(type) {
			case scope.Array:
				rt2.ValueOf(in.Parent)[i.Adr()] = a.Get(idx)
				rt2.RegOf(in.Parent)[in.Key] = i.Adr()
				return End()
			default:
				halt.As(100, reflect.TypeOf(a))
			}
			panic(890)
		})
	})
}

func getProc(in IN) OUT {
	p := in.IR.(node.ProcedureNode)
	sc := rt2.ThisScope(in.Frame)
	fn := sc.Provide(p.Object())
	assert.For(fn != nil, 40)
	rt2.ValueOf(in.Parent)[p.Adr()] = fn(nil)
	rt2.RegOf(in.Parent)[in.Key] = p.Adr()
	return End()
}

func getDeref(in IN) OUT {
	const left = "design:left"
	d := in.IR.(node.DerefNode)
	return GetDesignator(in, left, d.Left(), func(in IN) OUT {
		_v := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		switch v := _v.(type) {
		case scope.Array:
			assert.For(!d.Ptr(), 40)
			t, c := scope.Ops.TypeOf(v)
			switch cc := c.(type) {
			case object.ArrayType:
				assert.For(cc.Base() == object.CHAR || cc.Base() == object.SHORTCHAR, 41)
				rt2.ValueOf(in.Parent)[d.Adr()] = scope.TypeFromGo(scope.GoTypeFrom(v))
				rt2.RegOf(in.Parent)[in.Key] = d.Adr()
				return End()
			case object.DynArrayType:
				assert.For(cc.Base() == object.CHAR || cc.Base() == object.SHORTCHAR, 41)
				rt2.ValueOf(in.Parent)[d.Adr()] = scope.TypeFromGo(scope.GoTypeFrom(v))
				rt2.RegOf(in.Parent)[in.Key] = d.Adr()
				return End()
			default:
				halt.As(100, t, reflect.TypeOf(cc))
			}
		default:
			halt.As(100, reflect.TypeOf(v))
		}
		panic(0)
	})
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
			fmt.Println(res, l, r)
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
		case operation.NOT:
			res = scope.Ops.Not(lv)
		case operation.IS:
			/*sc := rt2.ScopeFor(f, n.Left().Object().Adr())
			rt2.ValueOf(f.Parent())[n.Adr()] = scope.Ops.Is(sc.Select(n.Left().Object().Adr()), n.Object().Complex())
			return frame.End()*/
			panic(0)
		case operation.ABS:
			res = scope.Ops.Abs(lv)
		case operation.ODD:
			res = scope.Ops.Odd(lv)
		case operation.CAP:
			res = scope.Ops.Cap(lv)
		case operation.BITS:
			res = scope.Ops.Bits(lv)
		case operation.MINUS:
			res = scope.Ops.Minus(lv)
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
