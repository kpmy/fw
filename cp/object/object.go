package object

type Mode int
type Type int

const (
	HEAD Mode = iota
	VARIABLE
	LOCAL_PROCEDURE
)

const (
	NOTYPE Type = iota
	INTEGER
)

type Object interface {
	SetName(name string)
	SetType(typ Type)
	Type() Type
}

func New(mode Mode) Object {
	switch mode {
	case HEAD:
		return new(headObject)
	case VARIABLE:
		return new(variableObject)
	case LOCAL_PROCEDURE:
		return new(localProcedureObject)
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
