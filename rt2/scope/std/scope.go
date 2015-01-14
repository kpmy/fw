package std

import (
	"container/list"
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rt_mod "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"runtime"
	"ypk/assert"
)

type manager struct {
	d     context.Domain
	areas *list.List
}

type KVarea interface {
	set(scope.ID, interface{})
	get(scope.ID) interface{}
}

type area struct {
	root  node.Node
	x     map[scope.ID]interface{}
	ready bool
}

func area_fin(a interface{}) {
	fmt.Println("scope cleared")
}

func (a *area) set(k scope.ID, v interface{}) {
	key := scope.ID{Name: k.Name}
	a.x[key] = v
}

func (a *area) get(k scope.ID) interface{} {
	key := scope.ID{Name: k.Name}
	return a.x[key]
}

type value interface {
	set(x interface{})
	get() interface{}
}

type reference interface {
	id(...scope.ID) scope.ID
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
	setField(field string, x interface{})
	getField(field string) interface{}
	init(root node.Node)
}

func (b *basic) set(i interface{}) { b.data = i }

func (b *basic) get() interface{} { return b.data }

type ref struct {
	link object.Object
	ref  scope.ID
}

func (r *ref) id(x ...scope.ID) scope.ID {
	if len(x) == 1 {
		r.ref = x[0]
	} else if len(x) > 1 {
		panic("there can be only one")
	}
	return r.ref
}

type arr struct {
	link object.Object
	par  int64
	data []interface{}
}

func (a *arr) get(i int64) interface{} {
	if len(a.data) == 0 {
		a.data = make([]interface{}, a.par)
	}
	return a.data[i]
}

func (a *arr) set(i int64, x interface{}) {
	if len(a.data) == 0 {
		a.data = make([]interface{}, a.par)
	}
	a.data[i] = x
}

func (a *arr) sel() []interface{} { return a.data }

func (a *arr) upd(x []interface{}) { a.data = x }

type rec struct {
	link object.Object
	root node.Node
	x    map[scope.ID]interface{}
}

func (r *rec) setField(f string, x interface{}) { r.set(scope.ID{Name: f}, x) }

func (r *rec) getField(f string) interface{} { return r.get(scope.ID{Name: f}) }

func (r *rec) init(n node.Node) {
	r.root = n
	r.x = make(map[scope.ID]interface{})
}

func (a *rec) set(k scope.ID, v interface{}) { a.x[k] = v }

func (a *rec) get(k scope.ID) interface{} { return a.x[k] }

func nm() scope.Manager {
	m := &manager{areas: list.New()}
	return m
}

func init() {
	scope.NewStack = nm
	scope.Designator = design
	scope.FindObjByName = FindObjByName

	scope.NewHeap = nh
}

func design(n ...node.Node) (id scope.ID) {
	switch x := n[0].(type) {
	case node.VariableNode, node.ParameterNode:
		id = scope.ID{Name: x.Object().Name()}
	case node.FieldNode:
		fmt.Println(x.Object().Name())
		if len(n) == 1 {
			id = scope.ID{Name: x.Left().Object().Name(), Path: x.Object().Name()}
		} else if n[1] != nil {
			id = scope.ID{Name: n[1].Object().Name(), Path: x.Object().Name()}
		} else {
			panic("wrong params")
		}
	case node.IndexNode:
		id = scope.ID{Name: x.Left().Object().Name()}
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(n)))
	}
	return id
}

func odesign(o object.Object) (id scope.ID) {
	switch x := o.(type) {
	case object.VariableObject, object.ParameterObject:
		switch x.Name() {
		case "@":
			id = scope.ID{Name: "@", Ref: new(int)}
			*id.Ref = *x.Ref()[0].(*desc).Ref
		default:
			id = scope.ID{Name: x.Name()}
		}
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(o)))
	}
	return id
}

func obj(o object.Object) (key scope.ID, val interface{}) {
	switch x := o.(type) {
	case object.ConstantObject, object.ProcedureObject, object.TypeObject:
	case object.VariableObject, object.FieldObject:
		//fmt.Println(x.Name())
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
		case object.PointerType:
			val = &ref{link: o, ref: scope.ID{Ref: new(int)}}
		default:
			fmt.Println("unexpected", reflect.TypeOf(t))
		}
	case object.ParameterObject:
		//fmt.Println("'" + x.Name())
		key = scope.ID{Name: x.Name()}
		val = &ref{link: o}
	default:
		fmt.Println(reflect.TypeOf(o))
	}
	return key, val
}

func alloc(root node.Node, h KVarea, k scope.ID, v interface{}) {
	h.set(k, v)
	switch rv := v.(type) {
	case record:
		rv.init(root)
		o := rv.(*rec).link
		switch t := o.Complex().(type) {
		case object.RecordType:
			for rec := t; rec != nil; {
				for x := rec.Link(); x != nil; x = x.Link() {
					//fmt.Println(o.Name(), ".", x.Name())
					k, f := obj(x)
					alloc(root, v.(KVarea), k, f)
				}
				rec = rec.BaseType()
			}
		}
	}
}

func (m *manager) Target(...scope.Allocator) scope.Allocator {
	return m
}

func (m *manager) Allocate(n node.Node, final bool) {
	h := &area{ready: final, root: n, x: make(map[scope.ID]interface{})}
	runtime.SetFinalizer(h, area_fin)
	mod := rt_mod.DomainModule(m.Domain())
	for _, o := range mod.Objects[n] {
		k, v := obj(o)
		//fmt.Println(k, v)
		alloc(n, h, k, v)
	}
	m.areas.PushFront(h)
	fmt.Println("allocate")
}

func (m *manager) Initialize(n node.Node, par scope.PARAM) (seq frame.Sequence, ret frame.WAIT) {
	e := m.areas.Front()
	assert.For(e != nil, 20)
	h := e.Value.(*area)
	assert.For(h.root == n, 21)
	assert.For(!h.ready, 22)
	val := par.Values
	fmt.Println("initialize")
	f := par.Frame
	end := func(frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		h.ready = true
		if par.Tail != nil {
			return par.Tail(f)
		} else {
			return frame.End()
		}
	}
	seq = end
	ret = frame.NOW
	for next := par.Objects; next != nil; next = next.Link() {
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
				k, v := scope.ID{Name: next.Name()}, &basic{link: next}
				h.set(k, v)
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
				h.get(scope.ID{Name: next.Name()}).(*ref).ref = design(ov)
			default:
				panic("unknown value")
			}
		case node.DerefNode:
			rt2.Push(rt2.New(ov), f)
			dn := next
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				fmt.Println(rt2.DataOf(f)[ov])
				switch dn.(type) {
				case object.VariableObject:
					m.Update(odesign(dn), func(old interface{}) interface{} {
						return rt2.DataOf(f)[ov]
					})
				case object.ParameterObject:
					h.get(scope.ID{Name: dn.Name()}).(*ref).ref = odesign(rt2.DataOf(f)[ov].(object.Object))
					fmt.Println("got param", odesign(rt2.DataOf(f)[ov].(object.Object)))
				default:
					panic(fmt.Sprintln("unknown value", reflect.TypeOf(next)))
				}
				return end, frame.NOW
			}
			ret = frame.LATER
		default:
			panic(fmt.Sprintln("unknown value", reflect.TypeOf(val)))
		}
		val = val.Link()
	}
	return seq, ret
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

func (m *manager) Select(i scope.ID) interface{} {
	type result struct {
		x interface{}
	}
	var res *result
	var sel func(interface{}) *result
	fmt.Println(i)
	sel = func(x interface{}) (ret *result) {
		switch x := x.(type) {
		case value:
			ret = &result{x: x.get()}
		case reference:
			i = x.id()
			fmt.Println("ref!", i)
			ret = nil
		case array:
			if i.Index != nil {
				ret = &result{x: x.get(*i.Index)}
			} else {
				ret = &result{x: x.sel()}
			}
		case record:
			if i.Path == "" {
				ret = &result{x: x.(*rec).link}
			} else {
				z := x.getField(i.Path)
				ret = sel(z)
			}
		case nil:
			if i.Name == "@" {
				fmt.Println("ptr")
				if hm := m.Domain().Discover(context.HEAP).(scope.Manager); hm != nil {
					ret = &result{x: hm.Select(i)}
				} else {
					panic(0)
				}
			}
		default:
			panic(0)
		}
		return ret
	}
	for e := m.areas.Front(); (e != nil) && (res == nil); e = e.Next() {
		h := e.Value.(*area)
		if h.ready {
			res = sel(h.get(i))
		}
	}
	assert.For(res != nil, 40)
	//fmt.Println(res.x)
	return res.x
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
	case int32:
		fmt.Println("not an array")
		return []interface{}{rune(0)}
	default:
		panic(fmt.Sprintln("unsupported", reflect.TypeOf(x)))
	}
}

func (m *manager) Update(i scope.ID, val scope.ValueFor) {
	assert.For(val != nil, 21)
	fmt.Println("update", i)
	var x interface{}
	var upd func(x interface{}) (ret interface{})
	upd = func(x interface{}) (ret interface{}) {
		fmt.Println(reflect.TypeOf(x))
		switch x := x.(type) {
		case value:
			old := x.get()
			tmp := val(old)
			assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
			//fmt.Println(tmp)
			x.set(tmp)
			ret = x
		case reference:
			fmt.Println(x.id())
			switch {
			case x.id().Ref != nil && *x.id().Ref == 0: //это нулевой указатель
				tmp := val(nil)
				ret = x
				switch id := tmp.(type) {
				case scope.ID:
					x.id(id)
				case object.VariableObject:
					x.id(odesign(id))
				default:
					panic(fmt.Sprintln("only id for nil pointer", reflect.TypeOf(id)))
				}
			case x.id().Ref != nil && *x.id().Ref != 0: //это указатель на объект в куче
				i.Name = x.id().Name
				i.Ref = x.id().Ref
				ret = nil
			default: //это параметр процедуры
				i.Name = x.id().Name
				ret = nil
			}
		case array:
			if i.Index != nil {
				old := x.get(*i.Index)
				tmp := val(old)
				assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
				//fmt.Println(tmp)
				x.set(*i.Index, tmp)
			} else {
				old := x.sel()
				tmp := val(old)
				assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
				//fmt.Println(tmp)
				x.upd(arrConv(tmp))
			}
			ret = x
		case record:
			if i.Path == "" {
				//fmt.Println(i, depth)
				panic(0) //случай выбора всей записи целиком
			} else {
				z := x.getField(i.Path)
				ret = upd(z)
			}
		case nil:
			if i.Name == "@" {
				fmt.Println("ptr")
				if hm := m.Domain().Discover(context.HEAP).(scope.Manager); hm != nil {
					hm.Update(i, val)
					ret = hm
				} else {
					panic(0)
				}
			} else {
				ret = x
			}
		default:
			panic(fmt.Sprintln("unhandled", reflect.TypeOf(x)))
		}
		return ret
	}
	for e := m.areas.Front(); (e != nil) && (x == nil); e = e.Next() {
		h := e.Value.(*area)
		x = upd(h.get(i))
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
		x := h.get(scope.ID{Name: name})
		switch x.(type) {
		case *basic:
			ret = x.(*basic).link
		default:
			//fmt.Println("no such object")
		}
	}
	return ret
}
