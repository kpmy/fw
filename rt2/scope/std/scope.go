package std

import (
	"container/list"
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	rt_mod "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
)

type manager struct {
	d     context.Domain
	areas *list.List
}

type area struct {
	root  node.Node
	data  map[scope.ID]interface{}
	ready bool
}

type value interface {
	set(x interface{})
	get() interface{}
}

type reference interface {
	id() scope.ID
}

type array interface {
	set(i int64, x interface{})
	get(i int64) interface{}
	upd(x []interface{})
	sel() []interface{}
}

type basic struct {
	link object.Object
	data interface{}
}

type record interface {
	set(field string, x interface{})
	get(field string) interface{}
	init(heap *area)
}

func (b *basic) set(i interface{}) { b.data = i }

func (b *basic) get() interface{} { return b.data }

type ref struct {
	link object.Object
	ref  scope.ID
}

func (r *ref) id() scope.ID { return r.ref }

type arr struct {
	link object.Object
	par  int64
	data []interface{}
}

func (a *arr) get(i int64) interface{} { return a.data[i] }

func (a *arr) set(i int64, x interface{}) { a.data[i] = x }

func (a *arr) sel() []interface{} { return a.data }

func (a *arr) upd(x []interface{}) { a.data = x }

type rec struct {
	link object.Object
	heap *area
}

func (r *rec) set(f string, x interface{}) {
	panic(0)
}

func (r *rec) get(f string) interface{} { return nil }

func (r *rec) init(h *area) { r.heap = h }

func nm() scope.Manager {
	m := &manager{areas: list.New()}
	return m
}

func init() {
	scope.New = nm
	scope.Designator = design
	scope.FindObjByName = FindObjByName
}

func design(n node.Node) (id scope.ID) {
	switch x := n.(type) {
	case node.VariableNode, node.ParameterNode:
		id = scope.ID{Name: x.Object().Name()}
	case node.FieldNode:
		panic(0)
		//id = scope.ID{Name: x.Left().Object().Name(), Field: x.Object().Name()}
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(n)))
	}
	return id
}

func odesign(o object.Object) (id scope.ID) {
	switch x := o.(type) {
	case object.VariableObject, object.ParameterObject:
		id = scope.ID{Name: x.Name()}
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(o)))
	}
	return id
}

func obj(o object.Object) (key scope.ID, val interface{}) {
	switch x := o.(type) {
	case object.ConstantObject, object.ProcedureObject, object.TypeObject:
	case object.VariableObject, object.FieldObject:
		fmt.Println(x.Name())
		key = scope.ID{Name: x.Name()}
		switch t := x.Complex().(type) {
		case nil:
			val = &basic{link: o}
		case object.BasicType:
			val = &basic{link: o}
		case object.ArrayType:
			val = &arr{link: o, par: t.Len()}
		case object.DynArrayType:
			val = &arr{link: o}
		case object.RecordType:
			val = &rec{link: o}
		default:
			fmt.Println("unexpected", reflect.TypeOf(t))
		}
	case object.ParameterObject:
		fmt.Println("'" + x.Name())
		key = scope.ID{Name: x.Name()}
		val = &ref{link: o}
	default:
		fmt.Println(reflect.TypeOf(o))
	}
	return key, val
}

func (m *manager) Allocate(n node.Node, final bool) {
	h := &area{ready: final, root: n, data: make(map[scope.ID]interface{})}
	mod := rt_mod.DomainModule(m.Domain())
	var alloc func(h *area, o object.Object)
	alloc = func(h *area, o object.Object) {
		if k, v := obj(o); v != nil {
			h.data[k] = v
			switch rec := v.(type) {
			case record:
				hh := &area{ready: final, root: n, data: make(map[scope.ID]interface{})}
				rec.init(hh)
				switch t := o.Complex().(type) {
				case object.RecordType:
					for rec := t; rec != nil; {
						for x := rec.Link(); x != nil; x = x.Link() {
							fmt.Println(o.Name(), ".", x.Name())
							alloc(hh, x)
						}
						if rec.Base() != "" {
							rec = mod.TypeByName(n, rec.Base()).(object.RecordType)
						} else {
							rec = nil
						}
					}
				}
			}
		} else {
			fmt.Println("nil allocated", reflect.TypeOf(o))
		}
	}
	for _, o := range mod.Objects[n] {
		alloc(h, o)
	}
	m.areas.PushFront(h)
	fmt.Println("allocate")
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
		switch ov := val.(type) {
		case node.ConstantNode:
			switch next.(type) {
			case object.VariableObject:
				m.Update(odesign(next), func(old interface{}) interface{} {
					return ov.Data()
				})
			case object.ParameterObject:
				k, v := scope.ID{Name: next.Name()}, &basic{link: o}
				h.data[k] = v
				m.Update(odesign(next), func(old interface{}) interface{} {
					return ov.Data()
				})
			default:
				panic("unknown value")
			}
		case node.VariableNode, node.ParameterNode:
			switch next.(type) {
			case object.VariableObject:
				m.Update(odesign(next), func(old interface{}) interface{} {
					return m.Select(odesign(ov.Object()))
				})
			case object.ParameterObject:
				h.data[scope.ID{Name: next.Name()}].(*ref).ref = design(ov)
			}
		default:
			panic("unknown value")
		}
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

func (m *manager) Select(i scope.ID) (ret interface{}) {
	fmt.Println("select", i)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		if h.ready {
			ret = h.data[i]
			switch x := ret.(type) {
			case value:
				ret = x.get()
			case reference:
				i = x.id()
				ret = nil
			case array:
				ret = x.sel()
			case nil:
				//do nothing
			default:
				panic(0)
			}
		}
	}
	//	assert.For(ret != nil, 40)
	fmt.Println(ret)
	return ret
}

func arrConv(x interface{}) []interface{} {
	switch a := x.(type) {
	case string:
		s := []rune(a)
		ret := make([]interface{}, 0)
		for i := 0; i < len(s); i++ {
			ret = append(ret, s[i])
		}
		return ret
	case []interface{}:
		return a
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(x)))
	}
}

func (m *manager) Update(i scope.ID, val scope.ValueFor) {
	assert.For(val != nil, 21)
	var x interface{}
	fmt.Println("update", i)
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*area)
		x = h.data[i]
		switch x := x.(type) {
		case value:
			old := x.get()
			tmp := val(old)
			assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
			fmt.Println(tmp)
			x.set(tmp)
		case reference:
			i = x.id()
			x = nil
		case array:
			old := x.sel()
			tmp := val(old)
			assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
			fmt.Println(tmp)
			x.upd(arrConv(tmp))
		case record:

		case nil:
			//do nothing
		default:
			panic(fmt.Sprintln("unhandled", reflect.TypeOf(x)))
		}
	}
	assert.For(x != nil, 40)
}

func (m *manager) Init(d context.Domain) { m.d = d }

func (m *manager) Domain() context.Domain { return m.d }

func (m *manager) Handle(msg interface{}) {}

func FindObjByName(mgr scope.Manager, name string) (ret object.Object) {
	assert.For(name != "", 20)
	m := mgr.(*manager)
	for e := m.areas.Front(); (e != nil) && (ret == nil); e = e.Next() {
		h := e.Value.(*area)
		x := h.data[scope.ID{Name: name}]
		switch x.(type) {
		case *basic:
			ret = x.(*basic).link
		default:
			fmt.Println("no such object")
		}
	}
	return ret
}
