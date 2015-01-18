package scope

import (
	"fw/cp"
)

type Operations interface {
	Sum(Value, Value) Value
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

//средство обновления значенияx
type ValueFor func(in Value) (out Value)

func Simple(v Value) ValueFor {
	return func(Value) Value {
		return v
	}
}

var ValueFrom func(v Value) Value

var Ops Operations
