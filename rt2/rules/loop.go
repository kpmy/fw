package rules

import (
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
)

const flag = 0

func exitSeq(x frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	x.Root().ForEach(func(f frame.Frame) (ok bool) {
		n := rt2.Utils.NodeOf(f)
		_, ok = n.(node.LoopNode)
		if ok {
			rt2.Utils.DataOf(f)[flag] = true
		}
		ok = !ok
		return ok
	})
	return frame.End()
}

func loopSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.Utils.NodeOf(f)
	exit, ok := rt2.Utils.DataOf(f)[flag].(bool)
	if ok && exit {
		return frame.End()
	}
	if n.Left() != nil {
		rt2.Utils.Push(rt2.Utils.New(n.Left()), f)
		return loopSeq, frame.LATER
	} else if n.Left() == nil {
		return frame.End()
	} else {
		panic("unexpected loop seq")
	}

}
