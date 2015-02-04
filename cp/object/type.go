package object

import (
	"fmt"
	"fw/cp"
	"ypk/assert"
)

type Type int

const (
	WRONG Type = iota
	NOTYPE
	INTEGER
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
	COMPLEX
	STRING
	SHORTSTRING
	NIL
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
	case STRING:
		return "STRING"
	case SHORTSTRING:
		return "SHORTSTRING"
	case NIL:
		return "NIL"
	default:
		return fmt.Sprint("looks like new type here", int(t))
	}
}

type ComplexType interface {
	Link() Object
	SetLink(o Object)
	Equals(ComplexType) bool
	Qualident(...string) string
	cp.Id
}

type comp struct {
	link Object
	adr  cp.ID
	id   string
}

func (c *comp) Link() Object     { return c.link }
func (c *comp) SetLink(o Object) { c.link = o }

func (c *comp) Adr(a ...cp.ID) cp.ID {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		c.adr = a[0]
	}
	return c.adr
}

func (c *comp) Qualident(a ...string) string {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		c.id = a[0]
	}
	return c.id
}

type BasicType interface {
	ComplexType
	Type() Type
	Base(...Type) Type
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
	Complex(...ComplexType) ComplexType
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
	x.Adr(cp.Next(id))
	return x
}

type basic struct {
	comp
	typ, base Type
}

func (b *basic) Type() Type { return b.typ }
func (b *basic) Base(x ...Type) Type {
	if len(x) > 0 {
		b.base = x[0]
	}
	return b.base
}

func NewDynArrayType(b Type, id int) (ret DynArrayType) {
	ret = &dyn{base: b}
	ret.Adr(cp.Next(id))
	return ret
}

type dyn struct {
	comp
	base Type
}

func (d *dyn) Base() Type { return d.base }

func NewArrayType(b Type, len int64, id int) (ret ArrayType) {
	ret = &arr{base: b, length: len}
	ret.Adr(cp.Next(id))
	return ret
}

type arr struct {
	comp
	base   Type
	length int64
	cmp    ComplexType
}

func (a *arr) Base() Type { return a.base }
func (a *arr) Len() int64 { return a.length }
func (a *arr) Complex(t ...ComplexType) ComplexType {
	if len(t) == 1 {
		a.cmp = t[0]
	} else if len(t) > 1 {
		panic("too many args")
	}
	return a.cmp
}

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
		ret.Adr(cp.Next(id))
	} else {
		ret = &rec{name: n, base: par[0]}
		ret.Adr(cp.Next(id))
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

func NewPointerType(n string, id int) PointerType {
	//fmt.Println("new ptr type", n)
	p := &ptr{name: n}
	p.Adr(cp.Next(id))
	return p
}

func (p *ptr) Name() string { return p.name }
func (p *ptr) Base(x ...ComplexType) ComplexType {
	if len(x) == 1 {
		p.basetyp = x[0]
		//fmt.Println("pbasetyp", p.basetyp, reflect.TypeOf(p.basetyp))
	} else if len(x) > 1 {
		panic("there can be only one")
	}
	return p.basetyp
}
