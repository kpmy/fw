package object

type Mode int
type Type int

const (
	HEAD Mode = iota
	VARIABLE
	LOCAL_PROCEDURE
	CONSTANT
	PARAMETER
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
	Link() Object
	SetLink(o Object)
	Name() string
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

type ParameterObject interface {
	Object
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
	case PARAMETER:
		return new(parameterObject)
	default:
		panic("no such object mode")
	}
}

type objectFields struct {
	name string
	typ  Type
	link Object
}

func (of *objectFields) SetType(typ Type)    { of.typ = typ }
func (of *objectFields) SetName(name string) { of.name = name }
func (of *objectFields) Name() string        { return of.name }
func (of *objectFields) Type() Type          { return of.typ }
func (of *objectFields) Link() Object        { return of.link }
func (of *objectFields) SetLink(o Object)    { of.link = o }

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

type parameterObject struct {
	objectFields
}

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
	default:
		panic("looks like new type here")
	}
}
