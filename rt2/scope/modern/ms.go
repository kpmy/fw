package modern

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant/enter"
	cpm "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"

	"ypk/assert"
	"ypk/halt"
)

type level struct {
	k      map[cp.ID]int
	v      map[int]scope.Variable
	r      map[int]scope.Ref
	l      map[int]*level
	next   int
	ready  bool
	nested bool
}

type area struct {
	d    context.Domain
	data []*level
	all  scope.Allocator
}

type salloc struct {
	area *area
}

type ref struct {
	scope.Ref
	id   cp.ID
	link object.Object
	sc   scope.Manager
}

func (r *ref) String() string {
	var m string
	if r.sc != nil {
		m = rtm.DomainModule(r.sc.Domain()).Name
	}
	return fmt.Sprint(m, " ", r.link.Name(), "@", r.id)
}

func newRef(x object.Object) *ref {
	return &ref{link: x}
}

func newlvl() *level {
	return &level{next: 1,
		k: make(map[cp.ID]int), ready: true,
		v: make(map[int]scope.Variable),
		r: make(map[int]scope.Ref),
		l: make(map[int]*level)}
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
			return newConst(z)
		case object.ProcedureObject:
			return newProc(z)
		default:
			halt.As(100, reflect.TypeOf(z))
		}
		panic(0)
	}
}

//var alloc func(*level, []object.Object, map[cp.ID]interface{})
func (l *level) alloc(mod *cpm.Module, root node.Node, ol []object.Object, skip map[cp.ID]interface{}) {
	for _, o := range ol {
		imp := mod.ImportOf(o)
		utils.PrintScope(reflect.TypeOf(o), o.Adr())
		_, field := o.(object.FieldObject)
		if imp == "" && (skip[o.Adr()] == nil || (field && l.nested)) {
			utils.PrintScope("next", l.next)
			switch x := o.(type) {
			case object.VariableObject, object.FieldObject:
				switch t := o.Complex().(type) {
				case nil, object.BasicType:
					l.v[l.next] = newData(x)
					l.k[x.Adr()] = l.next
					l.next++
				case object.ArrayType, object.DynArrayType:
					l.v[l.next] = newData(x)
					l.k[x.Adr()] = l.next
					l.next++
				case object.RecordType:
					l.v[l.next] = newRec(x)
					nl := newlvl()
					nl.nested = true
					l.l[l.next] = nl
					l.k[x.Adr()] = l.next
					fl := make([]object.Object, 0)
					for rec := t; rec != nil; {
						for x := rec.Link(); x != nil; x = x.Link() {
							//fmt.Println(o.Name(), ".", x.Name(), x.Adr())
							fl = append(fl, x)
						}
						rec = rec.BaseType()
					}
					//fmt.Println("record")
					l.v[l.next].(*rec).l = nl
					nl.alloc(mod, root, fl, skip)
					l.next++
				case object.PointerType:
					l.v[l.next] = newPtr(x)
					l.k[x.Adr()] = l.next
					l.next++
				default:
					halt.As(20, reflect.TypeOf(t))
				}
			case object.TypeObject, object.ConstantObject, object.ProcedureObject, object.Module:
				//do nothing
			case object.ParameterObject:
				if root.(node.EnterNode).Enter() == enter.PROCEDURE {
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

func (a *salloc) Allocate(n node.Node, final bool) {
	mod := rtm.DomainModule(a.area.d)
	utils.PrintScope("ALLOCATE FOR", mod.Name, n.Adr())
	tl := mod.Types[n]
	skip := make(map[cp.ID]interface{}) //для процедурных типов в общей куче могут валяться переменные, скипаем их
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
	ol := mod.Objects[n]
	switch o := n.Object().(type) {
	case object.ProcedureObject:
		for l := o.Link(); l != nil; l = l.Link() {
			ol = append(ol, l)
		}
	case nil:
		for _, o := range ol {
			switch t := o.(type) {
			case object.ProcedureObject:
				for l := t.Link(); l != nil; l = l.Link() {
					skip[l.Adr()] = l
				}
			}
		}
	default:
		halt.As(100, reflect.TypeOf(o))

	}
	nl := newlvl()
	nl.ready = final
	a.area.data = append(a.area.data, nl)
	nl.alloc(mod, n, ol, skip)
}

func (a *salloc) Dispose(n node.Node) {
	x := a.area.data
	old := x[len(x)-1]
	old.k = nil
	old.v = nil
	old.l = nil
	old.r = nil
	a.area.data = nil
	for i := 0; i < len(x)-1; i++ {
		a.area.data = append(a.area.data, x[i])
	}
}

func (a *salloc) Initialize(n node.Node, par scope.PARAM) (seq frame.Sequence, ret frame.WAIT) {
	utils.PrintScope("INITIALIZE")
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
	var sm scope.Manager
	for next := par.Objects; next != nil; next = next.Link() {
		global := f.Domain().Discover(context.UNIVERSE).(context.Domain)
		mod := rtm.ModuleOfNode(f.Domain(), val)
		//mod := rtm.ModuleOfNode(f.Domain(), val.Object())
		if mod != nil {
			//fmt.Println(mod.Name)
			global = global.Discover(mod.Name).(context.Domain)
			sm = global.Discover(context.SCOPE).(scope.Manager)
		} else { //для фиктивных узлов, которые созданы рантаймом, типа INC/DEC
			sm = a.area
		}
		switch o := next.(type) {
		case object.VariableObject:
			switch nv := val.(type) {
			case node.ConstantNode:
				v := newConst(nv)
				l.v[l.k[o.Adr()]].Set(v)
			case node.VariableNode, node.ParameterNode:
				v := sm.Select(nv.Object().Adr())
				l.v[l.k[o.Adr()]].Set(v)
			case node.OperationNode:
				nf := rt2.New(nv)
				rt2.Push(nf, f)
				rt2.Assert(f, func(f frame.Frame) (bool, int) {
					return rt2.ValueOf(f)[nv.Adr()] != nil, 59
				})
				rt2.ReplaceDomain(nf, global)
				seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					v := rt2.ValueOf(f)[nv.Adr()]
					l.v[l.k[o.Adr()]].Set(v)
					return end, frame.NOW
				}
				ret = frame.LATER
			case node.FieldNode, node.DerefNode:
				nf := rt2.New(nv)
				rt2.Push(nf, f)
				rt2.Assert(f, func(f frame.Frame) (bool, int) {
					return rt2.ValueOf(f)[nv.Adr()] != nil, 60
				})
				rt2.ReplaceDomain(nf, global)
				seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					v := rt2.ValueOf(f)[nv.Adr()]
					l.v[l.k[o.Adr()]].Set(v)
					return end, frame.NOW
				}
				ret = frame.LATER
			default:
				halt.As(40, reflect.TypeOf(nv))
			}
		case object.ParameterObject:
			switch nv := val.(type) {
			case node.VariableNode:
				old := l.r[l.k[o.Adr()]].(*ref)
				l.r[l.k[o.Adr()]] = &ref{link: old.link, sc: sm, id: nv.Object().Adr()}
			case node.ConstantNode: //array :) заменяем ссылку на переменную
				old := l.r[l.k[o.Adr()]].(*ref)
				l.r[l.k[o.Adr()]] = nil
				data := newConst(nv)
				switch data.(type) {
				case STRING, SHORTSTRING:
					val := &dynarr{link: old.link}
					val.Set(data)
					l.v[l.k[o.Adr()]] = val
				default:
					halt.As(100, reflect.TypeOf(data))
				}
			case node.DerefNode:
				rt2.Push(rt2.New(nv), f)
				rt2.Assert(f, func(f frame.Frame) (bool, int) {
					return rt2.ValueOf(f)[nv.Adr()] != nil, 61
				})
				dn := next
				old := l.r[l.k[dn.Adr()]].(*ref)
				seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					switch dn.(type) {
					case object.VariableObject, object.ParameterObject:
						l.r[l.k[dn.Adr()]] = nil
						data := rt2.ValueOf(f)[nv.Adr()]
						switch deref := data.(type) {
						case STRING, SHORTSTRING:
							val := &dynarr{link: old.link}
							val.Set(deref)
							l.v[l.k[dn.Adr()]] = val
						case *rec:
							l.v[l.k[dn.Adr()]] = deref
						default:
							halt.As(100, reflect.TypeOf(data))
						}
					default:
						panic(fmt.Sprintln("unknown value", reflect.TypeOf(next)))
					}
					return end, frame.NOW
				}
				ret = frame.LATER
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

func (a *salloc) Join(m scope.Manager) { a.area = m.(*area) }

func (a *area) Update(id cp.ID, fval scope.ValueFor) {
	assert.For(id != 0, 20)
	var upd func(x int, id cp.ID)
	var k int
	upd = func(x int, id cp.ID) {
		utils.PrintScope("UPDATE", id)
		for i := x - 1; i >= 0 && k == 0; i-- {
			l := a.data[i]
			if l.ready {
				k = l.k[id]
				if k != 0 {
					v := a.data[i].v[k]
					if v == nil { //ref?
						r := l.r[k]
						if r != nil {
							utils.PrintScope("ref")
							if r.(*ref).sc == a {
								upd(i, r.(*ref).id)
							} else {
								k = -1
								r.(*ref).sc.Update(r.(*ref).id, fval)
							}
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

func (a *area) Select(id cp.ID, val ...scope.ValueOf) (ret scope.Value) {
	var sel func(x int, id cp.ID)
	sel = func(x int, id cp.ID) {
		utils.PrintScope("SELECT", id)
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
							if l.l[k] != nil { //rec
								panic(0)
							} else {
								utils.PrintScope("ref")
								if r.(*ref).sc == a {
									sel(i, r.(*ref).id)
								} else {
									ret = r.(*ref).sc.Select(r.(*ref).id)
								}
							}
							break
						}
					} else if len(val) > 0 {
						val[0](ret)
					}
				}
			}
		}
	}
	sel(len(a.data), id)
	assert.For(ret != nil, 60)
	return ret
}

func (a *area) Target(all ...scope.Allocator) scope.Allocator {
	if len(all) > 0 {
		a.all = all[0]
	}
	if a.all == nil {
		return &salloc{area: a}
	} else {
		a.all.Join(a)
		return a.all
	}
}

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
		} else if l.l[v] != nil {
			ret = fmt.Sprintln(ret, "{")
			ret = fmt.Sprintln(ret, l.l[v], "}")
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
	utils.PrintScope("FIND", name)
	a, ok := mgr.(*area)
	assert.For(ok, 20)
	assert.For(name != "", 21)
	for i := len(a.data) - 1; i >= 0 && ret == nil; i-- {
		l := a.data[i]
		for _, v := range l.v {
			switch vv := v.(type) {
			case *data:
				utils.PrintScope(vv.link.Name())
				if vv.link.Name() == name {
					ret = vv.link
				}
			default:
				utils.PrintScope(reflect.TypeOf(vv))
			}
		}
	}
	return ret
}

func nn(role string) scope.Manager {
	if role == context.SCOPE {
		return &area{all: &salloc{}}
	} else if role == context.HEAP {
		return &area{all: &halloc{}}
	} else {
		panic(0)
	}
}

func init() {
	scope.New = nn
	scope.FindObjByName = fn
}
