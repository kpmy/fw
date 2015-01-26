package modern

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"runtime"
	"ypk/assert"
	"ypk/halt"
)

type halloc struct {
	area *area
}

func fin(x interface{}) {
	switch p := x.(type) {
	case *ptrValue:
		defer func() {
			mod := rtm.ModuleOfType(p.scope.Domain(), p.link.Complex())
			ol := mod.Objects[mod.Enter]
			var fn object.ProcedureObject
			for _, _po := range ol {
				switch po := _po.(type) {
				case object.ProcedureObject:
					if po.Name() == "FINALIZE" && po.Link().Complex() == p.link.Complex() {
						fn = po
						break
					}
				}
			}
			if fn != nil {
				global := p.scope.Domain().Discover(context.UNIVERSE).(context.Domain)
				root := global.Discover(context.STACK).(frame.Stack)
				cn := node.New(constant.CALL, int(cp.SomeAdr()))
				ol := mod.NodeByObject(fn)
				assert.For(len(ol) <= 1, 40)
				cn.SetLeft(ol[0])
				cc := node.New(constant.CONSTANT, int(cp.SomeAdr())).(node.ConstantNode)
				cc.SetData(p)
				cc.SetType(object.COMPLEX)
				cn.SetRight(cc)
				nf := rt2.New(cn)
				nf.Init(global.Discover(mod.Name).(context.Domain))
				root.Queue(nf)
			}
			p.scope.Target().(scope.HeapAllocator).Dispose(p.id)
		}()
	}
}

func (h *halloc) Allocate(n node.Node, par ...interface{}) scope.ValueFor {
	//fmt.Println("HEAP ALLOCATE")
	mod := rtm.ModuleOfNode(h.area.d, n)
	if h.area.data == nil {
		h.area.data = append(h.area.data, newlvl())
	}
	var ol []object.Object
	skip := make(map[cp.ID]interface{})
	l := h.area.top()
	var res scope.Value
	f_res := func(x scope.Value) scope.Value {
		return res
	}
	var talloc func(t object.PointerType)
	talloc = func(t object.PointerType) {
		switch bt := t.Base().(type) {
		case object.RecordType:
			fake := object.New(object.VARIABLE, int(cp.SomeAdr()))
			fake.SetComplex(bt)
			fake.SetType(object.COMPLEX)
			fake.SetName("{}")
			l.alloc(mod, nil, append(ol, fake), skip)
			res = &ptrValue{scope: h.area, id: fake.Adr(), link: n.Object()}
		case object.DynArrayType:
			assert.For(len(par) > 0, 20)
			fake := object.New(object.VARIABLE, int(cp.SomeAdr()))
			fake.SetComplex(bt)
			fake.SetType(object.COMPLEX)
			fake.SetName("[]")
			l.alloc(mod, nil, append(ol, fake), skip)
			h.area.Select(fake.Adr(), func(v scope.Value) {
				arr, ok := v.(*dynarr)
				assert.For(ok, 60)
				arr.Set(par[0].(scope.Value))
			})
			res = &ptrValue{scope: h.area, id: fake.Adr(), link: n.Object()}
		default:
			halt.As(100, fmt.Sprintln("cannot allocate", reflect.TypeOf(bt)))
		}
	}
	switch v := n.(type) {
	case node.VariableNode:
		switch t := v.Object().Complex().(type) {
		case object.PointerType:
			talloc(t)
			//h.area.data[0].alloc(mod, nil, )
		default:
			halt.As(100, reflect.TypeOf(t))
		}
	default:
		halt.As(101, reflect.TypeOf(v))
	}
	assert.For(res != nil, 60)
	runtime.SetFinalizer(res, fin)
	return f_res
}

func (h *halloc) Dispose(id cp.ID) {
	h.area.Select(id, func(v scope.Value) {
		//fmt.Println("dispose", v)
		a := h.area.top()
		k := a.k[id]
		delete(a.l, k)
		delete(a.r, k)
		delete(a.v, k)
		delete(a.k, id)
	})
}

func (a *halloc) Join(m scope.Manager) { a.area = m.(*area) }
