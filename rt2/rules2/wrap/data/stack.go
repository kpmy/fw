package data

import (
	"fmt"
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/rules2/wrap/data/items"
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

func (a *salloc) push(_o object.Object) {
	switch o := _o.(type) {
	case object.VariableObject:
		switch t := o.Complex().(type) {
		case nil, object.BasicType:
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
			fmt.Println(o.Adr(), o.Name())
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

func (a *salloc) Initialize(n node.Node, par scope.PARAM) (frame.Sequence, frame.WAIT) {
	a.area.il.End()
	return frame.End()
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
	if role == context.SCOPE {
		return &area{all: &salloc{}, il: items.New()}
	} else if role == context.HEAP {
		return &area{all: nil}
		//return &area{all: &halloc{}}
	} else {
		panic(0)
	}
}

func init() {
	scope.New = nn
	//scope.FindObjByName = fn
}
