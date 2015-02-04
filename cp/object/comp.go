package object

import (
	"fmt"
	"reflect"
)

import (
	"ypk/halt"
)

func (a *basic) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *basic:
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *arr) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *arr:
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *dyn) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *dyn:
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *rec) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *rec:
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *ptr) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *ptr:
		fmt.Println(a.Name(), ":", a.Qualident(), ",", b.Name(), ":", b.Qualident())
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}
