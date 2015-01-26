package module

import (
	mod "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/xev"
	"os"
	"ypk/assert"
)

type Loader func(*mod.Module)

type List interface {
	context.ContextAware
	AsList() []*mod.Module
	Load(string, ...Loader) (*mod.Module, error)
	Loaded(string) *mod.Module
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

func (l *list) Load(name string, ldr ...Loader) (ret *mod.Module, err error) {
	assert.For(name != "", 20)
	//fmt.Println("loading", name, "loaded", l.Loaded(name) != nil)
	ret = l.Loaded(name)
	var loader Loader = func(m *mod.Module) {}
	if len(ldr) > 0 {
		loader = ldr[0]
	}
	if ret == nil {
		path, _ := os.Getwd()
		ret = xev.Load(path, name+".oz")
		ret.Name = name
		for _, imp := range ret.Imports {
			//fmt.Println("imports", imp.Name, "loaded", l.Loaded(imp.Name) != nil)
			_, err = l.Load(imp.Name, loader)
		}
		if err == nil {
			l.inner[name] = ret
			loader(ret)
			//fmt.Println("loaded", name)
		}
	}
	return ret, err
}

func (l *list) Loaded(name string) *mod.Module {
	assert.For(name != "", 20)
	return l.inner[name]
}

func DomainModule(d context.Domain) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	name := uni.Id(d)
	assert.For(name != "", 40)
	ml := uni.Discover(context.MOD).(List)
	return ml.Loaded(name)
}

func ModuleOfNode(d context.Domain, x node.Node) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, n := range m.Nodes {
			if n == x {
				return m
			}
		}
	}
	return nil
}

func ModuleOfObject(d context.Domain, x object.Object) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, v := range m.Objects {
			for _, o := range v {
				if o == x {
					return m
				}
			}
		}
	}
	return nil
}

func ModuleOfType(d context.Domain, x object.ComplexType) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, v := range m.Types {
			for _, o := range v {
				if o == x {
					return m
				}
			}
		}
	}
	return nil
}
