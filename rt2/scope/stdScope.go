package scope

import (
	"container/list"
	"cp/node"
	"cp/object"
	"fmt"
	"reflect"
	"rt2/context"
	rt_mod "rt2/module"
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
	heap map[object.Object]interface{}
	root node.Node
}

type undefined struct{}
type param struct{}

var undef *undefined = new(undefined)
var par *param = new(param)

func (m *manager) init() *manager {
	m.areas = list.New()
	return m
}

func (m *manager) Allocate(n node.Node) {
	mod := rt_mod.DomainModule(m.Domain())
	h := new(area)
	h.heap = make(map[object.Object]interface{})
	h.root = n
	for _, o := range mod.Objects[n] {
		//fmt.Println(reflect.TypeOf(o))
		switch o.(type) {
		case object.VariableObject:
			h.heap[o] = undef
		case object.ParameterObject:
			h.heap[o] = par
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
		m.Update(o, func(old interface{}) interface{} {
			return m.Select(val.Object())
		})
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
	if ret == undef {
		ret = nil
	} else if ret == par {
		panic("")
	}
	return ret
}

func (m *manager) Update(o object.Object, val Value) {
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
	tmp := x.heap[o]
	if tmp == undef {
		tmp = val(nil)
	} else {
		tmp = val(tmp)
	}
	if tmp == nil {
		tmp = undef
	}
	x.heap[o] = tmp
	fmt.Println("set", x.heap[o], reflect.TypeOf(x.heap[o]))
}

func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}

func (m *manager) Handle(msg interface{}) {}
