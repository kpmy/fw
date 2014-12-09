package constant

type Class int

const (
	ENTER Class = iota
	ASSIGN
	VARIABLE
	DYADIC
	CONSTANT
	CALL
	PROCEDURE
)
