package modern

import (
	"fmt"
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/scope"
	"fw/utils"
	"math"
	"math/big"
	"reflect"
	"strings"
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

type ptr struct {
	link object.Object
	scope.Pointer
	val *ptrValue
}

type idx struct {
	arr *arr
	idx int
}

func (r *rec) String() string {
	return r.link.Name()
}

func (r *rec) Id() cp.ID {
	return r.link.Adr()
}

func (r *rec) Set(v scope.Value) {
	panic(0)
}

func (r *rec) Get(id cp.ID) scope.Value {
	k := r.l.k[id]
	if r.l.v[k] == nil { //ref
		return r.l.r[k]
	} else {
		return r.l.v[k]
	}
}

func newRec(o object.Object) *rec {
	_, ok := o.Complex().(object.RecordType)
	assert.For(ok, 20)
	return &rec{link: o}
}

func (p *proc) String() string {
	return fmt.Sprint(p.link.Adr(), p.link.Name())
}

func (x *data) Id() cp.ID { return x.link.Adr() }

func (x *arr) Id() cp.ID { return x.link.Adr() }

func (x *dynarr) Id() cp.ID { return x.link.Adr() }

func (a *arr) Set(v scope.Value) {
	switch x := v.(type) {
	case *arr:
		a.Set(STRING(x.tryString()))
	case STRING:
		v := make([]interface{}, int(a.length))
		for i := 0; i < int(a.length) && i < len(x); i++ {
			v[i] = CHAR(x[i])
		}
		a.val = v
	case SHORTSTRING:
		v := make([]interface{}, int(a.length))
		for i := 0; i < int(a.length) && i < len(x); i++ {
			v[i] = SHORTCHAR(x[i])
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
	case *dynarr:
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
	case INTEGER:
		a.val = make([]interface{}, int(x))
	default:
		halt.As(100, reflect.TypeOf(x))
	}
}

func (a *arr) tryString() (ret string) {
	for i := 0; i < len(a.val) && a.val[i] != nil; i++ {
		switch x := a.val[i].(type) {
		case CHAR:
			if int(x) == 0 {
				break
			}
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		case SHORTCHAR:
			if int(x) == 0 {
				break
			}
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	}
	return ret
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

func (a *arr) Get(id scope.Value) scope.Value {
	switch i := id.(type) {
	case *data:
		return a.Get(i.val.(scope.Value))
	case INTEGER:
		assert.For(int64(i) < a.length, 20)
		if len(a.val) == 0 {
			a.val = make([]interface{}, int(a.length))
		}
		return &idx{arr: a, idx: int(i)}
	default:
		halt.As(100, reflect.TypeOf(i))
	}
	panic(0)
}

func (i *idx) Id() cp.ID {
	return i.arr.Id()
}

func (i *idx) String() string {
	return fmt.Sprint("@", i.Id(), "[", i.idx, "]")
}

func (i *idx) Set(v scope.Value) {
	//fmt.Println(i, len(i.arr.val))
	switch x := v.(type) {
	case *idx:
		i.arr.val[i.idx] = x.arr.val[x.idx]
	case *data:
		i.Set(x.val.(scope.Value))
	case CHAR:
		i.arr.val[i.idx] = x
	default:
		halt.As(100, reflect.TypeOf(x))
	}
}

func (a *dynarr) tryString() (ret string) {
	for i := 0; i < len(a.val) && a.val[i] != nil; i++ {
		switch x := a.val[i].(type) {
		case CHAR:
			if int(x) == 0 {
				break
			}
			ret = fmt.Sprint(ret, string([]rune{rune(x)}))
		case SHORTCHAR:
			if int(x) == 0 {
				break
			}
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
	utils.PrintScope("set data")
	switch x := v.(type) {
	case *data:
		if d.link.Type() == x.link.Type() {
			d.val = x.val
		} else {
			d.Set(x.val.(scope.Value))
		}
	case *proc:
		assert.For(d.link.Type() == object.COMPLEX, 20)
		t, ok := d.link.Complex().(object.BasicType)
		assert.For(ok, 21, reflect.TypeOf(d.link.Complex()))
		assert.For(t.Type() == object.PROCEDURE, 22)
		d.val = x
	case *idx:
		d.val = x.arr.val[x.idx]
	case INTEGER:
		switch d.link.Type() {
		case object.INTEGER:
			d.val = x
		case object.LONGINT:
			d.val = LONGINT(x)
		case object.REAL:
			d.val = REAL(x)
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

func (p *ptr) String() string {
	return fmt.Sprint("pointer ", p.link.Complex().(object.PointerType).Name(), "&", p.val)
}

func (p *ptr) Id() cp.ID { return p.link.Adr() }

func (p *ptr) Set(v scope.Value) {
	switch x := v.(type) {
	case *ptr:
		p.Set(x.val)
	case *ptrValue:
		p.val = x
	case PTR:
		assert.For(x == NIL, 40)
		p.val = nil
	default:
		halt.As(100, reflect.TypeOf(x))
	}
}

func (p *ptr) Get() scope.Value {
	if p.val == nil {
		return NIL
	} else {
		return p.val.scope.Select(p.val.id)
	}
}

func newPtr(o object.Object) scope.Variable {
	_, ok := o.Complex().(object.PointerType)
	assert.For(ok, 20)
	return &ptr{link: o}
}

type ptrValue struct {
	scope *area
	id    cp.ID
	link  object.Object
}

func (p *ptrValue) String() string {
	return fmt.Sprint(p.id)
}

type PTR int

func (p PTR) String() string {
	return "NIL"
}

const NIL PTR = 0

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

func newData(o object.Object) (ret scope.Variable) {
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
	case *big.Int:
		return SET{bits: x}
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func newProc(o object.Object) scope.Value {
	p, ok := o.(object.ProcedureObject)
	assert.For(ok, 20, reflect.TypeOf(o))
	return &proc{link: p}
}

func newConst(n node.Node) scope.Value {
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
		case object.NIL:
			return NIL
		case object.COMPLEX: //не может существовать в реальности, используется для передачи параметров от рантайма
			switch d := x.Data().(type) {
			case *ptrValue:
				return d
			default:
				halt.As(100, reflect.TypeOf(d))
			}
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
		case object.CHAR:
			return n.val.(CHAR)
		case object.SET:
			return n.val.(SET)
		case object.BOOLEAN:
			return n.val.(BOOLEAN)
		case object.REAL:
			return n.val.(REAL)
		case object.LONGINT:
			return n.val.(LONGINT)
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
	case *rec:
		return n
	case *proc:
		return n.link
	case *dynarr:
		switch n.link.Complex().(object.DynArrayType).Base() {
		case object.SHORTCHAR:
			if n.val != nil {
				return n.tryString()
			} else {
				return ""
			}
		case object.CHAR:
			if n.val != nil {
				return n.tryString()
			} else {
				return ""
			}
		default:
			halt.As(100, n.link.Complex().(object.DynArrayType).Base())
		}
		panic(0)
	case *arr:
		switch n.link.Complex().(object.ArrayType).Base() {
		case object.SHORTCHAR:
			if n.val != nil {
				return n.tryString()
			} else {
				return ""
			}
		default:
			halt.As(100, n.link.Complex().(object.ArrayType).Base())
		}
		panic(0)
	case PTR:
		assert.For(n == NIL, 40)
		return nil
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
			case SET:
				switch y := b.(type) {
				case SET:
					return SET{bits: x.bits.Add(x.bits, y.bits)}
				case INTEGER: // INCL(SET, INTEGER)
					return SET{bits: x.bits.SetBit(x.bits, int(y), 1)}
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case *arr:
				switch y := b.(type) {
				case *arr:
					switch {
					case x.link.Type() == y.link.Type() && x.link.Complex().(object.ArrayType).Base() == object.CHAR:
						return STRING(x.tryString() + y.tryString())
					default:
						halt.As(100, x.link.Type(), y.link.Type())
					}
				case STRING:
					switch {
					case x.link.Complex().(object.ArrayType).Base() == object.CHAR:
						return STRING(x.tryString() + string(y))
					default:
						halt.As(100, x.link.Type())
					}
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case *dynarr:
				switch y := b.(type) {
				case STRING:
					switch {
					case x.link.Complex().(object.DynArrayType).Base() == object.CHAR:
						return STRING(x.tryString() + string(y))
					default:
						halt.As(100, x.link.Type())
					}
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case STRING:
				switch y := b.(type) {
				case STRING:
					return STRING(string(x) + string(y))
				default:
					halt.As(100, reflect.TypeOf(y))
				}
			case REAL:
				switch y := b.(type) {
				case REAL:
					return REAL(x + y)
				default:
					halt.As(100, reflect.TypeOf(y))
				}
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
					return LONGINT(x + y)
				default:
					halt.As(100, reflect.TypeOf(y))
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
			case SET:
				switch y := b.(type) {
				case SET:
					return SET{bits: x.bits.Sub(x.bits, y.bits)}
				case INTEGER:
					return SET{bits: x.bits.SetBit(x.bits, int(y), 0)}
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

func (o *ops) In(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.In(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.In(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case SET:
					fmt.Println("IN врет")
					return BOOLEAN(false)
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

func (o *ops) Min(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Min(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Min(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(int32(math.Min(float64(x), float64(y))))
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

func (o *ops) Max(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Max(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Max(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(int32(math.Max(float64(x), float64(y))))
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

func (o *ops) And(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.And(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.And(a, vfrom(b))
		default:
			switch x := a.(type) {
			case BOOLEAN:
				switch y := b.(type) {
				case BOOLEAN:
					return BOOLEAN(x && y)
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

func (o *ops) Or(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Or(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Or(a, vfrom(b))
		default:
			switch x := a.(type) {
			case BOOLEAN:
				switch y := b.(type) {
				case BOOLEAN:
					return BOOLEAN(x || y)
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

func (o *ops) Ash(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Max(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Max(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(x << uint(y))
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

func (o *ops) Div(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Div(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Div(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(x / y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
					return LONGINT(x / y)
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

func (o *ops) Mod(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Mod(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Mod(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(x % y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
					return LONGINT(x % y)
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

func (o *ops) Mult(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Mult(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Mult(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return INTEGER(x * y)
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

func (o *ops) Divide(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Divide(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Divide(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return REAL(float64(x) / float64(y))
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case REAL:
				switch y := b.(type) {
				case REAL:
					return REAL(float64(x) / float64(y))
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
			case *arr:
				ret = INTEGER(t.length)
			case *dynarr:
				ret = INTEGER(len(t.val))
			default:
				halt.As(100, "unsupported", reflect.TypeOf(t))
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

func (o *ops) Is(a scope.Value, typ object.ComplexType) scope.Value {
	var compare func(x, a object.ComplexType) bool
	compare = func(_x, _a object.ComplexType) bool {
		switch x := _x.(type) {
		case object.RecordType:
			switch a := _a.(type) {
			case object.RecordType:
				switch {
				case x.Name() == a.Name():
					//	fmt.Println("eq")
					return true //опасно сравнивать имена конеш
				case x.BaseType() != nil:
					//	fmt.Println("go base")
					return compare(x.BaseType(), a)
				default:
					return false
				}
			default:
				halt.As(100, reflect.TypeOf(a))
			}
		case object.PointerType:
			switch a := _a.(type) {
			case object.PointerType:
				switch {
				case x.Name() == a.Name():
					//	fmt.Println("eq")
					return true //опасно сравнивать имена конеш
				case x.Base() != nil:
					//	fmt.Println("go base")
					return compare(x.Base(), a)
				default:
					return false
				}
			default:
				halt.As(100, reflect.TypeOf(a))
			}
		default:
			halt.As(100, reflect.TypeOf(a))
		}
		panic(0)
	}
	switch x := a.(type) {
	case *rec:
		z, a := x.link.Complex().(object.RecordType)
		y, b := typ.(object.RecordType)
		//fmt.Println("compare", x.link.Complex(), typ, a, b, a && b && compare(z, y))
		return BOOLEAN(a && b && compare(z, y))
	case *ptr:
		z, a := x.link.Complex().(object.PointerType)
		y, b := typ.(object.PointerType)
		//fmt.Println("compare", x.link.Complex(), typ, a, b, a && b && compare(z, y))
		return BOOLEAN(a && b && compare(z, y))
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(0)
}

func (o *ops) Conv(a scope.Value, typ object.Type, comp ...object.ComplexType) scope.Value {
	switch typ {
	case object.INTEGER:
		switch x := a.(type) {
		case *data:
			return o.Conv(vfrom(x), typ)
		case BYTE:
			return INTEGER(x)
		case SET:
			return INTEGER(x.bits.Int64())
		case REAL:
			return INTEGER(x)
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	case object.SET:
		switch x := a.(type) {
		case *data:
			return o.Conv(vfrom(x), typ)
		case INTEGER:
			return SET{bits: big.NewInt(int64(x))}
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	case object.REAL:
		switch x := a.(type) {
		case *data:
			return o.Conv(vfrom(x), typ)
		case INTEGER:
			return REAL(float64(x))
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	case object.CHAR:
		switch x := a.(type) {
		case *data:
			return o.Conv(vfrom(x), typ)
		case LONGINT:
			return CHAR(rune(x))
		default:
			halt.As(100, reflect.TypeOf(x))
		}
	case object.NOTYPE:
		assert.For(len(comp) > 0, 20)
		switch t := comp[0].(type) {
		case object.BasicType:
			switch t.Type() {
			case object.SHORTSTRING:
				switch x := a.(type) {
				case *dynarr:
					return SHORTSTRING(x.tryString())
				case *arr:
					return SHORTSTRING(x.tryString())
				default:
					halt.As(100, reflect.TypeOf(x))
				}
			default:
				halt.As(100, t.Type())
			}
		default:
			halt.As(100, reflect.TypeOf(t))
		}
	default:
		halt.As(100, typ)
	}
	panic(100)
}

func (o *ops) Not(a scope.Value) scope.Value {
	switch x := a.(type) {
	case *data:
		return o.Not(vfrom(x))
	case BOOLEAN:
		return BOOLEAN(!x)
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func (o *ops) Abs(a scope.Value) scope.Value {
	switch x := a.(type) {
	case *data:
		return o.Abs(vfrom(x))
	case INTEGER:
		return INTEGER(int32(math.Abs(float64(x))))
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func (o *ops) Odd(a scope.Value) scope.Value {
	switch x := a.(type) {
	case *data:
		return o.Odd(vfrom(x))
	case INTEGER:
		return BOOLEAN(int32(math.Abs(float64(x)))%2 == 1)
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func (o *ops) Cap(a scope.Value) scope.Value {
	switch x := a.(type) {
	case *data:
		return o.Cap(vfrom(x))
	case CHAR:
		return CHAR([]rune(strings.ToUpper(string(x)))[0])
	default:
		halt.As(100, reflect.TypeOf(x))
	}
	panic(100)
}

func (o *ops) Bits(a scope.Value) scope.Value {
	switch x := a.(type) {
	case *data:
		return o.Bits(vfrom(x))
	case INTEGER:
		return SET{bits: big.NewInt(int64(x))}
	default:
		halt.As(100, reflect.TypeOf(x))
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
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
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

func (o *ops) Neq(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Neq(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Neq(a, vfrom(b))
		default:
			switch x := a.(type) {
			case *ptr:
				switch y := b.(type) {
				case PTR:
					assert.For(y == NIL, 40)
					return BOOLEAN(x.val != nil && x.val.id != 0)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x != y)
				default:
					panic(fmt.Sprintln(reflect.TypeOf(y)))
				}
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
					return BOOLEAN(x != y)
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
			case LONGINT:
				switch y := b.(type) {
				case LONGINT:
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

func (o *ops) Gtr(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Gtr(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Gtr(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x > y)
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

func (o *ops) Geq(a, b scope.Value) scope.Value {
	switch a.(type) {
	case *data:
		return o.Geq(vfrom(a), b)
	default:
		switch b.(type) {
		case *data:
			return o.Geq(a, vfrom(b))
		default:
			switch x := a.(type) {
			case INTEGER:
				switch y := b.(type) {
				case INTEGER:
					return BOOLEAN(x >= y)
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

func (o *ops) TypeOf(x scope.Value) (object.Type, object.ComplexType) {
	switch v := x.(type) {
	case *ptr:
		if v.val != nil {
			return v.val.link.Type(), v.val.link.Complex()
		}
	case *rec:
		return v.link.Type(), v.link.Complex()
	default:
		halt.As(100, reflect.TypeOf(v))
	}
	return object.NOTYPE, nil
}

func init() {
	scope.ValueFrom = vfrom
	scope.GoTypeFrom = gfrom
	scope.TypeFromGo = fromg
	scope.Ops = &ops{}
}
