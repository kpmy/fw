package rt

import (
	"cp/object"
	"fmt"
	"reflect"
)

type INTEGER int

type Variable interface {
	Set(interface{})
}

func (i *INTEGER) Set(val interface{}) {
	if val == nil {
		panic("cannot be nil")
	}
	switch val.(type) {
	case int:
		*i = INTEGER(val.(int))
	case INTEGER:
		*i = val.(INTEGER)
	case *INTEGER:
		*i = *val.(*INTEGER)
	default:
		fmt.Print(reflect.TypeOf(val), " ")
		panic("wrong type for INTEGER")
	}
	fmt.Println("set", int(*i))
}

type Heap interface {
	This(obj object.Object) Variable
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

func (h *stdHeap) This(obj object.Object) (ptr Variable) {
	p := h.inner[obj]
	if p == nil {
		switch obj.Type() {
		case object.INTEGER:
			ptr = new(INTEGER)
			h.inner[obj] = ptr
		default:
			fmt.Println(obj.Type())
			panic("unknown object type")
		}
	} else {
		ptr = p.(Variable)
	}
	return ptr
}
