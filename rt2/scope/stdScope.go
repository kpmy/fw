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

type heapObj struct {
	heap map[object.Object]interface{}
	root node.Node
}

type undefined struct{}

var undef *undefined = new(undefined)

func (m *manager) init() *manager {
	m.areas = list.New()
	return m
}

func (m *manager) Allocate(n node.Node) {
	mod := rt_mod.DomainModule(m.Domain())
	h := new(heapObj)
	h.heap = make(map[object.Object]interface{})
	h.root = n
	for _, o := range mod.Objects[n] {
		fmt.Println(reflect.TypeOf(o))
		switch o.(type) {
		case object.VariableObject:
			h.heap[o] = undef
		}
	}
	m.areas.PushFront(h)
	fmt.Println("allocate", len(h.heap), "obj")
}

func (m *manager) Dispose(n node.Node) {
	e := m.areas.Front()
	assert.For(e != nil, 20)
	h := e.Value.(*heapObj)
	assert.For(h.root == n, 21)
	m.areas.Remove(e)
	fmt.Println("dispose")
}

func (m *manager) Select(o object.Object) (ret interface{}) {
	assert.For(o != nil, 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*heapObj)
		ret = h.heap[o]
	}
	assert.For(ret != nil, 40)
	if ret == undef {
		ret = nil
	}
	return ret
}

func (m *manager) Update(o object.Object, val Value) {
	assert.For(o != nil, 20)
	assert.For(val != nil, 21)
	var x *heapObj
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*heapObj)
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
