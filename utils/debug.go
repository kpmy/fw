package utils

import "fmt"

var debug = true

func Println(x ...interface{}) {
	if debug {
		fmt.Println(x[0], x[1], x[2])
	}
}
