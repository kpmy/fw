package data

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
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
	"runtime"
	"ypk/assert"
	"ypk/halt"
)

type area struct {
	d      context.Domain
	all    scope.Allocator
	il     items.Data
	unsafe bool
}

type salloc struct {
	area *area
}

type halloc struct {
	area *area
}

type key struct {
	items.Key
	id cp.ID
}

func (k *key) String() string {
	return fmt.Sprint(k.id)
}

func (k *key) Hash() int {
	return int(k.id)
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

func push(dom context.Domain, il items.Data, _o object.Object) {
	switch o := _o.(type) {
	case object.VariableObject, object.FieldObject:
		var x interface{}
		switch t := o.Complex().(type) {
		case nil, object.BasicType:
			x = newData(o)
		case object.ArrayType, object.DynArrayType:
			x = newData(o)
		case object.RecordType:
			ml := dom.Global().Discover(context.MOD).(rtm.List)
			r := newRec(o)
			x = r
			fl := make([]object.Object, 0)
			for rec := t; rec != nil; {
				for x := rec.Link(); x != nil; x = x.Link() {
					switch x.(type) {
					case object.FieldObject:
						//fmt.Println(o.Name(), ".", x.Name(), x.Adr())
						fl = append(fl, x)
					case object.ParameterObject, object.ProcedureObject, object.VariableObject:
						//do nothing
					default:
						halt.As(100, reflect.TypeOf(x))
					}
				}
				if rec.BaseRec() == nil {
					x := ml.NewTypeCalc()
					x.ConnectTo(rec)
					_, frec := x.ForeignBase()
					//fmt.Println(frec)
					rec, _ = frec.(object.RecordType)
				} else {
					rec = rec.BaseRec()
				}
			}
			r.fi = items.New()
			r.fi.Begin()
			for _, f := range fl {
				push(dom, r.fi, f)
			}
			r.fi.End()
		case object.PointerType:
			x = newPtr(o)
		default:
			halt.As(100, reflect.TypeOf(t))
		}
		assert.For(x != nil, 40)
		//fmt.Println(_o.Name(), x)
		d := &item{}
		d.Data(x)
		il.Set(&key{id: o.Adr()}, d)
	case object.ParameterObject:
		il.Hold(&key{id: o.Adr()})
	default:
		halt.As(100, reflect.TypeOf(o))
	}
}

func fin(x interface{}) {
	switch p := x.(type) {
	case *ptrValue:
		defer func() {
			mod := rtm.ModuleOfType(p.scope.Domain(), p.ct)
			ol := mod.Objects[mod.Enter]
			var fn object.ProcedureObject
			for _, _po := range ol {
				switch po := _po.(type) {
				case object.ProcedureObject:
					if po.Name() == "FINALIZE" && po.Link().Complex().Equals(p.ct) {
						fn = po
						break
					}
				}
			}
			if fn != nil {
				global := p.scope.Domain().Discover(context.UNIVERSE).(context.Domain)
				root := global.Discover(context.STACK).(frame.Stack)
				cn := node.New(constant.CALL, cp.Some())
				ol := mod.NodeByObject(fn)
				assert.For(len(ol) <= 1, 40)
				cn.SetLeft(ol[0])
				cc := node.New(constant.CONSTANT, cp.Some()).(node.ConstantNode)
				cc.SetData(p)
				cc.SetType(object.COMPLEX)
				cn.SetRight(cc)
				nf := rt2.New(cn)
				nf.Init(global.Discover(mod.Name).(context.Domain))
				root.Queue(nf)
			}
			p.scope.Target().(scope.HeapAllocator).Dispose(p.id)
		}()
	}
}

func (h *halloc) Allocate(name string, t object.PointerType, par ...interface{}) scope.Value {
	utils.PrintScope("HEAP ALLOCATE")
	//mod := rtm.ModuleOfType(h.area.d, t)
	assert.For(t != nil, 20)
	var res scope.Value
	var talloc func(t object.PointerType)
	talloc = func(t object.PointerType) {
		switch bt := t.Complex().(type) {
		case object.RecordType:
			fake := object.New(object.VARIABLE, cp.Some())
			fake.SetComplex(bt)
			fake.SetType(object.COMPLEX)
			fake.SetName("{" + "}")
			push(h.area.d, h.area.il, fake)
			res = &ptrValue{scope: h.area, id: fake.Adr(), ct: t}
		case object.DynArrayType:
			assert.For(len(par) > 0, 20)
			fake := object.New(object.VARIABLE, cp.Some())
			fake.SetComplex(bt)
			fake.SetType(object.COMPLEX)
			fake.SetName("[]")
			push(h.area.d, h.area.il, fake)
			h.area.Select(fake.Adr(), func(v scope.Value) {
				arr, ok := v.(*dynarr)
				assert.For(ok, 60)
				arr.Set(par[0].(scope.Value))
			})
			res = &ptrValue{scope: h.area, id: fake.Adr(), ct: t}
		default:
			halt.As(100, fmt.Sprintln("cannot allocate", reflect.TypeOf(bt)))
		}
	}
	talloc(t)
	assert.For(res != nil, 60)
	runtime.SetFinalizer(res, fin)
	return res
}

func (h *halloc) Dispose(i cp.ID) {
	h.area.Select(i, func(v scope.Value) {
		utils.PrintScope("dispose", v)
		h.area.il.Remove(&key{id: i})
	})
}

func (a *halloc) Join(m scope.Manager) { a.area = m.(*area) }

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
		case object.TypeObject:
			skip[o.Adr()] = o
		case object.ConstantObject:
			skip[o.Adr()] = o
		case object.Module:
			skip[o.Adr()] = o
		}
	}
	a.area.il.Begin()
	a.area.unsafe = true
	for _, o := range ol {
		if skip[o.Adr()] == nil {
			utils.PrintScope(o.Adr(), o.Name())
			push(a.area.d, a.area.il, o)
		}
	}
	if final {
		a.area.il.End()
		a.area.unsafe = false
	}
}

func (a *salloc) Dispose(n node.Node) {
	utils.PrintScope("DISPOSE")
	a.area.il.Drop()
}

func (a *salloc) proper_init(root node.Node, _val node.Node, _par object.Object, tail eval.Do, in eval.IN) eval.Do {
	utils.PrintScope("INITIALIZE")
	const link = "initialize:par"
	end := func(in eval.IN) eval.OUT {
		a.area.il.End()
		a.area.unsafe = false
		return eval.Later(tail)
	}
	var next eval.Do
	do := func(val node.Node, par object.Object) (out eval.OUT) {
		utils.PrintScope(par.Adr(), par.Name(), ":=", reflect.TypeOf(val))
		out = eval.Now(next)
		switch par.(type) {
		case object.VariableObject:
			out = eval.GetExpression(in, link, val, func(in eval.IN) eval.OUT {
				it := a.area.il.Get(&key{id: par.Adr()}).(*item)
				v := it.Value().(scope.Variable)
				val := rt2.ValueOf(in.Frame)[eval.KeyOf(in, link)]
				v.Set(val)
				return eval.Later(next)
			})
		case object.ParameterObject:
			switch val.(type) {
			case node.Designator:
				out = eval.GetDesignator(in, link, val, func(in eval.IN) eval.OUT {
					if mt, _ := rt2.RegOf(in.Frame)[context.META].(*eval.Meta); mt != nil && mt.Scope != nil {
						fa := mt.Scope.(*area).il
						a.area.il.Link(&key{id: par.Adr()}, items.ID{In: fa, This: &key{id: mt.Id}})
					} else { //поля, элементы массива, некоторые результаты разыменований
						d := &item{}
						v := rt2.ValueOf(in.Frame)[eval.KeyOf(in, link)]
						d.Data(v)
						a.area.il.Put(&key{id: par.Adr()}, d)
					}
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
		ret := &area{all: &halloc{}, il: items.New(), unsafe: true}
		ret.il.Begin()
		return ret
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
