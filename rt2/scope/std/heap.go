package std

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/scope"
	"reflect"
)

type heap struct {
	d    context.Domain
	data *area
	next int64
}

func nh() scope.Manager {
	return &heap{data: &area{ready: true, root: nil, x: make(map[scope.ID]interface{})}}
}

func (h *heap) Allocate(n node.Node) scope.ValueFor {
	switch v := n.(type) {
	case node.VariableNode:
		switch t := v.Object().Complex().(type) {
		case object.PointerType:

		default:
			panic(fmt.Sprintln("unsupported type", reflect.TypeOf(t)))
		}
	default:
		panic(fmt.Sprintln("unsupported node", reflect.TypeOf(v)))
	}
}

func (h *heap) Dispose(n node.Node) {
}

func (h *heap) Target(...scope.Allocator) scope.Allocator {
	return h
}

func (h *heap) Update(i scope.ID, val scope.ValueFor) {}

func (h *heap) Select(i scope.ID) interface{} { return nil }

func (h *heap) Init(d context.Domain) { h.d = d }

func (h *heap) Domain() context.Domain { return h.d }

func (h *heap) Handle(msg interface{}) {}
