package object

import (
	"ypk/assert"
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
)

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
}

type Ref interface {
	Object() Object
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
}

type ProcedureObject interface {
	Object
	self() ProcedureObject
}

type TypeObject interface {
	Object
	self() TypeObject
}

func New(mode Mode) Object {
	switch mode {
	case HEAD:
		return new(headObject)
	case VARIABLE:
		return new(variableObject)
	case LOCAL_PROC:
		return new(localProcedureObject)
	case CONSTANT:
		return new(constantObject)
	case PARAMETER:
		return new(parameterObject)
	case EXTERNAL_PROC:
		return new(externalProcedureObject)
	case TYPE_PROC:
		return new(typeProcedureObject)
	case FIELD:
		return new(fieldObject)
	case TYPE:
		return new(typeObject)
	default:
		panic("no such object mode")
	}
}

type objectFields struct {
	name string
	typ  Type
	link Object
	comp ComplexType
	ref  []Ref
}

func (of *objectFields) SetType(typ Type)         { of.typ = typ }
func (of *objectFields) SetName(name string)      { of.name = name }
func (of *objectFields) Name() string             { return of.name }
func (of *objectFields) Type() Type               { return of.typ }
func (of *objectFields) Link() Object             { return of.link }
func (of *objectFields) SetLink(o Object)         { of.link = o }
func (of *objectFields) SetComplex(t ComplexType) { of.comp = t }
func (of *objectFields) Complex() ComplexType     { return of.comp }

func (of *objectFields) SetRef(n Ref) {
	assert.For(n != nil, 20)
	exists := func() bool {
		for _, v := range of.ref {
			if v == n {
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
}

func (v *fieldObject) self() FieldObject { return v }

type typeObject struct {
	objectFields
}

func (v *typeObject) self() TypeObject { return v }
