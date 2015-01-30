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
	Gtr(Value, Value) Value
	Max(Value, Value) Value
	Min(Value, Value) Value
	Div(Value, Value) Value
	Mod(Value, Value) Value
	Mult(Value, Value) Value
	Divide(Value, Value) Value
	In(Value, Value) Value
	Ash(Value, Value) Value
	And(Value, Value) Value
	Or(Value, Value) Value
	Geq(Value, Value) Value

	Not(Value) Value
	Abs(Value) Value
	Odd(Value) Value
	Cap(Value) Value
	Bits(Value) Value //это не BITS из КП, BITS(x) = {x}
	Minus(Value) Value

	Is(Value, object.ComplexType) Value
	Conv(Value, object.Type, ...object.ComplexType) Value
	Len(object.Object, Value, Value) Value
	TypeOf(Value) (object.Type, object.ComplexType)
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

type Pointer interface {
	Variable
	Get() Value
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
