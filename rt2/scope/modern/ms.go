package modern

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant/enter"
	cpm "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

type level struct {
	k     map[cp.ID]int
	v     map[int]scope.Variable
	r     map[int]scope.Ref
	next  int
	ready bool
}

type area struct {
	d    context.Domain
	data []*level
}

type salloc struct {
	area *area
}

type ref struct {
	id   cp.ID
	link object.Object
}

func (r *ref) String() string {
	return fmt.Sprint(r.link.Name(), "@", r.id)
}

func newRef(x object.Object) *ref {
	return &ref{link: x}
}

func (a *area) allocate(mod *cpm.Module, n node.EnterNode, r bool) {
	ol := mod.Objects[n]
	l := &level{next: 1, ready: r,
		k: make(map[cp.ID]int),
		v: make(map[int]scope.Variable),
		r: make(map[int]scope.Ref)}
	a.data = append(a.data, l)
	skip := make(map[cp.ID]interface{}) //для процедурных типов в общей куче могут валяться переменные, скипаем их
	tl := mod.Types[n]
	for _, t := range tl {
		switch x := t.(type) {
		case object.BasicType:
			for link := x.Link(); link != nil; link = link.Link() {
				skip[link.Adr()] = link
			}
		case object.RecordType:
			for link := x.Link(); link != nil; link = link.Link() {
				skip[link.Adr()] = link
			}
		}
	}
	for _, o := range ol {
		imp := mod.ImportOf(o)
		fmt.Println(reflect.TypeOf(o), o.Adr())
		if imp == "" && skip[o.Adr()] == nil {
			fmt.Println(l.next)
			switch x := o.(type) {
			case object.VariableObject:
				switch t := o.Complex().(type) {
				case nil, object.BasicType, object.ArrayType, object.DynArrayType:
					l.v[l.next] = NewData(x)
					l.k[x.Adr()] = l.next
					l.next++
				case object.RecordType:
					l.v[l.next] = newRec(x)
					l.k[x.Adr()] = l.next
					l.next++
				case object.PointerType:
					fmt.Println("pointer")
				default:
					halt.As(20, reflect.TypeOf(t))
				}
			case object.TypeObject:
				//do nothing
			case object.ConstantObject:
				//do nothing
			case object.ProcedureObject:
				//do nothing
			case object.ParameterObject:
				if n.Enter() == enter.PROCEDURE {
					l.r[l.next] = newRef(x)
					l.k[x.Adr()] = l.next
					l.next++
				}
			default:
				halt.As(20, reflect.TypeOf(x))
			}
		}
	}
}

func (a *area) top() *level {
	if len(a.data) > 0 {
		return a.data[len(a.data)-1]
	}
	return nil
}

func (a *area) Provide(x interface{}) scope.ValueFor {
	return func(scope.Value) scope.Value {
		switch z := x.(type) {
		case node.ConstantNode:
			return NewConst(z)
		case object.ProcedureObject:
			return NewProc(z)
		default:
			halt.As(100, reflect.TypeOf(z))
		}
		panic(0)
	}
}

func (a *salloc) Allocate(n node.Node, final bool) {
	fmt.Println("ALLOCATE")
	mod := rtm.DomainModule(a.area.d)
	a.area.allocate(mod, n.(node.EnterNode), final)
}

func (a *salloc) Dispose(n node.Node) {
	x := a.area.data
	a.area.data = nil
	for i := 0; i < len(x)-1; i++ {
		a.area.data = append(a.area.data, x[i])
	}
}

func (a *salloc) Initialize(n node.Node, par scope.PARAM) (seq frame.Sequence, ret frame.WAIT) {
	fmt.Println("INITIALIZE")
	l := a.area.top()
	assert.For(l != nil && !l.ready, 20)
	val := par.Values
	f := par.Frame
	end := func(frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		l.ready = true
		if par.Tail != nil {
			return par.Tail(f)
		} else {
			return frame.End()
		}
	}
	seq = end
	ret = frame.NOW
	for next := par.Objects; next != nil; next = next.Link() {
		switch o := next.(type) {
		case object.VariableObject:
			switch nv := val.(type) {
			case node.ConstantNode:
				v := NewConst(nv)
				l.v[l.k[o.Adr()]].Set(v)
			case node.VariableNode:
				v := a.area.Select(nv.Object().Adr())
				l.v[l.k[o.Adr()]].Set(v)
			default:
				halt.As(40, reflect.TypeOf(nv))
			}
		case object.ParameterObject:
			switch nv := val.(type) {
			case node.VariableNode:
				old := l.r[l.k[o.Adr()]].(*ref)
				l.r[l.k[o.Adr()]] = &ref{link: old.link, id: nv.Object().Adr()}
			case node.ConstantNode: //array :) заменяем ссылку на переменную
				old := l.r[l.k[o.Adr()]].(*ref)
				l.r[l.k[o.Adr()]] = nil
				data := NewConst(nv)
				switch data.(type) {
				case STRING, SHORTSTRING:
					val := &dynarr{link: old.link}
					val.Set(data)
					l.v[l.k[o.Adr()]] = val
				default:
					halt.As(100, reflect.TypeOf(data))
				}
			default:
				halt.As(40, reflect.TypeOf(nv))
			}
		default:
			halt.As(40, reflect.TypeOf(o))
		}
		val = val.Link()
	}
	return seq, ret
}

func (a *area) Update(id cp.ID, fval scope.ValueFor) {
	assert.For(id != 0, 20)
	fmt.Println("UPDATE", id)
	var upd func(x int, id cp.ID)
	var k int
	upd = func(x int, id cp.ID) {
		for i := x - 1; i >= 0 && k == 0; i-- {
			l := a.data[i]
			if l.ready {
				k = l.k[id]
				if k != 0 {
					v := a.data[i].v[k]
					if v == nil { //ref?
						r := l.r[k]
						if r != nil {
							upd(i, r.(*ref).id)
							break
						}
					} else {
						v.Set(fval(a.data[i].v[k]))
					}
				}
			}
		}
	}
	k = 0
	upd(len(a.data), id)
	assert.For(k != 0, 60)
}
func (a *area) Select(id cp.ID) (ret scope.Value) {
	fmt.Println("SELECT", id)
	var sel func(x int, id cp.ID)
	sel = func(x int, id cp.ID) {
		for i := x - 1; i >= 0 && ret == nil; i-- {
			l := a.data[i]
			k := 0
			if l.ready {
				k = l.k[id]
				if k != 0 {
					ret = l.v[k]
					if ret == nil { //ref?
						r := l.r[k]
						if r != nil {
							sel(i, r.(*ref).id)
							break
						}
					}
				}
			}
		}
	}
	sel(len(a.data), id)
	assert.For(ret != nil, 60)
	return ret
}

func (a *area) Target(...scope.Allocator) scope.Allocator { return &salloc{area: a} }
func (a *area) String() (ret string) {
	for _, l := range a.data {
		ret = fmt.Sprintln(ret, l)
	}
	return ret
}

func (l *level) String() (ret string) {
	for k, v := range l.k {
		ret = fmt.Sprint(ret, "@", k, v, l.v[v])
		if l.v[v] == nil {
			ret = fmt.Sprintln(ret, l.r[v])
		} else {
			ret = fmt.Sprintln(ret)
		}
	}
	return ret
}

func (a *area) Init(d context.Domain) { a.d = d }

func (a *area) Domain() context.Domain { return a.d }

func (a *area) Handle(msg interface{}) {}

func fn(mgr scope.Manager, name string) (ret object.Object) {
	fmt.Println("FIND", name)
	a, ok := mgr.(*area)
	assert.For(ok, 20)
	assert.For(name != "", 21)
	for i := len(a.data) - 1; i >= 0 && ret == nil; i-- {
		l := a.data[i]
		for _, v := range l.v {
			switch vv := v.(type) {
			case *data:
				fmt.Println(vv.link.Name())
				if vv.link.Name() == name {
					ret = vv.link
				}
			default:
				fmt.Println(reflect.TypeOf(vv))
			}
		}
	}
	return ret
}

func nn() scope.Manager {
	return &area{}
}

func init() {
	scope.New = nn
	scope.FindObjByName = fn
}
