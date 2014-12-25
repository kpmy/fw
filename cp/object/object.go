package object

type Mode int

const (
	HEAD Mode = iota
	VARIABLE
	LOCAL_PROC
	EXTERNAL_PROC
	CONSTANT
	PARAMETER
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
	default:
		panic("no such object mode")
	}
}

type objectFields struct {
	name string
	typ  Type
	link Object
	comp ComplexType
}

func (of *objectFields) SetType(typ Type)         { of.typ = typ }
func (of *objectFields) SetName(name string)      { of.name = name }
func (of *objectFields) Name() string             { return of.name }
func (of *objectFields) Type() Type               { return of.typ }
func (of *objectFields) Link() Object             { return of.link }
func (of *objectFields) SetLink(o Object)         { of.link = o }
func (of *objectFields) SetComplex(t ComplexType) { of.comp = t }
func (of *objectFields) Complex() ComplexType     { return of.comp }

type variableObject struct {
	objectFields
}

type headObject struct {
	objectFields
}

type localProcedureObject struct {
	objectFields
}

type externalProcedureObject struct {
	objectFields
}

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
