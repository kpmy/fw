package eval

import (
	"fw/cp"
	"fw/rt2"
	"ypk/assert"
)

func KeyOf(in IN, key interface{}) cp.ID {
	id, ok := rt2.RegOf(in.Frame)[key].(cp.ID)
	assert.For(ok, 40)
	return id
}
