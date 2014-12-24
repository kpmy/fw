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
	v.mgr.Update(v.ref, func(old interface{}) interface{} {
		return x
	})
}

func (v *indirect) Get() interface{} {
	assert.For(v.ref != nil, 20)
	return v.mgr.Select(v.ref)
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
			h.heap[o] = &direct{data: def}
		case object.ParameterObject:
			h.heap[o] = &indirect{mgr: m}
		}
	}
	m.areas.PushFront(h)
	fmt.Println("allocate", len(h.heap), "obj")
}

func (m *manager) set(a *area, o object.Object, val node.Node) {
	switch val.(type) {
	case node.ConstantNode:
		m.Update(o, func(old interface{}) interface{} {
			return val.(node.ConstantNode).Data()
		})
	case node.VariableNode, node.ParameterNode:
		switch o.(type) {
		case object.VariableObject:
			m.Update(o, func(old interface{}) interface{} {
				return m.Select(val.Object())
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

func (m *manager) Select(o object.Object) (ret interface{}) {
	assert.For(o != nil, 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		ret = h.heap[o]
	}
	assert.For(ret != nil, 40)
	return ret.(value).Get()
}

func (m *manager) Update(o object.Object, val ValueFor) {
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

func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}

func (m *manager) Handle(msg interface{}) {}
