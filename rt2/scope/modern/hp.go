package modern

import (
	"fmt"
	//cpm "fw/cp/module"
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

type halloc struct {
	area *area
}

func (h *halloc) Allocate(n node.Node, par ...interface{}) scope.ValueFor {
	fmt.Println("HEAP ALLOCATE")
	mod := rtm.ModuleOfNode(h.area.d, n)
	if h.area.data == nil {
		h.area.data = append(h.area.data, newlvl())
	}
	var ol []object.Object
	skip := make(map[cp.ID]interface{})
	l := h.area.data[0]
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
			res = &ptrValue{scope: h.area, id: fake.Adr()}
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
			res = &ptrValue{scope: h.area, id: fake.Adr()}
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
	return f_res
}

func (h *halloc) Dispose(n node.Node) {}

func (a *halloc) Join(m scope.Manager) { a.area = m.(*area) }
