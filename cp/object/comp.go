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
		fmt.Println("basic comp", a.Qualident(), ",", b.Qualident())
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *arr) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *arr:
		fmt.Println("arr comp", ":", a.Qualident(), ",", b.Qualident())
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *dyn) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *dyn:
		fmt.Println("dyn comp", a.Qualident(), ",", b.Qualident())
		return a == b
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *rec) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *rec:
		fmt.Println("rec comp", a.Name(), ":", a.Qualident(), ",", b.Name(), ":", b.Qualident())
		return a == b
	case *ptr:
		return false
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}

func (a *ptr) Equals(x ComplexType) bool {
	switch b := x.(type) {
	case *ptr:
		fmt.Println("pointer comp", a.Name(), ":", a.Qualident(), ",", b.Name(), ":", b.Qualident())
		return a.Qualident() == b.Qualident()
	case *rec:
		return false
	default:
		halt.As(100, reflect.TypeOf(b))
	}
	panic(0)
}
