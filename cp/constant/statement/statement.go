package statement

import (
	"strconv"
)

type Statement int

const (
	ASSIGN Statement = iota
	INC
	DEC
	INCL
	EXCL
)

func (s Statement) String() string {
	switch s {
	case ASSIGN:
		return ":="
	case INC:
		return "INC"
	case DEC:
		return "DEC"
	case INCL:
		return "INCL"
	case EXCL:
		return "EXCL"
	default:
		return strconv.Itoa(int(s))
	}
}
