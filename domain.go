package main

import (
	"fw/rt2/context"
	"ypk/assert"
)

type stdDomain struct {
	list   map[string]context.ContextAware
	parent context.Domain
	global context.Domain
}

func (d *stdDomain) ConnectTo(name string, x context.ContextAware) {
	assert.For(x != nil, 20)
	assert.For(name != context.UNIVERSE, 21)
	if d.list == nil {
		d.list = make(map[string]context.ContextAware)
	}
	assert.For(d.list[name] == nil, 40)
	x.Init(d)
	d.list[name] = x
}

func (d *stdDomain) Discover(name string) (ret context.ContextAware) {
	assert.For(name != "", 20)
	if d.list != nil {
		ret = d.list[name]
	}
	switch name {
	case context.UNIVERSE:
		ret = d.global
	case context.HEAP:
		ret = d.global.Discover(name)
	}
	return ret
}

func (d *stdDomain) Domain() context.Domain {
	return d.parent
}

func (d *stdDomain) Handle(msg interface{}) {}

func (d *stdDomain) Init(dd context.Domain) {
	glob := dd.(*stdDomain).global
	assert.For(glob == nil, 20) //допустим только один уровень вложенности доменов пока
	d.parent = dd
	if dd.(*stdDomain).global == nil {
		d.global = dd
	} else {
		d.global = glob
	}
}

func (d *stdDomain) Id(c context.ContextAware) (ret string) {
	for k, v := range d.list {
		if v == c {
			ret = k
			break //стыд-то какой
		}
	}
	return ret
}
