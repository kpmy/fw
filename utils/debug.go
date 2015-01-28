package utils

import "fmt"

var debugFrame = false
var debugScope = false
var debugTrap = true

func Debug(x bool) {
	debugFrame = x
	debugScope = x
}

func PrintFrame(x ...interface{}) {
	if debugFrame {
		fmt.Println(x...)
	}
}

func PrintScope(x ...interface{}) {
	if debugScope {
		fmt.Println(x...)
	}
}

func PrintTrap(x ...interface{}) {
	if debugTrap {
		fmt.Println(x...)
	}
}
