package rules

import (
	"cp/node"
	"cp/statement"
	"rt2/frame"
	"rt2/nodeframe"
)

func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	a := fu.NodeOf(f)
	switch a.(node.AssignNode).Statement() {
	case statement.ASSIGN:
		switch a.Left().(type) {
		case node.VariableNode:
			switch a.Right().(type) {
			case node.ConstantNode:
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					//тут присвоение
					return frame.End()
				}
				ret = frame.DO
			case node.VariableNode:
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					//тут присвоение
					return frame.End()
				}
				ret = frame.DO
			case node.OperationNode:
				fu.Push(fu.New(a.Right()), f)
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					//тут чтение результата операции
					return frame.End()
				}
				ret = frame.SKIP
			default:
				panic("wrong right")
			}
		default:
			panic("wrong left")
		}
	default:
		panic("wrong statement")
	}
	return seq, ret
}
