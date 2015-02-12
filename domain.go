package main

import (
	"fw/rt2/context"
	"ypk/assert"
)

type stdDomain struct {
	list map[string]context.ContextAware
	//parent context.Domain
	global context.Domain
	god    bool
}

func (d *stdDomain) New() context.Domain { return &stdDomain{global: d.global} }

func (d *stdDomain) Global() context.Domain { return d.global }
func (d *stdDomain) Attach(name string, x context.ContextAware) {
	assert.For(x != nil, 20)
	assert.For(name != context.UNIVERSE, 21)
	if d.list == nil {
		d.list = make(map[string]context.ContextAware)
	}
	assert.For(d.list[name] == nil, 40)
	x.Init(d)
	d.list[name] = x
}

func (d *stdDomain) Discover(name string, opts ...interface{}) (ret context.ContextAware) {
	assert.For(name != "", 20)
	if name == context.VSCOPE {
		assert.For(len(opts) != 0, 20)
	}
	if d.list != nil {
		ret = d.list[name]
	}
	if ret == nil {
		switch {
		case name == context.UNIVERSE:
			ret = d.global
		case name == context.HEAP && !d.god:
			ret = d.global.Discover(name)
		}
	}
	assert.For(ret != nil, 60) //все плохо
	return ret
}

func (d *stdDomain) Domain() context.Domain {
	return d.global
	//return d.parent
}

func (d *stdDomain) Handle(msg interface{}) {}

func (d *stdDomain) Init(dd context.Domain) {
	glob := dd.(*stdDomain)
	assert.For(glob.god == true, 20) //допустим только один уровень вложенности доменов пока
	//	d.parent = dd
	d.global = dd
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
