package data

import (
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/scope"
	"reflect"
	"ypk/halt"
)

type area struct {
	d   context.Domain
	all scope.Allocator
}

type salloc struct {
	area *area
}

func (a *area) Select(id cp.ID, val scope.ValueOf) {}

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
		return &area{all: &salloc{}}
	} else if role == context.HEAP {
		//return &area{all: &halloc{}}
		panic(0)
	} else {
		panic(0)
	}
}

func init() {
	scope.New = nn
	//scope.FindObjByName = fn
}
