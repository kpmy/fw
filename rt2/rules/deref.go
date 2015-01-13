package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"ypk/assert"
)

func derefSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	fmt.Println("deref from ptr", n.(node.DerefNode).Ptr())
	_, ok := rt2.DataOf(f.Parent())[n].(scope.ID)
	if ok {
		rt2.DataOf(f.Parent())[n] = scope.Designator(n.Left())
	} else {
		data := sc.Select(scope.Designator(n.Left()))
		assert.For(data != nil, 40)
		rt2.DataOf(f.Parent())[n] = data
	}
	return frame.End()
}
