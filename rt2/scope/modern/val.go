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

type arr struct {
	link   object.Object
	val    []interface{}
	length int64
}

type dynarr struct {
	link object.Object
	val  []interface{}
}

type proc struct {
	link object.Object
}

type rec struct {
	link object.Object
	scope.Record
	l *level
}

func (r *rec) String() string {
	return r.link.Name()
}

func (r *rec) Id() cp.ID {
	return r.link.Adr()
}

func (r *rec) Set(v scope.Value) {

}

func newRec(o object.Object) *rec {
	_, ok := o.Complex().(object.RecordType)
	assert.For(ok, 20)
	return &rec{link: o}
}

func (p *proc) String() string {
	return fmt.Sprint(p.link.Adr(), p.link.Name())
}

func (x *data) Id() cp.ID {
	return x.link.Adr()
}

func (x *arr) Id() cp.ID {
	return x.link.Adr()
}

func (x *dynarr) Id() cp.ID {
	return x.link.Adr()
}

func (a *arr) Set(v scope.Value) {
	switch x := v.(type) {
	case STRING:
		v := make([]interface{}, int(a.length))
		for i := 0; i < int(a.length) && i < len(x); i++ {
			v[i] = CHAR(x[i])
		}
		a.val = v
	default:
		halt.As(100, reflect.TypeOf(x))
	}
}

func (a *dynarr) Set(v scope.Value) {
	switch x := v.(type) {
	case *arr:
		a.val = x.val
	case STRING:
		v := make([]interface{}, len(x))
		for i := 0; i < len(x); i++ {
			v[i] = CHAR(x[i])
		}
		a.val = v
	case SHORTSTRING:
		v := make([]interface{}, len(x))
		for i := 0; i < len(x); i++ {
			v[i] = SHORTCHAR(x[i])
		}
		a.val = v
	default:
		halt.As(100, reflect.TypeOf(x))
	}
}

func (a *arr) String() (ret string) {
	ret = fmt.Sprint("array", "[", a.length, "]")
	for i := 0; i < len(a.val) && a.val[i] != nil; i++ {
		switch x := a.val[i].(type) {
		case CHAR:
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		case SHORTCHAR:
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	}
	return ret
}

func (a *dynarr) String() (ret string) {
	ret = fmt.Sprint("dyn array")
	for i := 0; i < len(a.val) && a.val[i] != nil; i++ {
		switch x := a.val[i].(type) {
		case CHAR:
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		case SHORTCHAR:
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	}
	return ret
}

func (d *data) Set(v scope.Value) {
	fmt.Println("set data")
	switch x := v.(type) {
	case *data:
		assert.For(d.link.Type() == x.link.Type(), 20)
		d.val = x.val
	case *proc:
		assert.For(d.link.Type() == object.COMPLEX, 20)
		t, ok := d.link.Complex().(object.BasicType)
		assert.For(ok, 21, reflect.TypeOf(d.link.Complex()))
		assert.For(t.Type() == object.PROCEDURE, 22)
		d.val = x
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
type STRING string
type SHORTSTRING string

func (x SHORTSTRING) String() string { return string(x) }
func (x STRING) String() string      { return string(x) }
func (x SHORTCHAR) String() string   { return fmt.Sprint(rune(x)) }
func (x SHORTREAL) String() string   { return fmt.Sprint(float32(x)) }
func (x REAL) String() string        { return fmt.Sprint(float64(x)) }
func (x CHAR) String() string        { return fmt.Sprint(rune(x)) }
func (x SET) String() string         { return fmt.Sprint(x.bits) }
func (x LONGINT) String() string     { return fmt.Sprint(int64(x)) }
func (x SHORTINT) String() string    { return fmt.Sprint(int16(x)) }
func (x BYTE) String() string        { return fmt.Sprint(int8(x)) }
func (x INTEGER) String() string     { return fmt.Sprint(int32(x)) }
func (x BOOLEAN) String() string     { return fmt.Sprint(bool(x)) }

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
	case object.COMPLEX:
		switch t := o.Complex().(type) {
		case object.BasicType:
			switch t.Type() {
			case object.PROCEDURE:
				ret = &data{link: o, val: nil}
			default:
				halt.As(100, t.Type())
			}
		case object.ArrayType:
			ret = &arr{link: o, length: t.Len()}
		case object.DynArrayType:
			ret = &dynarr{link: o}
		default:
			halt.As(100, reflect.TypeOf(t))

		}
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

func NewProc(o object.Object) scope.Value {
	p, ok := o.(object.ProcedureObject)
	assert.For(ok, 20, reflect.TypeOf(o))
	return &proc{link: p}
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
		case object.STRING:
			return STRING(x.Data().(string))
		case object.SHORTSTRING:
			return SHORTSTRING(x.Data().(string))
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
		if n.val == nil {
			return nil
		} else {
			return gfrom(n.val.(scope.Value))
		}
	case *proc:
		return n.link
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

func (o *ops) Len(a object.Object, _a, _b scope.Value) (ret scope.Value) {
	//assert.For(a != nil, 20)
	assert.For(_b != nil, 21)
	var b int32 = gfrom(_b).(int32)
	assert.For(b == 0, 22)
	if a != nil {
		assert.For(a.Type() == object.COMPLEX, 23)
		switch typ := a.Complex().(type) {
		case object.ArrayType:
			ret = INTEGER(int32(typ.Len()))
		case object.DynArrayType:
			switch t := _a.(type) {
			//case string:
			//	ret = int64(utf8.RuneCountInString(_a.(string)))
			default:
				ret = INTEGER(0)
				fmt.Sprintln("unsupported", reflect.TypeOf(t))
			}
		default:
			panic(fmt.Sprintln("unsupported", reflect.TypeOf(a.Complex())))
		}
	} else {
		switch a := _a.(type) {
		//		case string:
		//			ret = int64(utf8.RuneCountInString(_a.(string)))
		//		case []interface{}:
		//			ret = int64(len(_a.([]interface{})))
		case *arr:
			ret = INTEGER(int32(a.length))
		case *dynarr:
			ret = INTEGER(int32(len(a.val)))
		default:
			panic(fmt.Sprintln("unsupported", reflect.TypeOf(a)))
		}
	}
	return ret
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
