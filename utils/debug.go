package utils

import "log"

var debugFrame = false
var debugScope = false
var debugTrap = true

func Debug(x bool) {
	debugFrame = x
	debugScope = x
}

func Debug2(x, y bool) {
	debugFrame = x
	debugScope = y
}

func PrintFrame(x ...interface{}) {
	if debugFrame {
		log.Println(x...)
	}
}

func PrintScope(x ...interface{}) {
	if debugScope {
		log.Println(x...)
	}
}

func PrintTrap(x ...interface{}) {
	if debugTrap {
		log.Println(x...)
	}
}

func Do(do func()) {
	if debugFrame {
		do()
	}
}
