package object

import (
	"fmt"
	"fw/cp"
	"reflect"
	"ypk/assert"
)

type Type int

const (
	INTEGER Type = iota
	SHORTINT
	LONGINT
	BYTE
	BOOLEAN
	SHORTREAL
	REAL
	CHAR
	SHORTCHAR
	SET
	PROCEDURE
	//фиктивные типы
	NOTYPE
	COMPLEX
	STRING
	SHORTSTRING
)

func (t Type) String() string {
	switch t {
	case NOTYPE:
		return "NO TYPE"
	case INTEGER:
		return "INTEGER"
	case SHORTINT:
		return "SHORTINT"
	case LONGINT:
		return "LONGINT"
	case BYTE:
		return "BYTE"
	case BOOLEAN:
		return "BOOLEAN"
	case SHORTREAL:
		return "SHORTREAL"
	case REAL:
		return "REAL"
	case CHAR:
		return "CHAR"
	case SHORTCHAR:
		return "SHORTCHAR"
	case SET:
		return "SET"
	case COMPLEX:
		return "COMPLEX"
	case PROCEDURE:
		return "PROCEDURE"
	default:
		return fmt.Sprint("looks like new type here", int(t))
	}
}

type ComplexType interface {
	Link() Object
	SetLink(o Object)
	cp.Id
}

type comp struct {
	link Object
	adr  int
}

func (c *comp) Link() Object     { return c.link }
func (c *comp) SetLink(o Object) { c.link = o }

func (c *comp) Adr(a ...int) int {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		c.adr = a[0]
	}
	return c.adr
}

type BasicType interface {
	ComplexType
	Type() Type
}

type PointerType interface {
	ComplexType
	Base(...ComplexType) ComplexType
	Name() string
}

type ArrayType interface {
	ComplexType
	Base() Type
	Len() int64
}

type DynArrayType interface {
	ComplexType
	Base() Type
}

type RecordType interface {
	ComplexType
	Base() string
	BaseType() RecordType
	SetBase(ComplexType)
	Name() string
}

func NewBasicType(t Type, id int) BasicType {
	x := &basic{typ: t}
	x.Adr(id)
	return x
}

type basic struct {
	comp
	typ Type
}

func (b *basic) Type() Type { return b.typ }

func NewDynArrayType(b Type, id int) (ret DynArrayType) {
	ret = &dyn{base: b}
	ret.Adr(id)
	return ret
}

type dyn struct {
	comp
	base Type
}

func (d *dyn) Base() Type { return d.base }

func NewArrayType(b Type, len int64, id int) (ret ArrayType) {
	ret = &arr{base: b, length: len}
	ret.Adr(id)
	return ret
}

type arr struct {
	comp
	base   Type
	length int64
}

func (a *arr) Base() Type { return a.base }
func (a *arr) Len() int64 { return a.length }

type rec struct {
	comp
	name, base string
	basetyp    RecordType
}

func (r *rec) Name() string { return r.name }
func (r *rec) Base() string { return r.base }

func NewRecordType(n string, id int, par ...string) (ret RecordType) {
	if len(par) == 0 {
		ret = &rec{}
		ret.Adr(id)
	} else {
		ret = &rec{name: n, base: par[0]}
		ret.Adr(id)
	}
	return ret
}

func (r *rec) BaseType() RecordType { return r.basetyp }
func (r *rec) SetBase(t ComplexType) {
	r.basetyp = t.(RecordType)
}

type ptr struct {
	comp
	basetyp ComplexType
	name    string
}

func NewPointerType(n string) PointerType {
	fmt.Println("new ptr type", n)
	return &ptr{name: n}
}

func (p *ptr) Name() string { return p.name }
func (p *ptr) Base(x ...ComplexType) ComplexType {
	if len(x) == 1 {
		p.basetyp = x[0]
		fmt.Println("pbasetyp", p.basetyp, reflect.TypeOf(p.basetyp))
	} else if len(x) > 1 {
		panic("there can be only one")
	}
	return p.basetyp
}
