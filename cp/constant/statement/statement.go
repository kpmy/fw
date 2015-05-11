package statement

import (
	"github.com/kpmy/ypk/assert"
	"strconv"
)

type Statement int

const (
	WRONG Statement = iota
	ASSIGN
	INC
	DEC
	INCL
	EXCL
	NEW
)

var this map[string]Statement

func init() {
	this = make(map[string]Statement)
	this[ASSIGN.String()] = ASSIGN
	this[INC.String()] = INC
	this[DEC.String()] = DEC
	this[INCL.String()] = INCL
	this[EXCL.String()] = EXCL
	this[NEW.String()] = NEW
}

func This(s string) (ret Statement) {
	ret = this[s]
	assert.For(ret != WRONG, 60)
	return ret
}

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
	case NEW:
		return "NEW"
	default:
		return strconv.Itoa(int(s))
	}
}
