package cp

import (
	"fmt"
)

type ID int

type Id interface {
	Adr(...ID) ID
}

func (i ID) String() string {
	if i > 0 {
		return fmt.Sprint(int(i), "_", getId(i))
	} else {
		return fmt.Sprint(int(i))
	}
}

func Init() {
	var (
		fake int        = 0
		this int        = 0
		list map[ID]int = make(map[ID]int)
	)
	Next = func(id int) ID {
		if id >= 0 {
			this++
			list[ID(this)] = id
			return ID(this)
		} else {
			return ID(id)
		}
	}
	Some = func() int {
		fake--
		return fake
	}
	getId = func(id ID) int {
		return list[id]
	}
}

var Next func(id int) ID
var Some func() int
var getId func(id ID) int
