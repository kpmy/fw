package data

import (
	"fmt"
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/rules2/wrap/data/items"
	"fw/rt2/rules2/wrap/eval"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

type area struct {
	d   context.Domain
	all scope.Allocator
	il  items.Data
}

type salloc struct {
	area *area
}

type key struct {
	items.Key
	id cp.ID
}

func (k *key) String() string {
	return fmt.Sprint(k.id)
}

func (k *key) EqualTo(to items.Key) int {
	kk, ok := to.(*key)
	if ok && kk.id == k.id {
		return 0
	} else {
		return -1
	}
}

type item struct {
	items.Item
	k items.Key
	d interface{}
}

func (i *item) KeyOf(k ...items.Key) items.Key {
	if len(k) == 1 {
		i.k = k[0]
	}
	return i.k
}

func (i *item) Copy(from items.Item) { panic(0) }

func (i *item) Data(d ...interface{}) interface{} {
	if len(d) == 1 {
		i.d = d[0]
	}
	return i.d
}

func (i *item) Value() scope.Value {
	return i.d.(scope.Value)
}

func (a *area) Select(this cp.ID, val scope.ValueOf) {
	utils.PrintScope("SELECT", this)
	d, ok := a.il.Get(&key{id: this}).(*item)
	assert.For(ok, 20, this)
	val(d.Value())
}

func (a *area) Exists(this cp.ID) bool {
	utils.PrintScope("SEARCH", this)
	return a.il.Exists(&key{id: this})
}

func (a *salloc) push(_o object.Object) {
	switch o := _o.(type) {
	case object.VariableObject:
		switch t := o.Complex().(type) {
		case nil, object.BasicType:
			x := newData(o)
			d := &item{}
			d.Data(x)
			a.area.il.Set(&key{id: o.Adr()}, d)
		case object.ArrayType, object.DynArrayType:
			x := newData(o)
			d := &item{}
			d.Data(x)
			a.area.il.Set(&key{id: o.Adr()}, d)
		default:
			halt.As(100, reflect.TypeOf(t))
		}
	case object.ParameterObject:
		a.area.il.Hold(&key{id: o.Adr()})
	default:
		halt.As(100, reflect.TypeOf(o))
	}
}

func (a *salloc) Allocate(n node.Node, final bool) {
	mod := rtm.ModuleOfNode(a.area.d, n)
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
	//все объекты скоупа
	ol := mod.Objects[n]
	//добавим либо переменные внутри процедуры либо если мы создаем скоуп для модуля то процедурные объекты добавим в скиплист
	switch o := n.Object().(type) {
	case object.ProcedureObject:
		for l := o.Link(); l != nil; l = l.Link() {
			ol = append(ol, l)
		}
	case nil: //do nothing
	default:
		halt.As(100, reflect.TypeOf(o))
	}

	for _, o := range ol {
		switch t := o.(type) {
		case object.ProcedureObject:
			for l := t.Link(); l != nil; l = l.Link() {
				skip[l.Adr()] = l
			}
			skip[o.Adr()] = o
		case object.ConstantObject:
			skip[o.Adr()] = o
		}
	}
	a.area.il.Begin()
	for _, o := range ol {
		if skip[o.Adr()] == nil {
			utils.PrintScope(o.Adr(), o.Name())
			a.push(o)
		}
	}
	if final {
		a.area.il.End()
	}
}

func (a *salloc) Dispose(n node.Node) {
	a.area.il.Drop()
}

func (a *salloc) proper_init(root node.Node, _val node.Node, _par object.Object, tail eval.Do, in eval.IN) eval.Do {
	utils.PrintScope("INITIALIZE")
	const link = "initialize:par"
	end := func(in eval.IN) eval.OUT {
		a.area.il.End()
		return eval.Later(tail)
	}
	var next eval.Do
	do := func(val node.Node, par object.Object) (out eval.OUT) {
		utils.PrintScope(par.Adr(), par.Name(), ":=", reflect.TypeOf(val))
		out = eval.Now(next)
		switch par.(type) {
		case object.VariableObject:
			out = eval.GetExpression(in, link, val, func(in eval.IN) eval.OUT {
				it := a.area.il.Get(&key{id: par.Adr()}, items.INIT).(*item)
				v := it.Value().(scope.Variable)
				val := rt2.ValueOf(in.Frame)[eval.KeyOf(in, link)]
				v.Set(val)
				return eval.Later(next)
			})
		case object.ParameterObject:
			switch val.(type) {
			case node.Designator:
				out = eval.GetDesignator(in, link, val, func(in eval.IN) eval.OUT {
					mt := rt2.RegOf(in.Frame)[context.META].(*eval.Meta)
					fa := mt.Scope.(*area).il
					a.area.il.Link(&key{id: par.Adr()}, items.ID{In: fa, This: &key{id: mt.Id}})
					return eval.Later(next)
				})
			case node.Expression: //array заменяем ссылку на переменную
				out = eval.GetExpression(in, link, val, func(in eval.IN) eval.OUT {
					d := &item{}
					data := rt2.ValueOf(in.Frame)[eval.KeyOf(in, link)]
					switch data.(type) {
					case STRING, SHORTSTRING:
						val := &dynarr{link: par}
						val.Set(data)
						d.Data(val)
					default:
						halt.As(100, reflect.TypeOf(data))
					}
					a.area.il.Put(&key{id: par.Adr()}, d)
					return eval.Later(next)
				})
			default:
				halt.As(100, reflect.TypeOf(val))
			}
		default:
			halt.As(100, reflect.TypeOf(par))
		}
		return
	}
	val := _val
	par := _par
	next = func(eval.IN) eval.OUT {
		if val == nil {
			return eval.Later(end)
		} else {
			step := do(val, par)
			val = val.Link()
			par = par.Link()
			return step
		}
	}
	return next
}

func (a *salloc) Initialize(n node.Node, par scope.PARAM) (frame.Sequence, frame.WAIT) {
	var tail eval.Do
	if par.Tail != nil {
		tail = eval.Expose(par.Tail)
	} else {
		tail = eval.Tail(eval.STOP)
	}
	return eval.Propose(a.proper_init(n, par.Values, par.Objects, tail, eval.IN{Frame: par.Frame, Parent: par.Frame.Parent()})), frame.NOW
}

func (a *salloc) Join(m scope.Manager) { a.area = m.(*area) }

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

func (a *area) String() string { return "fixme" }

func (a *area) Provide(x interface{}) scope.Value {
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

func (a *area) Init(d context.Domain) { a.d = d }

func (a *area) Domain() context.Domain { return a.d }

func (a *area) Handle(msg interface{}) {}

func nn(role string) scope.Manager {
	switch role {
	case context.SCOPE, context.CALL:
		return &area{all: &salloc{}, il: items.New()}
	case context.HEAP:
		return &area{all: nil}
		//return &area{all: &halloc{}}
	default:
		panic(0)
	}
}

func fn(mgr scope.Manager, name string) (ret object.Object) {
	utils.PrintScope("FIND", name)
	a, ok := mgr.(*area)
	assert.For(ok, 20)
	assert.For(name != "", 21)
	a.il.ForEach(func(in items.Value) (ok bool) {
		var v scope.Value
		switch val := in.(type) {
		case *item:
			v = val.Value()
		}
		switch vv := v.(type) {
		case *data:
			utils.PrintScope(vv.link.Name())
			if vv.link.Name() == name {
				ret = vv.link
				ok = true
			}
		default:
			utils.PrintScope(reflect.TypeOf(vv))
		}
		return
	})
	return ret
}

func init() {
	scope.New = nn
	scope.FindObjByName = fn
}
