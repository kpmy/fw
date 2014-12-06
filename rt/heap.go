package rt

import (
	"cp/object"
	"fmt"
)

type INTEGER int

type Heap interface {
	ThisVariable(obj object.Object) interface{}
}

type stdHeap struct {
	inner map[interface{}]interface{}
}

func NewHeap() Heap {
	return new(stdHeap).Init()
}

func (h *stdHeap) Init() *stdHeap {
	h.inner = make(map[interface{}]interface{}, 0)
	return h
}

func (h *stdHeap) ThisVariable(obj object.Object) (ptr interface{}) {
	fmt.Println(obj)
	ptr = h.inner[obj]
	if ptr == nil {
		switch obj.Type() {
		case object.INTEGER:
			ptr = new(int)
			h.inner[obj] = ptr
		default:
			fmt.Println(obj.Type())
			panic("unknown object type")
		}
	}
	return ptr
}
