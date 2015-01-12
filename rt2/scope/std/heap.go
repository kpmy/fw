package std

import (
	"fw/rt2/context"
	"fw/rt2/scope"
)

type heap struct {
	d context.Domain
}

func nh() scope.Manager {
	return &heap{}
}

func (h *heap) Target(...scope.Allocator) scope.Allocator {
	return nil
}

func (h *heap) Update(i scope.ID, val scope.ValueFor) {}

func (h *heap) Select(i scope.ID) interface{} { return nil }

func (h *heap) Init(d context.Domain) { h.d = d }

func (h *heap) Domain() context.Domain { return h.d }

func (h *heap) Handle(msg interface{}) {}
