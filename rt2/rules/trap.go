package rules

import (
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/traps"
	"fw/rt2"
	"fw/rt2/frame"
)

func doTrap(f frame.Frame, err traps.TRAP) (frame.Sequence, frame.WAIT) {
	trap := node.New(constant.TRAP, cp.Some()).(node.TrapNode)
	code := node.New(constant.CONSTANT, cp.Some()).(node.ConstantNode)
	code.SetData(int32(err))
	trap.SetLeft(code)
	rt2.Push(rt2.New(trap), f)
	return frame.Tail(frame.STOP), frame.LATER
}

func thisTrap(f frame.Frame, err traps.TRAP) (out OUT) {
	out.do = Expose(func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		return doTrap(f, err)
	})
	out.next = LATER
	return
}
