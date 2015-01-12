package operation

import (
	"ypk/assert"
)

type Operation int

const (
	WRONG Operation = iota
	ALIEN_CONV
	ALIEN_MSK
	PLUS
	MINUS
	EQUAL
	LESSER
	LESS_EQUAL
	LEN
	NOT
	NOT_EQUAL
	IS
	GREATER
	ABS
	ODD
	CAP
	BITS
	MIN
	MAX
	DIV
	MOD
	TIMES
	SLASH
	IN
	OR
	AND
	ASH
	GREAT_EQUAL
)

var this map[string]Operation = make(map[string]Operation)

func init() {
	this[PLUS.String()] = PLUS
	this[MINUS.String()] = MINUS
	this[ALIEN_CONV.String()] = ALIEN_CONV
	this[ALIEN_MSK.String()] = ALIEN_MSK
	this[EQUAL.String()] = EQUAL
	this[LESSER.String()] = LESSER
	this[LESS_EQUAL.String()] = LESS_EQUAL
	this[LEN.String()] = LEN
	this[NOT.String()] = NOT
	this[NOT_EQUAL.String()] = NOT_EQUAL
	this[IS.String()] = IS
	this[GREATER.String()] = GREATER
	this[ABS.String()] = ABS
	this[ODD.String()] = ODD
	this[CAP.String()] = CAP
	this[BITS.String()] = BITS
	this[MIN.String()] = MIN
	this[MAX.String()] = MAX
	this[DIV.String()] = DIV
	this[MOD.String()] = MOD
	this[TIMES.String()] = TIMES
	this[SLASH.String()] = SLASH
	this[IN.String()] = IN
	this[OR.String()] = OR
	this[AND.String()] = AND
	this[ASH.String()] = ASH
	this[GREAT_EQUAL.String()] = GREAT_EQUAL
}

func This(s string) (ret Operation) {
	ret = this[s]
	assert.For(ret != WRONG, 40)
	return ret
}

func (o Operation) String() string {
	switch o {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case ALIEN_CONV:
		return "CONV"
	case EQUAL:
		return "="
	case LESSER:
		return "<"
	case LESS_EQUAL:
		return "<="
	case LEN:
		return "LEN"
	case NOT:
		return "~"
	case NOT_EQUAL:
		return "#"
	case IS:
		return "IS"
	case GREATER:
		return ">"
	case ABS:
		return "ABS"
	case ODD:
		return "ODD"
	case CAP:
		return "CAP"
	case BITS:
		return "BITS"
	case MIN:
		return "MIN"
	case MAX:
		return "MAX"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case TIMES:
		return "*"
	case SLASH:
		return "/"
	case IN:
		return "IN"
	case OR:
		return "OR"
	case AND:
		return "&"
	case ASH:
		return "ASH"
	case GREAT_EQUAL:
		return ">="
	case ALIEN_MSK:
		return "MSK"
	default:
		return "?"
	}
}
