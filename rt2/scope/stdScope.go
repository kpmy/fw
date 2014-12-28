package scope

import (
	"container/list"
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	rt_mod "fw/rt2/module"
	"fw/utils"
	"math/rand"
	"reflect"
	"strconv"
	"ypk/assert"
)

func Id(of interface{}) (ret string) {
	assert.For(of != nil, 20)
	switch of.(type) {
	case object.Object:
		fmt.Println("id", of.(object.Object).Name(), reflect.TypeOf(of))
		//panic("fuck objects, use nodes")
		ret = of.(object.Object).Name()
	default:
		fmt.Println(reflect.TypeOf(of))
		panic("cannot identify")
	}
	assert.For(ret != "", 60)
	return ret
}

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
	opts  struct {
		startFromArea *area
	}
}

type area struct {
	heap  map[string]value
	cache map[string]object.Object
	root  node.Node
	ready bool
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
	ref  object.Object
	mgr  *manager
	area *area
}

// маскирует объект-параметр
type mask struct {
	object.Object
	id int
}

func nextMask() int {
	return rand.Int()
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
	v.mgr.Update(Id(v.ref), func(old interface{}) interface{} {
		return x
	})
}

func (v *indirect) Get() (ret interface{}) {
	assert.For(v.ref != nil, 20)
	_, ok := v.ref.(*mask)
	if !ok {
		v.mgr.opts.startFromArea = v.area
	}
	ret = v.mgr.Select(Id(v.ref))
	v.mgr.opts.startFromArea = nil
	return ret
}

func (m *manager) init() *manager {
	m.areas = list.New()
	return m
}

func (m *manager) Allocate(n node.Node, final bool) {
	mod := rt_mod.DomainModule(m.Domain())
	h := new(area)
	h.ready = final
	h.heap = make(map[string]value)
	h.cache = make(map[string]object.Object)
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
					h.heap[Id(o)] = &direct{data: def}
					h.cache[Id(o)] = o
				}
			default:
				h.heap[Id(o)] = &direct{data: def}
				h.cache[Id(o)] = o
			}
		case object.ParameterObject:
			h.heap[Id(o)] = &indirect{mgr: m, area: h}
			h.cache[Id(o)] = o
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
			m.Update(Id(o), func(old interface{}) interface{} {
				return val.(node.ConstantNode).Data()
			})
		case object.ParameterObject:
			assert.For(a.heap[Id(o)].(*indirect).ref == nil, 40)
			m := &mask{id: nextMask()}
			a.heap[Id(o)].(*indirect).ref = m
			a.heap[Id(m)] = &direct{data: val.(node.ConstantNode).Data()}
		default:
			panic("unknown value")
		}
	case node.VariableNode, node.ParameterNode:
		switch o.(type) {
		case object.VariableObject:
			m.Update(Id(o), func(old interface{}) interface{} {
				return m.Select(Id(val.Object()))
			})
		case object.ParameterObject:
			a.heap[Id(o)].(*indirect).ref = val.Object()
			fmt.Println(val.Object())
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
	assert.For(!h.ready, 22)
	val := _val
	fmt.Println("initialize")
	for next := o; next != nil; next = next.Link() {
		assert.For(val != nil, 40)
		fmt.Println(reflect.TypeOf(next), next.Name(), ":", next.Type())
		fmt.Println(reflect.TypeOf(val))
		m.set(h, next, val)
		val = val.Link()
	}
	h.ready = true
}

func (m *manager) Dispose(n node.Node) {
	e := m.areas.Front()
	if e != nil {
		h := e.Value.(*area)
		assert.For(h.root == n, 21)
		m.areas.Remove(e)
		fmt.Println("dispose")
	}
}

func FindObjByName(mgr Manager, name string) (ret object.Object) {
	assert.For(name != "", 20)
	m := mgr.(*manager)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		ret = h.cache[name]
	}
	return ret
}

func (m *manager) Select(id string) (ret interface{}) {
	assert.For(id != "", 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		if (h != m.opts.startFromArea) && h.ready {
			ret = h.heap[id]
		}
	}
	assert.For(ret != nil, 40)
	ret = ret.(value).Get()
	if _, ok := ret.(*dummy); ok {
		ret = nil
	}
	return ret
}

func (m *manager) Update(id string, val ValueFor) {
	assert.For(id != "", 20)
	assert.For(val != nil, 21)
	var x *area
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*area)
		if h.heap[id] != nil {
			x = h
		}
	}
	assert.For(x != nil, 40)
	old := x.heap[id].Get()
	if old == def {
		old = nil
	}
	tmp := val(old)
	assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
	x.heap[id].Set(tmp)
}

func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}

func (m *manager) Handle(msg interface{}) {}

func (m *mask) Name() string { return strconv.Itoa(m.id) }
