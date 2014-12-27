package rules

import (
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
)

func derefSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.Utils.NodeOf(f)
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	rt2.Utils.DataOf(f.Parent())[n] = sc.SelectObj(n.Left().Object())
	return frame.End()
}
