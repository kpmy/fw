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
	GREATER
	ABS
	ODD
	CAP
	BITS
	MIN
	MAX
)

func (o Operation) String() string {
	switch o {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case CONVERT:
		return "CONVERT"
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
	default:
		return "?"
	}
}
