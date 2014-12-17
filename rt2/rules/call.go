package rules

import (
	"cp/node"
	"rt2/frame"
	mod "rt2/module"
	"rt2/nodeframe"
)

func callSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	switch n.Left().(type) {
	case node.ProcedureNode:
		m := mod.DomainModule(f.Domain())
		proc := m.NodeByObject(n.Left().Object())
		f.Root().Push(fu.New(proc))
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			return frame.End()
		}
		ret = frame.SKIP
	default:
		panic("unknown call left")
	}
	return seq, ret
}
