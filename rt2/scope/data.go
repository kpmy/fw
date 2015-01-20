package scope

import (
	"fw/cp"
	"fw/cp/object"
)

type Operations interface {
	Sum(Value, Value) Value
	Sub(Value, Value) Value

	Eq(Value, Value) Value
	Neq(Value, Value) Value
	Lss(Value, Value) Value
	Leq(Value, Value) Value

	Not(Value) Value
	Conv(Value, object.Type) Value
	Len(object.Object, Value, Value) Value
}

type Value interface {
	String() string
}

type Constant interface {
	Value
}

type Variable interface {
	Id() cp.ID
	Set(Value)
	Value
}

type Ref interface {
	Value
}

type Record interface {
	Variable
	Get(cp.ID) Value
}

type Array interface {
	Variable
	Get(Value) Value
}

//средство обновления значенияx
type ValueFor func(in Value) (out Value)
type ValueOf func(in Value)

func Simple(v Value) ValueFor {
	return func(Value) Value {
		return v
	}
}

var ValueFrom func(v Value) Value
var GoTypeFrom func(v Value) interface{}
var TypeFromGo func(v interface{}) Value
var Ops Operations
