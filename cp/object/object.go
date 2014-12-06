package object

type Mode int
type Type int

const (
	HEAD Mode = iota
	VARIABLE
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
