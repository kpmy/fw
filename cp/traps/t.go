package traps

import (
	"strconv"
)

type TRAP int

const (
	Default TRAP = iota
	NILderef
)

func This(i interface{}) TRAP { return TRAP(i.(int32)) }
func (t TRAP) String() string {
	switch t {
	case Default:
		return "0"
	case NILderef:
		return "NIL dereference"
	default:
		return strconv.Itoa(int(t))
	}
}
