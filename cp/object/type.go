package object

import (
	"fmt"
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
	default:
		return fmt.Sprint("looks like new type here", int(t))
	}
}

type ComplexType interface {
	Link() Object
	SetLink(o Object)
}

type comp struct {
	link Object
}

func (c *comp) Link() Object     { return c.link }
func (c *comp) SetLink(o Object) { c.link = o }

type BasicType interface {
	ComplexType
	Type() Type
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
}

func NewBasicType(t Type) BasicType {
	x := &basic{typ: t}
	return x
}

type basic struct {
	comp
	typ Type
}

func (b *basic) Type() Type { return b.typ }

func NewDynArrayType(b Type) DynArrayType {
	return &dyn{base: b}
}

type dyn struct {
	comp
	base Type
}

func (d *dyn) Base() Type { return d.base }

func NewArrayType(b Type, len int64) ArrayType {
	return &arr{base: b, length: len}
}

type arr struct {
	comp
	base   Type
	length int64
}

func (a *arr) Base() Type { return a.base }
func (a *arr) Len() int64 { return a.length }
