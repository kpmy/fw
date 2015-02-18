package eval

import (
	"fmt"
	"fw/cp/constant/operation"
	"fw/cp/node"
	"fw/cp/object"
	"fw/cp/traps"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/scope"
	"math/big"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

func getConst(in IN) OUT {
	c := in.IR.(node.ConstantNode)
	sc := rt2.ModScope(in.Frame)
	fn := sc.Provide(c)
	assert.For(fn != nil, 40)
	rt2.ValueOf(in.Parent)[c.Adr()] = fn
	rt2.RegOf(in.Parent)[in.Key] = c.Adr()
	return End()
}

func getVar(in IN) OUT {
	v := in.IR.(node.VariableNode)
	sc := rt2.ScopeFor(in.Frame, v.Object().Adr(), func(val scope.Value) {
		rt2.ValueOf(in.Parent)[v.Adr()] = val
		rt2.RegOf(in.Parent)[in.Key] = v.Adr()
	})
	rt2.RegOf(in.Parent)[context.META] = &Meta{Scope: sc, Id: v.Object().Adr()}
	return End()
}

func getVarPar(in IN) OUT {
	v := in.IR.(node.ParameterNode)
	sc := rt2.ScopeFor(in.Frame, v.Object().Adr(), func(val scope.Value) {
		rt2.ValueOf(in.Parent)[v.Adr()] = val
		rt2.RegOf(in.Parent)[in.Key] = v.Adr()
		rt2.RegOf(in.Parent)[context.META] = &Meta{}
	})
	rt2.RegOf(in.Parent)[context.META] = &Meta{Scope: sc, Id: v.Object().Adr()}
	return End()
}

func getField(in IN) OUT {
	const left = "field:left"
	f := in.IR.(node.FieldNode)
	return GetDesignator(in, left, f.Left(), func(in IN) OUT {
		_v := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		switch v := _v.(type) {
		case scope.Record:
			fld := v.Get(f.Object().Adr()).(scope.Variable)
			rt2.ValueOf(in.Parent)[f.Adr()] = fld
			rt2.RegOf(in.Parent)[in.Key] = f.Adr()
			rt2.RegOf(in.Parent)[context.META] = &Meta{Scope: nil, Rec: v, Id: fld.Id()}
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
				rt2.RegOf(in.Parent)[context.META] = &Meta{Arr: a, Id: a.Id()}
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
	sc := rt2.ModScope(in.Frame)
	fn := sc.Provide(p.Object())
	assert.For(fn != nil, 40)
	rt2.ValueOf(in.Parent)[p.Adr()] = fn
	rt2.RegOf(in.Parent)[in.Key] = p.Adr()
	rt2.RegOf(in.Parent)[context.META] = &Meta{}
	return End()
}

func getDeref(in IN) OUT {
	const left = "design:left"
	d := in.IR.(node.DerefNode)
	return GetDesignator(in, left, d.Left(), func(in IN) (out OUT) {
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
				rt2.RegOf(in.Parent)[context.META] = &Meta{}
				return End()
			case object.DynArrayType:
				assert.For(cc.Base() == object.CHAR || cc.Base() == object.SHORTCHAR, 41)
				rt2.ValueOf(in.Parent)[d.Adr()] = scope.TypeFromGo(scope.GoTypeFrom(v))
				rt2.RegOf(in.Parent)[in.Key] = d.Adr()
				rt2.RegOf(in.Parent)[context.META] = &Meta{}
				return End()
			default:
				halt.As(100, t, reflect.TypeOf(cc))
			}
		case scope.Pointer:
			switch {
			case scope.GoTypeFrom(v.Get()) == nil:
				out = makeTrap(in.Frame, traps.NILderef)
			case d.Ptr():
				out = End()
				switch r := v.Get().(type) {
				case scope.Record:
					rec := r.(scope.Record)
					rt2.ValueOf(in.Parent)[d.Adr()] = rec
					rt2.RegOf(in.Parent)[in.Key] = d.Adr()
					rt2.RegOf(in.Parent)[context.META] = &Meta{Scope: rt2.ScopeFor(in.Frame, rec.Id()), Id: rec.Id()}
				case scope.Array:
					arr := r.(scope.Array)
					rt2.ValueOf(in.Parent)[d.Adr()] = arr
					rt2.RegOf(in.Parent)[in.Key] = d.Adr()
					rt2.RegOf(in.Parent)[context.META] = &Meta{Scope: rt2.ScopeFor(in.Frame, arr.Id()), Id: arr.Id()}
				default:
					halt.As(100, reflect.TypeOf(r))
				}
			case !d.Ptr():
				out = End()
				switch r := v.Get().(type) {
				case scope.Array:
					arr := r.(scope.Array)
					rt2.ValueOf(in.Parent)[d.Adr()] = scope.TypeFromGo(scope.GoTypeFrom(arr))
					rt2.RegOf(in.Parent)[in.Key] = d.Adr()
					rt2.RegOf(in.Parent)[context.META] = &Meta{Arr: arr, Id: arr.Id()}
				default:
					halt.As(100, reflect.TypeOf(r))
				}
			default:
				halt.As(100, d.Adr(), d.Ptr(), v, v.Get())
			}
			return out
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
			res = scope.Ops.Is(lv, op.Object().Complex())
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

func getGuard(in IN) OUT {
	const left = "guard:left"
	g := in.IR.(node.GuardNode)
	return GetDesignator(in, left, g.Left(), func(IN) OUT {
		v := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		assert.For(v != nil, 20)
		if scope.GoTypeFrom(scope.Ops.Is(v, g.Type())).(bool) {
			rt2.ValueOf(in.Parent)[g.Adr()] = v
			rt2.RegOf(in.Parent)[in.Key] = g.Adr()
			rt2.RegOf(in.Parent)[context.META] = rt2.RegOf(in.Frame)[context.META] //&Meta{Id: vv.Id(), }
			return End()
		} else {
			return makeTrap(in.Frame, 0)
		}

	})
}

func bit_range(_f scope.Value, _t scope.Value) scope.Value {
	f := scope.GoTypeFrom(_f).(int32)
	t := scope.GoTypeFrom(_t).(int32)
	ret := big.NewInt(0)
	for i := f; i <= t; i++ {
		ret = ret.SetBit(ret, int(i), 1)
	}
	fmt.Println("bits", ret)
	return scope.TypeFromGo(ret)
}

func getRange(in IN) OUT {
	const (
		left  = "range:left"
		right = "range:right"
	)
	r := in.IR.(node.RangeNode)
	return GetExpression(in, left, r.Left(), func(IN) OUT {
		return GetExpression(in, right, r.Right(), func(in IN) OUT {
			lv := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
			rv := rt2.ValueOf(in.Frame)[KeyOf(in, right)]
			res := bit_range(lv, rv)
			rt2.ValueOf(in.Parent)[r.Adr()] = res
			rt2.RegOf(in.Parent)[in.Key] = r.Adr()
			return End()
		})
	})
}
