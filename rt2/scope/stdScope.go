package scope

import (
	"cp/node"
	"fmt"
	"rt2/context"
	rt_mod "rt2/module"
	"ypk/assert"
)

func This(i interface{}) Manager {
	assert.For(i != nil, 20)
	return i.(Manager)
}

func New() Manager {
	return new(manager).init()
}

type manager struct {
	d context.Domain
}

func (m *manager) init() *manager {
	return m
}

func (m *manager) Allocate(n node.Node) {
	mod := rt_mod.DomainModule(m.Domain())
	fmt.Println("allocate", len(mod.Objects[n]), "obj")
}

func (m *manager) Dispose(n node.Node) {
	fmt.Println("dispose")
}

func (m *manager) Calculate(n node.Node) Area {
	return nil
}

func (m *manager) Init(d context.Domain) {
	m.d = d
}

func (m *manager) Domain() context.Domain {
	return m.d
}
