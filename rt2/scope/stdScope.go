package scope

import (
	"container/list"
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	rt_mod "fw/rt2/module"
	"fw/utils"
	"reflect"
	"ypk/assert"
)

func This(i interface{}) Manager {
	assert.For(i != nil, 20)
	return i.(Manager)
}

func New() Manager {
	return new(manager).init()
}

type manager struct {
	d     context.Domain
	areas *list.List
}

type area struct {
	heap map[object.Object]value
	root node.Node
}

type value interface {
	Set(x interface{})
	Get() interface{}
}

type direct struct {
	value
	data interface{}
}

type indirect struct {
	value
	ref object.Object
	mgr Manager
}

type mask struct {
	object.Object
}

type dummy struct{}

var def *dummy = &dummy{}

func (v *direct) Set(x interface{}) {
	assert.For(x != nil, 20)
	v.data = x
	utils.Println("set", x, reflect.TypeOf(x))
}

func (v *direct) Get() interface{} { return v.data }

func (v *indirect) Set(x interface{}) {
	assert.For(x != nil, 20)
	assert.For(v.ref != nil, 21)
	v.mgr.UpdateObj(v.ref, func(old interface{}) interface{} {
		return x
	})
}

func (v *indirect) Get() interface{} {
	assert.For(v.ref != nil, 20)
	return v.mgr.SelectObj(v.ref)
}

func (m *manager) init() *manager {
	m.areas = list.New()
	return m
}

func (m *manager) Allocate(n node.Node) {
	mod := rt_mod.DomainModule(m.Domain())
	h := new(area)
	h.heap = make(map[object.Object]value)
	h.root = n
	for _, o := range mod.Objects[n] {
		//fmt.Println(reflect.TypeOf(o))
		switch o.(type) {
		case object.VariableObject:
			switch o.(object.VariableObject).Type() {
			case object.COMPLEX:
				switch o.(object.VariableObject).Complex().(type) {
				case object.RecordType:
					for rec := o.(object.VariableObject).Complex().(object.RecordType); rec != nil; {
						for x := rec.Link(); x != nil; x = x.Link() {
							fmt.Println(o.Name(), ".", x.Name())
						}
						if rec.Base() != "" {
							rec = mod.TypeByName(n, rec.Base()).(object.RecordType)
						} else {
							rec = nil
						}
					}
				default:
					h.heap[o] = &direct{data: def}
				}
			default:
				h.heap[o] = &direct{data: def}
			}
		case object.ParameterObject:
			h.heap[o] = &indirect{mgr: m}
		default:
			fmt.Println("wrong object type", reflect.TypeOf(o))
		}
	}
	m.areas.PushFront(h)
	fmt.Println("allocate", len(h.heap), "obj")
}

func (m *manager) set(a *area, o object.Object, val node.Node) {
	switch val.(type) {
	case node.ConstantNode:
		switch o.(type) {
		case object.VariableObject:
			m.UpdateObj(o, func(old interface{}) interface{} {
				return val.(node.ConstantNode).Data()
			})
		case object.ParameterObject:
			assert.For(a.heap[o].(*indirect).ref == nil, 40)
			m := &mask{}
			a.heap[o].(*indirect).ref = m
			a.heap[m] = &direct{data: val.(node.ConstantNode).Data()}
		default:
			panic("unknown value")
		}
	case node.VariableNode, node.ParameterNode, node.FieldNode:
		switch o.(type) {
		case object.VariableObject:
			m.UpdateObj(o, func(old interface{}) interface{} {
				return m.SelectObj(val.Object())
			})
		case object.ParameterObject:
			a.heap[o].(*indirect).ref = val.Object()
		}
	default:
		panic("unknown value")
	}
}

func (m *manager) Initialize(n node.Node, o object.Object, _val node.Node) {
	e := m.areas.Front()
	assert.For(e != nil, 20)
	h := e.Value.(*area)
	assert.For(h.root == n, 21)
	val := _val
	fmt.Println("initialize")
	for next := o; next != nil; next = next.Link() {
		assert.For(val != nil, 40)
		fmt.Println(reflect.TypeOf(next), next.Name(), ":", next.Type())
		fmt.Println(reflect.TypeOf(val))
		m.set(h, next, val)
		val = val.Link()
	}
}

func (m *manager) Dispose(n node.Node) {
	e := m.areas.Front()
	assert.For(e != nil, 20)
	h := e.Value.(*area)
	assert.For(h.root == n, 21)
	m.areas.Remove(e)
	fmt.Println("dispose")
}

func (m *manager) FindObjByName(name string) (ret object.Object) {
	assert.For(name != "", 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		for k, _ := range h.heap {
			if k.Name() == name {
				ret = k
			}
		}
	}
	return ret
}

func (m *manager) SelectObj(o object.Object) (ret interface{}) {
	assert.For(o != nil, 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		ret = h.heap[o]
	}
	assert.For(ret != nil, 40)
	ret = ret.(value).Get()
	if _, ok := ret.(*dummy); ok {
		ret = nil
	}
	return ret
}

func (m *manager) UpdateObj(o object.Object, val ValueFor) {
	assert.For(o != nil, 20)
	assert.For(val != nil, 21)
	var x *area
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*area)
		if h.heap[o] != nil {
			x = h
		}
	}
	assert.For(x != nil, 40)
	old := x.heap[o].Get()
	if old == def {
		old = nil
	}
	tmp := val(old)
	assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
	x.heap[o].Set(tmp)
}

func (m *manager) SelectNode(n node.Node) interface{} {
	return nil
}

func (m *manager) UpdateNode(n node.Node, val ValueFor) {

}
func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}

func (m *manager) Handle(msg interface{}) {}

func (m *mask) Name() string { return "" }
