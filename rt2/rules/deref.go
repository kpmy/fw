package rules

import (
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"ypk/assert"
)

func derefSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.Utils.NodeOf(f)
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	data := sc.Select(scope.Id(n.Left().Object()))
	assert.For(data != nil, 40)
	rt2.Utils.DataOf(f.Parent())[n] = data
	return frame.End()
}
