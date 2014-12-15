package rules

import (
	"cp/node"
	"rt2/frame"
	"rt2/nodeframe"
)

func callSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	switch n.Left().(type) {
	case node.ProcedureNode:
		//proc := f.p.thisMod.NodeByObject(f.ir.Left().Object())
		//f.Root().Push(fu.New(?))
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			return frame.End()
		}
		ret = frame.STOP
		seq = nil //frame.SKIP //uncomment when truly work
	default:
		panic("unknown call left")
	}
	return seq, ret
}
