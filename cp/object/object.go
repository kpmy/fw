package object

type Mode int
type Type int

const (
	HEAD Mode = iota
	VARIABLE
	LOCAL_PROCEDURE
	CONSTANT
)

const (
	NOTYPE Type = iota
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
)

type Object interface {
	SetName(name string)
	SetType(typ Type)
	Type() Type
}

type VariableObject interface {
	Object
	This() VariableObject
}

type ConstantObject interface {
	Object
	SetData(x interface{})
	Data() interface{}
}

func New(mode Mode) Object {
	switch mode {
	case HEAD:
		return new(headObject)
	case VARIABLE:
		return new(variableObject)
	case LOCAL_PROCEDURE:
		return new(localProcedureObject)
	case CONSTANT:
		return new(constantObject)
	default:
		panic("no such object mode")
	}
}

type objectFields struct {
	name string
	typ  Type
}

func (of *objectFields) SetType(typ Type) {
	of.typ = typ
}

func (of *objectFields) SetName(name string) {
	of.name = name
}

func (of *objectFields) Type() Type {
	return of.typ
}

type variableObject struct {
	objectFields
}

type headObject struct {
	objectFields
}

type localProcedureObject struct {
	objectFields
}

func (v *variableObject) This() VariableObject { return v }

type constantObject struct {
	objectFields
	val interface{}
}

func (o *constantObject) SetData(x interface{}) {
	o.val = x
}

func (o *constantObject) Data() interface{} { return o.val }
