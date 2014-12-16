package assert

import (
	"fmt"
)

func For(cond bool, code int) {
	e := fmt.Sprint(code)
	if !cond {
		switch {
		case (code >= 20) && (code < 40):
			e = fmt.Sprintln(code, "precondition violated")
		case (code >= 40) && (code < 59):
			e = fmt.Sprintln(code, "subcondition violated")
		case (code >= 60) && (code < 80):
			e = fmt.Sprintln(code, "postcondition violated")
		default:
		}
		panic(e)
	}
}
