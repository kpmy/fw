package rt

import (
	"fmt"
	"reflect"
)

func intOf(x interface{}) (a int) {
	fmt.Println(reflect.TypeOf(x))
	switch x.(type) {
	case *INTEGER:
		z := *x.(*INTEGER)
		a = int(z)
	case *int:
		z := *x.(*int)
		a = z
	case int:
		a = x.(int)
	default:
		panic("unsupported type")
	}
	return a
}
func Sum(_a interface{}, _b interface{}) interface{} {
	var a int = intOf(_a)
	var b int = intOf(_b)
	return a + b
}
