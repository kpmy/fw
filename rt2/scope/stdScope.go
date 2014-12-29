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

const NONE int64 = -1

func Id(of interface{}) (ret ID) {
	assert.For(of != nil, 20)
	ret.Index = -1
	switch of.(type) {
	case object.FieldObject:
		panic("fuck objects, use nodes")
	case object.VariableObject, object.ParameterObject:
		//fmt.Println("id", of.(object.Object).Name(), reflect.TypeOf(of))
		ret.Name = of.(object.Object).Name()
	case node.FieldNode:
		f := of.(node.FieldNode)
		ret.Name = f.Left().Object().Name() + "." + f.Object().Name()
	case node.IndexNode:
		f := of.(node.IndexNode)
		ret.Name = f.Left().Object().Name()
	case *mask:
		ret.Name = of.(*mask).Name()
	default:
		fmt.Println(reflect.TypeOf(of))
		panic("cannot identify")
	}
	assert.For(ret.Name != "", 60)
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
	Set(id ID, x interface{})
	Get(id ID) interface{}
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

func defval(x ...interface{}) interface{} {
	if x != nil {
		return x[0]
	}
	return def
}

func (v *direct) Set(id ID, x interface{}) {
	assert.For(x != nil, 20)
	if id.Index == NONE {
		v.data = x
		utils.PrintScope("set", x, reflect.TypeOf(x))
	} else {
		_, ok := v.data.([]rune)
		if ok {
			utils.PrintScope("set indexed", x)
		} else {
			panic("unsupported")
		}
	}
}

func (v *direct) Get(id ID) interface{} {
	if id.Index == -1 {
		return v.data
	} else {
		v, ok := v.data.([]rune)
		if ok {
			v := string(v)
			return v[id.Index]
		}
		panic("unsupported")
	}
}

func (v *indirect) Set(id ID, x interface{}) {
	assert.For(x != nil, 20)
	assert.For(v.ref != nil, 21)
	t := Id(v.ref)
	t.Index = id.Index
	v.mgr.Update(t, func(old interface{}) interface{} {
		return x
	})
}

func (v *indirect) Get(id ID) (ret interface{}) {
	assert.For(v.ref != nil, 20)
	_, ok := v.ref.(*mask)
	if !ok {
		v.mgr.opts.startFromArea = v.area
	}
	t := Id(v.ref)
	t.Index = id.Index
	ret = v.mgr.Select(t)
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
							//fmt.Println(o.Name(), ".", x.Name())
							h.heap[o.Name()+"."+x.Name()] = &direct{data: defval()}
						}
						if rec.Base() != "" {
							rec = mod.TypeByName(n, rec.Base()).(object.RecordType)
						} else {
							rec = nil
						}
					}
				case object.ArrayType:
					arr := o.(object.VariableObject).Complex().(object.ArrayType)
					h.heap[Id(o).Name] = &direct{data: defval(make([]rune, arr.Len()), arr.Len())}
					h.cache[Id(o).Name] = o
				default:
					h.heap[Id(o).Name] = &direct{data: defval()}
					h.cache[Id(o).Name] = o
				}
			default:
				h.heap[Id(o).Name] = &direct{data: defval()}
				h.cache[Id(o).Name] = o
			}
		case object.ParameterObject:
			h.heap[Id(o).Name] = &indirect{mgr: m, area: h}
			h.cache[Id(o).Name] = o
		case object.ProcedureObject, object.ConstantObject, object.FieldObject:
			//do nothing
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
			assert.For(a.heap[Id(o).Name].(*indirect).ref == nil, 40)
			m := &mask{id: nextMask()}
			a.heap[Id(o).Name].(*indirect).ref = m
			a.heap[Id(m).Name] = &direct{data: val.(node.ConstantNode).Data()}
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
			a.heap[Id(o).Name].(*indirect).ref = val.Object()
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

func (m *manager) Select(id ID) (ret interface{}) {
	assert.For(id.Name != "", 20)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		if (h != m.opts.startFromArea) && h.ready {
			ret = h.heap[id.Name]
		}
	}
	assert.For(ret != nil, 40)
	ret = ret.(value).Get(id)
	if _, ok := ret.(*dummy); ok {
		ret = nil
	}
	return ret
}

func (m *manager) Update(id ID, val ValueFor) {
	assert.For(id.Name != "", 20)
	assert.For(val != nil, 21)
	var x *area
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*area)
		if h.heap[id.Name] != nil {
			x = h
		}
	}
	assert.For(x != nil, 40)
	old := x.heap[id.Name].Get(id)
	if old == def {
		old = nil
	}
	tmp := val(old)
	assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
	x.heap[id.Name].Set(id, tmp)
}

func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}

func (m *manager) Handle(msg interface{}) {}

func (m *mask) Name() string { return strconv.Itoa(m.id) }
