package utils

import "fmt"

var debugFrame = false
var debugScope = true

func PrintFrame(x ...interface{}) {
	if debugFrame {
		fmt.Println(x[0], x[1], x[2])
	}
}

func PrintScope(x ...interface{}) {
	if debugScope {
		fmt.Println(x[0], x[1])
	}
}
