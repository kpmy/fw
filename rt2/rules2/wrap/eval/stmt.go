package eval

import (
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/node"
	"fw/cp/object"
	"fw/cp/traps"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"fw/utils"
	"log"
	"math/big"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

func makeTrap(f frame.Frame, err traps.TRAP) (out OUT) {
	trap := node.New(constant.TRAP, cp.Some()).(node.TrapNode)
	code := node.New(constant.CONSTANT, cp.Some()).(node.ConstantNode)
	code.SetType(object.INTEGER)
	code.SetData(int32(err))
	trap.SetLeft(code)
	rt2.Push(rt2.New(trap), f)
	return Later(Tail(STOP))
}

func doEnter(in IN) OUT {
	e := in.IR.(node.EnterNode)
	var next Do
	tail := func(IN) (out OUT) {
		body := in.IR.Right()
		switch {
		case body == nil:
			return End()
		case body != nil && in.Parent != nil:
			rt2.Push(rt2.New(body), in.Frame)
			return Later(Tail(STOP))
		case body != nil && in.Parent == nil: //секция BEGIN
			rt2.Push(rt2.New(body), in.Frame)
			end := in.IR.Link()
			if end != nil { //секция CLOSE
				out.Do = func(in IN) OUT {
					in.Frame.Root().PushFor(rt2.New(end), in.Frame)
					return OUT{Do: Tail(STOP), Next: LATER}
				}
			} else {
				out.Do = Tail(STOP)
			}
			out.Next = BEGIN
		}
		return
	}
	var sm scope.Manager
	switch e.Enter() {
	case enter.MODULE:
		sm = rt2.ModScope(in.Frame)
	case enter.PROCEDURE:
		sm = rt2.CallScope(in.Frame)
	}
	if e.Object() != nil { //параметры процедуры
		par, ok := rt2.RegOf(in.Frame)[e.Object()].(node.Node)
		//fmt.Println(rt2.DataOf(f)[n.Object()])
		//fmt.Println(ok)
		if ok {
			sm.Target().(scope.ScopeAllocator).Allocate(e, false)
			next = func(in IN) OUT {
				seq, _ := sm.Target().(scope.ScopeAllocator).Initialize(e,
					scope.PARAM{Objects: e.Object().Link(),
						Values: par,
						Frame:  in.Frame,
						Tail:   Propose(tail)})
				return Later(Expose(seq))
			}
		} else {
			sm.Target().(scope.ScopeAllocator).Allocate(e, true)
			next = tail
		}
	} else {
		sm.Target().(scope.ScopeAllocator).Allocate(in.IR, true)
		next = tail
	}
	return Now(next)
}

func inc_dec_seq(in IN, code operation.Operation) OUT {
	n := in.IR
	a := node.New(constant.ASSIGN, cp.Some()).(node.AssignNode)
	a.SetStatement(statement.ASSIGN)
	a.SetLeft(n.Left())
	op := node.New(constant.DYADIC, cp.Some()).(node.OperationNode)
	op.SetOperation(code)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	a.SetRight(op)
	rt2.Push(rt2.New(a), in.Frame)
	/*seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(n.Left().Object().Adr(), scope.Simple(rt2.ValueOf(f)[op.Adr()]))
		return frame.End()
	}
	ret = frame.LATER */
	return Later(Tail(STOP))
}

func doAssign(in IN) (out OUT) {
	const (
		right = "assign:right"
		left  = "assign:left"
	)
	a := in.IR.(node.AssignNode)
	switch a.Statement() {
	case statement.ASSIGN:
		out = GetExpression(in, right, a.Right(), func(in IN) OUT {
			id := KeyOf(in, right)
			val := rt2.ValueOf(in.Frame)[id]
			assert.For(val != nil, 40, id)
			return GetDesignator(in, left, a.Left(), func(in IN) OUT {
				id := KeyOf(in, left)
				v, ok := rt2.ValueOf(in.Frame)[id].(scope.Variable)
				assert.For(ok, 41, reflect.TypeOf(v))
				v.Set(val)
				return End()
			})
		})
	case statement.INC, statement.INCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode, node.FieldNode:
			out = inc_dec_seq(in, operation.PLUS)
		default:
			halt.As(100, "wrong left", reflect.TypeOf(a.Left()))
		}
	case statement.DEC, statement.EXCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode, node.FieldNode:
			out = inc_dec_seq(in, operation.MINUS)
		default:
			halt.As(100, "wrong left", reflect.TypeOf(a.Left()))
		}
	case statement.NEW:
		obj := a.Left().Object()
		assert.For(obj != nil, 20)
		heap := in.Frame.Domain().Discover(context.HEAP).(scope.Manager).Target().(scope.HeapAllocator)
		if a.Right() != nil {
			out = GetExpression(in, right, a.Right(), func(in IN) OUT {
				size := rt2.ValueOf(in.Frame)[KeyOf(in, right)]
				return GetDesignator(in, left, a.Left(), func(in IN) OUT {
					v := rt2.ValueOf(in.Frame)[KeyOf(in, left)].(scope.Variable)
					fn := heap.Allocate(obj, obj.Complex().(object.PointerType), size)
					v.Set(fn)
					return End()
				})
			})
		} else {
			out = GetDesignator(in, left, a.Left(), func(IN) OUT {
				v := rt2.ValueOf(in.Frame)[KeyOf(in, left)].(scope.Variable)
				fn := heap.Allocate(obj, obj.Complex().(object.PointerType))
				v.Set(fn)
				return End()
			})
		}
	default:
		halt.As(100, "unsupported assign statement", a.Statement())
	}
	return
}

func doIf(in IN) OUT {
	const left = "if:left:if"
	i := in.IR.(node.IfNode)
	return GetExpression(in, left, i.Left(), func(in IN) OUT {
		val := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		assert.For(val != nil, 20)
		rt2.ValueOf(in.Parent)[i.Adr()] = val
		rt2.RegOf(in.Parent)[in.Key] = i.Adr()
		return End()
	})
}

func doCondition(in IN) OUT {
	const left = "if:left"
	i := in.IR.(node.ConditionalNode)
	rt2.RegOf(in.Frame)[0] = i.Left() // if
	var next Do
	next = func(in IN) OUT {
		last := rt2.RegOf(in.Frame)[0].(node.Node)
		fi := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		done := scope.GoTypeFrom(fi).(bool)
		rt2.RegOf(in.Frame)[0] = nil
		rt2.ValueOf(in.Frame)[KeyOf(in, left)] = nil

		if done && last.Right() != nil {
			rt2.Push(rt2.New(last.Right()), in.Frame)
			return Later(Tail(STOP))
		} else if last.Right() == nil {
			return End()
		} else if last.Link() != nil { //elsif
			rt2.RegOf(in.Frame)[0] = last.Link()
			return GetStrange(in, left, i.Left(), next)
		} else if i.Right() != nil { //else
			rt2.Push(rt2.New(i.Right()), in.Frame)
			return Later(Tail(STOP))
		} else if i.Right() == nil {
			return End()
		} else if i.Right() == last {
			return End()
		} else {
			halt.As(100, "wrong if then else")
			panic(100)
		}
	}

	return GetStrange(in, left, i.Left(), next)
}

func doWhile(in IN) OUT {
	const left = "while:left"
	w := in.IR.(node.WhileNode)
	var next Do

	next = func(IN) OUT {
		key := KeyOf(in, left)
		fi := rt2.ValueOf(in.Frame)[key]
		rt2.ValueOf(in.Frame)[key] = nil
		done := scope.GoTypeFrom(fi).(bool)
		if done && w.Right() != nil {
			rt2.Push(rt2.New(w.Right()), in.Frame)
			return Later(func(IN) OUT { return GetExpression(in, left, w.Left(), next) })
		} else if !done {
			return End()
		} else if w.Right() == nil {
			return End()
		} else {
			panic("unexpected while seq")
		}
	}

	return GetExpression(in, left, w.Left(), next)
}

const exitFlag = 1812

func doExit(in IN) OUT {
	in.Frame.Root().ForEach(func(f frame.Frame) (ok bool) {
		n := rt2.NodeOf(f)
		_, ok = n.(node.LoopNode)
		if ok {
			rt2.RegOf(f)[exitFlag] = true
		}
		ok = !ok
		return ok
	})
	return End()
}

func doLoop(in IN) OUT {
	l := in.IR.(node.LoopNode)
	exit, ok := rt2.RegOf(in.Frame)[exitFlag].(bool)
	if ok && exit {
		return End()
	}
	if l.Left() != nil {
		rt2.Push(rt2.New(l.Left()), in.Frame)
		return Later(doLoop)
	} else if l.Left() == nil {
		return End()
	} else {
		panic("unexpected loop seq")
	}
}

func doRepeat(in IN) OUT {
	const right = "return:right"
	r := in.IR.(node.RepeatNode)
	rt2.ValueOf(in.Frame)[r.Right().Adr()] = scope.TypeFromGo(false)
	var next Do
	next = func(in IN) OUT {
		fi := rt2.ValueOf(in.Frame)[r.Right().Adr()]
		done := scope.GoTypeFrom(fi).(bool)
		if !done && r.Left() != nil {
			rt2.Push(rt2.New(r.Left()), in.Frame)
			return Later(func(IN) OUT { return GetExpression(in, right, r.Right(), next) })
		} else if done {
			return End()
		} else if r.Left() == nil {
			return End()
		} else {
			panic("unexpected repeat seq")
		}
	}
	return Now(next)
}

func doTrap(in IN) OUT {
	const left = "trap:left"

	return GetExpression(in, left, in.IR.Left(), func(IN) OUT {
		val := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		log.Println("TRAP:", traps.This(scope.GoTypeFrom(val)))
		return Now(Tail(WRONG))
	})

}

func doReturn(in IN) OUT {
	const left = "return:left"
	r := in.IR.(node.ReturnNode)
	return GetExpression(in, left, r.Left(), func(IN) OUT {
		val := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		if val == nil {
			val, _ = rt2.RegOf(in.Frame)[context.RETURN].(scope.Value)
		}
		assert.For(val != nil, 40)
		rt2.ValueOf(in.Parent)[r.Object().Adr()] = val
		rt2.RegOf(in.Parent)[context.RETURN] = val
		return End()
	})
}

func doCall(in IN) (out OUT) {
	const (
		right = "call:right"
	)
	c := in.IR.(node.CallNode)

	call := func(proc node.Node, d context.Domain) {
		_, ok := proc.(node.EnterNode)
		assert.For(ok, 20, "try call", reflect.TypeOf(proc), proc.Adr(), proc.Object().Adr())
		nf := rt2.New(proc)
		rt2.Push(nf, in.Frame)
		if d != nil {
			rt2.ReplaceDomain(nf, d)
		}
		//передаем ссылку на цепочку значений параметров в данные фрейма входа в процедуру
		if (c.Right() != nil) && (proc.Object() != nil) {
			rt2.RegOf(nf)[proc.Object()] = c.Right()
		} else {
			//fmt.Println("no data for call")
		}
		out = Later(func(in IN) OUT {
			if in.Key != nil {
				val := rt2.ValueOf(in.Frame)[c.Left().Object().Adr(0, 0)]
				if val == nil {
					if rv, _ := rt2.RegOf(in.Frame)[context.RETURN].(scope.Value); rv != nil {
						val = rv
					}
				}
				if val != nil {
					id := cp.ID(cp.Some())
					rt2.ValueOf(in.Parent)[id] = val
					rt2.RegOf(in.Parent)[in.Key] = id
				} else {
					utils.PrintFrame("possibly no return", in.Key, rt2.RegOf(in.Frame), rt2.ValueOf(in.Frame), rt2.RegOf(in.Parent), rt2.ValueOf(in.Parent))
				}
			}
			return End()
		})
	}

	switch p := c.Left().(type) {
	case node.ProcedureNode:
		m := rtm.DomainModule(in.Frame.Domain())
		ml := in.Frame.Domain().Global().Discover(context.MOD).(rtm.List)
		switch p.Object().Mode() {
		case object.LOCAL_PROC, object.EXTERNAL_PROC:
			if imp := p.Object().Imp(); imp == "" || imp == m.Name {
				proc := m.NodeByObject(p.Object())
				assert.For(proc != nil, 40, m.Name, imp, p.Object().Imp(), p.Object().Adr(0, 0), p.Object().Name())
				call(proc[0], nil)
			} else {
				m := ml.Loaded(imp)
				pl := m.ObjectByName(m.Enter, c.Left().Object().Name())
				var proc object.ProcedureObject
				var nl []node.Node
				for _, n := range pl {
					if n.Mode() == p.Object().Mode() {
						proc = n.(object.ProcedureObject)
					}

				}
				nl = m.NodeByObject(proc)
				utils.PrintFrame("foreign call", len(nl), "proc refs", proc)
				call(nl[0], in.Frame.Domain().Global().Discover(imp).(context.Domain))
			}
		case object.TYPE_PROC:
			//sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
			out = GetExpression(in, right, c.Right(), func(IN) (out OUT) {
				var (
					proc []node.Node
					dm   context.Domain
				)
				id := KeyOf(in, right)
				v := rt2.ValueOf(in.Frame)[id]
				t, ct := scope.Ops.TypeOf(v)
				if ct == nil {
					return makeTrap(in.Frame, traps.Default)
				}
				assert.For(ct != nil, 40, id, v, t)
				x := ml.NewTypeCalc()
				x.ConnectTo(ct)
				sup := p.Super()
				for _, ml := range x.MethodList() {
					for _, m := range ml {
						//fmt.Println(m.Obj.Name(), m.Obj.Adr())
						if m.Obj.Name() == p.Object().Name() {
							if !sup {
								proc = append(proc, m.Enter)
								break
							} else {
								sup = false //супервызов должен быть вторым методом с тем же именем
							}
							dm = in.Frame.Domain().Global().Discover(m.Mod.Name).(context.Domain)
						}
					}
					if len(proc) > 0 {
						break
					}
				}
				assert.For(len(proc) > 0, 40, p.Object().Name())
				call(proc[0], dm)
				out = Later(Tail(STOP))
				return
			})
		default:
			halt.As(100, "wrong proc mode ", p.Object().Mode(), p.Object().Adr(), p.Object().Name())
		}
	case node.VariableNode:
		m := rtm.DomainModule(in.Frame.Domain())
		sc := rt2.ScopeFor(in.Frame, p.Object().Adr())
		var obj interface{}
		sc.Select(p.Object().Adr(), func(in scope.Value) {
			obj = scope.GoTypeFrom(in)
		})
		if obj, ok := obj.(object.Object); ok {
			proc := m.NodeByObject(obj)
			call(proc[0], nil)
		} else {
			name := p.Object().Name()
			switch {
			case sys[name] != nil:
				return syscall(in)
			default:
				halt.As(100, "unknown sysproc variable", name)
			}
		}
	default:
		halt.As(100, reflect.TypeOf(p))
	}
	return
}

const isFlag = 1700

func doWith(in IN) OUT {
	const left = "with:left"
	w := in.IR.(node.WithNode)
	f := in.Frame
	rt2.RegOf(f)[isFlag] = w.Left() //if
	var seq Do
	seq = func(in IN) OUT {
		last := rt2.RegOf(f)[isFlag].(node.Node)
		id := KeyOf(in, left)
		v := rt2.ValueOf(f)[id]
		assert.For(v != nil, 40)
		done := scope.GoTypeFrom(v).(bool)
		rt2.ValueOf(f)[id] = nil
		if done && last.Right() != nil {
			rt2.Push(rt2.New(last.Right()), f)
			return Later(Tail(STOP))
		} else if last.Right() == nil {
			return End()
		} else if last.Link() != nil { //elsif
			rt2.RegOf(f)[isFlag] = last.Link()
			return Later(func(IN) OUT { return GetStrange(in, left, last.Link(), seq) })
		} else if w.Right() != nil { //else
			rt2.Push(rt2.New(w.Right()), f)
			return Later(Tail(STOP))
		} else if w.Right() == nil {
			return End()
		} else if last == w.Right() {
			return End()
		} else {
			panic("conditional sequence wrong")
		}
	}
	return GetStrange(in, left, w.Left(), seq)
}

func int32Of(x interface{}) (a int32) {
	//fmt.Println(reflect.TypeOf(x))
	switch v := x.(type) {
	case *int32:
		z := *x.(*int32)
		a = z
	case int32:
		a = x.(int32)
	case *big.Int:
		a = int32(v.Int64())
	default:
		//panic(fmt.Sprintln("unsupported type", reflect.TypeOf(x)))
	}
	return a
}

func doCase(in IN) OUT {
	const left = "case:left"
	c := in.IR.(node.CaseNode)

	var e int

	do := func(in IN) (out OUT) {
		cond := c.Right().(node.ElseNode)
		//fmt.Println("case?", e, cond.Min(), cond.Max())
		if e < cond.Min() || e > cond.Max() { //case?
			if cond.Right() != nil {
				rt2.Push(rt2.New(cond.Right()), in.Frame)
			}
			out.Do = Tail(STOP)
			out.Next = LATER
		} else {
			for next := cond.Left(); next != nil && out.Do == nil; next = next.Link() {
				var ok bool
				//				_c := next.Left()
				for _c := next.Left(); _c != nil && !ok; _c = _c.Link() {
					c := _c.(node.ConstantNode)
					//fmt.Println("const", c.Data(), c.Min(), c.Max())
					if (c.Min() != nil) && (c.Max() != nil) {
						//fmt.Println(e, *c.Max(), "..", *c.Min())
						ok = e >= *c.Min() && e <= *c.Max()
					} else {
						//fmt.Println(e, c.Data())
						ok = int32Of(c.Data()) == int32(e)
					}
				}
				//fmt.Println(ok)
				if ok {
					rt2.Push(rt2.New(next.Right()), in.Frame)
					out.Do = Tail(STOP)
					out.Next = LATER
				}
			}
			if out.Do == nil && cond.Right() != nil {
				rt2.Push(rt2.New(cond.Right()), in.Frame)
				out.Do = Tail(STOP)
				out.Next = LATER
			}
		}
		assert.For(out.Do != nil, 60)
		return out
	}

	return GetExpression(in, left, c.Left(), func(IN) (out OUT) {
		v := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		_x := scope.GoTypeFrom(v)
		switch x := _x.(type) {
		case nil:
			panic("nil")
		case int32:
			e = int(x)
			//fmt.Println("case", e)
			out.Do = do
			out.Next = NOW
			return out
		default:
			halt.As(100, "unsupported case expr", reflect.TypeOf(_x))
		}
		panic(0)
	})
}
