package eval

import (
	"fw/cp"
	"fw/rt2"
	"fw/rt2/context"
	"github.com/kpmy/ypk/assert"
)

func KeyOf(in IN, key interface{}) cp.ID {
	id, ok := rt2.RegOf(in.Frame)[key].(cp.ID)
	assert.For(ok, 40)
	return id
}

func MetaOf(in IN) (ret *Meta) {
	ret, ok := rt2.RegOf(in.Frame)[context.META].(*Meta)
	assert.For(ok, 60)
	return
}
