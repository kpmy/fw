package modern

import (
	"fmt"
	//cpm "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"reflect"
	"ypk/halt"
)

type halloc struct {
	area *area
}

func (h *halloc) Allocate(n node.Node, par ...interface{}) scope.ValueFor {
	fmt.Println("HEAP ALLOCATE")
	_ = rtm.ModuleOfNode(h.area.d, n)
	if h.area.data == nil {
		h.area.data = append(h.area.data, newlvl())
	}
	switch v := n.(type) {
	case node.VariableNode:
		switch t := v.Object().Complex().(type) {
		case object.PointerType:
			panic(0)
			//h.area.data[0].alloc(mod, nil, )
		default:
			halt.As(100, reflect.TypeOf(t))
		}
	default:
		halt.As(101, reflect.TypeOf(v))
	}
	return nil
}

func (h *halloc) Dispose(n node.Node) {}

func (a *halloc) Join(m scope.Manager) { a.area = m.(*area) }
