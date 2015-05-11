package object

import (
	"fw/cp"
	"github.com/kpmy/ypk/assert"
	"strconv"
)

type Mode int

const (
	HEAD Mode = iota
	VARIABLE
	LOCAL_PROC
	EXTERNAL_PROC
	TYPE_PROC
	CONSTANT
	PARAMETER
	FIELD
	TYPE
	MODULE
)

func (m Mode) String() string {
	switch m {
	case HEAD:
		return "HEAD"
	case VARIABLE:
		return "VARIABLE"
	case LOCAL_PROC:
		return "LOCAL PROCEDURE"
	case EXTERNAL_PROC:
		return "EXTERNAL PROCEDURE"
	case TYPE_PROC:
		return "METHOD"
	case CONSTANT:
		return "CONSTANT"
	case PARAMETER:
		return "PARAMETER"
	case FIELD:
		return "FIELD"
	case TYPE:
		return "TYPE"
	case MODULE:
		return "MODULE"
	default:
		return strconv.Itoa(int(m))
	}
}

type Object interface {
	SetName(name string)
	SetType(typ Type)
	Type() Type
	SetComplex(typ ComplexType)
	Complex() ComplexType
	Link() Object
	SetLink(o Object)
	Name() string
	SetRef(n Ref)
	Ref() []Ref
	Imp(...string) string
	cp.Id
	Mode(...Mode) Mode
}

type Ref interface {
	cp.Id
}

type VariableObject interface {
	Object
	self() VariableObject
}

type ConstantObject interface {
	Object
	SetData(x interface{})
	Data() interface{}
}

type ParameterObject interface {
	Object
	self() ParameterObject
}

type FieldObject interface {
	Object
	self() FieldObject
	TypeOf(...ComplexType) ComplexType
}

type ProcedureObject interface {
	Object
	self() ProcedureObject
}

type TypeObject interface {
	Object
	self() TypeObject
}

type Module interface {
	self() Module
}

func New(mode Mode, id int) (ret Object) {
	switch mode {
	case HEAD:
		ret = new(headObject)
	case VARIABLE:
		ret = new(variableObject)
	case LOCAL_PROC:
		ret = new(localProcedureObject)
	case CONSTANT:
		ret = new(constantObject)
	case PARAMETER:
		ret = new(parameterObject)
	case EXTERNAL_PROC:
		ret = new(externalProcedureObject)
	case TYPE_PROC:
		ret = new(typeProcedureObject)
	case FIELD:
		ret = new(fieldObject)
	case TYPE:
		ret = new(typeObject)
	case MODULE:
		ret = new(mod)
	default:
		panic("no such object mode")
	}
	ret.Mode(mode)
	ret.Adr(cp.Next(id))
	return ret
}

type objectFields struct {
	name string
	typ  Type
	link Object
	comp ComplexType
	ref  []Ref
	adr  cp.ID
	mod  Mode
	imp  string
}

func (of *objectFields) SetType(typ Type)         { of.typ = typ }
func (of *objectFields) SetName(name string)      { of.name = name }
func (of *objectFields) Name() string             { return of.name }
func (of *objectFields) Type() Type               { return of.typ }
func (of *objectFields) Link() Object             { return of.link }
func (of *objectFields) SetLink(o Object)         { of.link = o }
func (of *objectFields) SetComplex(t ComplexType) { of.comp = t }
func (of *objectFields) Complex() ComplexType     { return of.comp }

func (of *objectFields) Adr(a ...cp.ID) cp.ID {
	assert.For(len(a) <= 2, 20)
	if len(a) == 1 {
		of.adr = a[0]
	} else if len(a) == 0 && of.imp != "" {
		panic(123)
	}
	return of.adr
}

func (of *objectFields) Imp(a ...string) string {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		of.imp = a[0]
	}
	return of.imp
}

func (of *objectFields) Mode(a ...Mode) Mode {
	assert.For(len(a) <= 1, 20)
	if len(a) == 1 {
		of.mod = a[0]
	}
	return of.mod
}

func (of *objectFields) SetRef(n Ref) {
	assert.For(n != nil, 20)
	exists := func() bool {
		for _, v := range of.ref {
			if v.Adr() == n.Adr() {
				return true
			}
		}
		return false
	}
	if !exists() {
		of.ref = append(of.ref, n)
	}
}

func (of *objectFields) Ref() []Ref { return of.ref }

type variableObject struct {
	objectFields
}

type headObject struct {
	objectFields
}

type localProcedureObject struct {
	objectFields
}

func (p *localProcedureObject) self() ProcedureObject { return p }

type externalProcedureObject struct {
	objectFields
}

func (p *externalProcedureObject) self() ProcedureObject { return p }

type typeProcedureObject struct {
	objectFields
}

func (p *typeProcedureObject) self() ProcedureObject { return p }

func (v *variableObject) self() VariableObject { return v }

type constantObject struct {
	objectFields
	val interface{}
}

func (o *constantObject) SetData(x interface{}) {
	o.val = x
}

func (o *constantObject) Data() interface{} { return o.val }

type parameterObject struct {
	objectFields
}

func (v *parameterObject) self() ParameterObject { return v }

type fieldObject struct {
	objectFields
	comp ComplexType
}

func (v *fieldObject) self() FieldObject { return v }

func (v *fieldObject) TypeOf(t ...ComplexType) ComplexType {
	if len(t) == 1 {
		v.comp = t[0]
	}
	return v.comp
}

type typeObject struct {
	objectFields
}

func (v *typeObject) self() TypeObject { return v }

type mod struct {
	objectFields
}

func (v *mod) self() Module { return v }
