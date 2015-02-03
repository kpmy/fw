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

type Digest interface{}

type dig struct {
	fake int
	this int
	list map[ID]int
}

func Init() Digest {
	var (
		d *dig = &dig{list: make(map[ID]int)}
	)
	Next = func(id int) ID {
		if id >= 0 {
			d.this++
			d.list[ID(d.this)] = id
			return ID(d.this)
		} else {
			return ID(id)
		}
	}
	Some = func() int {
		d.fake--
		return d.fake
	}
	getId = func(id ID) int {
		return d.list[id]
	}
	return d
}

var Next func(id int) ID
var Some func() int
var getId func(id ID) int
