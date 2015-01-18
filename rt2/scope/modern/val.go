package modern

import (
	"fmt"
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/scope"
	"math/big"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

type data struct {
	link object.Object
	val  interface{}
}

func (x *data) Id() cp.ID {
	return 0
}

func (d *data) Set(v scope.Value) {
	fmt.Println("set data")
	switch x := v.(type) {
	case *data:
		assert.For(d.link.Type() == x.link.Type(), 20)
		d.val = x.val
	case INTEGER:
		switch d.link.Type() {
		case object.INTEGER:
			d.val = x
		case object.LONGINT:
			d.val = LONGINT(x)
		default:
			halt.As(20, d.link.Type())
		}
	case BOOLEAN:
		assert.For(d.link.Type() == object.BOOLEAN, 20)
		d.val = x
	case SHORTCHAR:
		assert.For(d.link.Type() == object.SHORTCHAR, 20)
		d.val = x
	case CHAR:
		assert.For(d.link.Type() == object.CHAR, 20)
		d.val = x
	case SHORTINT:
		assert.For(d.link.Type() == object.SHORTINT, 20)
		d.val = x
	case LONGINT:
		assert.For(d.link.Type() == object.LONGINT, 20)
		d.val = x
	case BYTE:
		assert.For(d.link.Type() == object.BYTE, 20)
		d.val = x
	case SET:
		assert.For(d.link.Type() == object.SET, 20)
		d.val = x
	case REAL:
		assert.For(d.link.Type() == object.REAL, 20)
		d.val = x
	case SHORTREAL:
		assert.For(d.link.Type() == object.SHORTREAL, 20)
		d.val = x
	default:
		panic(fmt.Sprintln(reflect.TypeOf(x)))
	}
}

func (d *data) String() string {
	return fmt.Sprint(d.link.Name(), "=", d.val)
}

type INTEGER int32
type BOOLEAN bool
type BYTE int8
type SHORTINT int16
type LONGINT int64
type SET struct {
	bits *big.Int
}
type CHAR rune
type REAL float64
type SHORTREAL float32
type SHORTCHAR rune

func (x SHORTCHAR) String() string { return fmt.Sprint(rune(x)) }
func (x SHORTREAL) String() string { return fmt.Sprint(float32(x)) }
func (x REAL) String() string      { return fmt.Sprint(float64(x)) }
func (x CHAR) String() string      { return fmt.Sprint(rune(x)) }
func (x SET) String() string       { return fmt.Sprint(x.bits) }
func (x LONGINT) String() string   { return fmt.Sprint(int64(x)) }
func (x SHORTINT) String() string  { return fmt.Sprint(int16(x)) }
func (x BYTE) String() string      { return fmt.Sprint(int8(x)) }
func (x INTEGER) String() string   { return fmt.Sprint(int32(x)) }
func (x BOOLEAN) String() string   { return fmt.Sprint(bool(x)) }

func NewData(o object.Object) (ret scope.Variable) {
	switch o.Type() {
	case object.INTEGER:
		ret = &data{link: o, val: INTEGER(0)}
	case object.BOOLEAN:
		ret = &data{link: o, val: BOOLEAN(false)}
	case object.BYTE:
		ret = &data{link: o, val: BYTE(0)}
	case object.CHAR:
		ret = &data{link: o, val: CHAR(0)}
	case object.LONGINT:
		ret = &data{link: o, val: LONGINT(0)}
	case object.SHORTINT:
		ret = &data{link: o, val: SHORTINT(0)}
	case object.SET:
		ret = &data{link: o, val: SET{bits: big.NewInt(0)}}
	case object.REAL:
		ret = &data{link: o, val: REAL(0)}
	case object.SHORTREAL:
		ret = &data{link: o, val: SHORTREAL(0)}
	case object.SHORTCHAR:
		ret = &data{link: o, val: SHORTCHAR(0)}
	default:
		panic(fmt.Sprintln("unsupported type", o.Type()))
	}
	return ret
}

func fromg(x interface{}) scope.Value {
	switch x := x.(type) {
	case int32:
		return INTEGER(x)
	case bool:
		return BOOLEAN(x)
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func NewConst(n node.Node) scope.Value {
	switch x := n.(type) {
	case node.ConstantNode:
		switch x.Type() {
		case object.INTEGER:
			return INTEGER(x.Data().(int32))
		case object.REAL:
			return REAL(x.Data().(float64))
		case object.BOOLEAN:
			return BOOLEAN(x.Data().(bool))
		case object.SHORTCHAR:
			return SHORTCHAR(x.Data().(rune))
		case object.LONGINT:
			return LONGINT(x.Data().(int64))
		case object.SHORTINT:
			return SHORTINT(x.Data().(int16))
		case object.SHORTREAL:
			return SHORTREAL(x.Data().(float32))
		case object.BYTE:
			return BYTE(x.Data().(int8))
		case object.SET:
			return SET{bits: x.Data().(*big.Int)}
		case object.CHAR:
			return CHAR(x.Data().(rune))
		default:
			panic(fmt.Sprintln(x.Type()))
		}
	}
	panic(0)
}

func vfrom(v scope.Value) scope.Value {
	switch n := v.(type) {
	case *data:
		switch n.link.Type() {
		case object.INTEGER:
			return n.val.(INTEGER)
		case object.BYTE:
			return n.val.(BYTE)
		default:
			halt.As(100, n.link.Type())
		}
	case INTEGER:
		return n
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	return nil
}

func gfrom(v scope.Value) interface{} {
	switch n := v.(type) {
	case *data:
		return gfrom(n.val.(scope.Value))
	case INTEGER:
		return int32(n)
	case BOOLEAN:
		return bool(n)
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	return nil
}

type ops struct{}

func (o *ops) Sum(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Sum(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Sum(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(int32(x) + int32(y))
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			default:
				panic(fmt.Sprintln(reflect.TypeOf(x)))
			}
		}
	}
	panic(0)
}

func (o *ops) Sub(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Sub(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Sub(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(int32(x) - int32(y))
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			default:
				panic(fmt.Sprintln(reflect.TypeOf(x)))
			}
		}
	}
	panic(0)

}
func (o *ops) Conv(a scope.Value, typ object.Type) scope.Value {
	switch typ {
	case object.INTEGER:
		switch x := a.(type) {
		case *data:
			return o.Conv(vfrom(x), typ)
		case BYTE:
			return INTEGER(x)
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	default:
		halt.As(100, typ)
	}
	panic(100)
}

func (o *ops) Eq(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Eq(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Eq(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x == y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			default:
				panic(fmt.Sprintln(reflect.TypeOf(x)))
			}
		}
	}
	panic(0)
}

func (o *ops) Lss(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Lss(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Lss(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x < y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			default:
				panic(fmt.Sprintln(reflect.TypeOf(x)))
			}
		}
	}
	panic(0)
}

func (o *ops) Leq(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Leq(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Leq(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x <= y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			default:
				panic(fmt.Sprintln(reflect.TypeOf(x)))
			}
		}
	}
	panic(0)
}

func init() {
	scope.ValueFrom = vfrom
	scope.GoTypeFrom = gfrom
	scope.TypeFromGo = fromg
	scope.Ops = &ops{}
}
