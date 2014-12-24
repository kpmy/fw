package utils

import "fmt"

var debug = true

func Println(x ...interface{}) {
	if debug {
		fmt.Println(x)
	}
}
