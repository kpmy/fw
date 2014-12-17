package module

import (
	mod "cp/module"
	"os"
	"rt2/context"
	"xev"
	"ypk/assert"
)

type List interface {
	context.ContextAware
	AsList() []*mod.Module
	Load(name string) (*mod.Module, error)
	Loaded(name string) *mod.Module
}

func New() List {
	return new(list).init()
}

type list struct {
	inner map[string]*mod.Module
	d     context.Domain
}

func (l *list) init() *list {
	l.inner = make(map[string]*mod.Module)
	return l
}

func (l *list) AsList() (ret []*mod.Module) {
	if len(l.inner) > 0 {
		ret = make([]*mod.Module, 0)
	}
	for _, v := range l.inner {
		ret = append(ret, v)
	}
	return ret
}

func (l *list) Domain() context.Domain {
	return l.d
}

func (l *list) Init(d context.Domain) {
	l.d = d
}

func (l *list) Handle(msg interface{}) {}
func (l *list) Load(name string) (*mod.Module, error) {
	assert.For(name != "", 20)
	ret := l.Loaded(name)
	if ret == nil {
		path, _ := os.Getwd()
		ret = xev.Load(path, name+".oxf")
		l.inner[name] = ret
	}
	return ret, nil
}

func (l *list) Loaded(name string) *mod.Module {
	assert.For(name != "", 20)
	return l.inner[name]
}

func DomainModule(d context.Domain) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	name := uni.Id(d)
	ml := uni.Discover(context.MOD).(List)
	return ml.Loaded(name)
}
