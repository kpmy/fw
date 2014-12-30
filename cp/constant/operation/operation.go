package operation

type Operation int

const (
	PLUS Operation = iota
	MINUS
	CONVERT
	EQUAL
	LESSER
	LESS_EQUAL
	LEN
	NOT
	NOT_EQUAL
	IS
)
