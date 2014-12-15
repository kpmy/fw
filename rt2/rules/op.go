package rules

import (
	"cp/constant/operation"
	"cp/node"
	"rt2/frame"
	"rt2/nodeframe"
)

func opSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils

	op := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.(node.OperationNode).Operation() {
		case operation.PLUS:
			//складываем
			return frame.End()
		default:
			panic("unknown operation")
		}
	}

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Right().(type) {
		case node.ConstantNode:
			//f.ret[f.ir.Right()] = f.ir.Right().(node.ConstantNode).Data()
			return op, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//f.ret[f.ir.Right()] = f.p.heap.This(f.ir.Right().Object())
				return op, frame.DO
			}
			ret = frame.DO
			return seq, ret
		default:
			panic("wrong right")
		}
	}

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		n := fu.NodeOf(f)
		switch n.Left().(type) {
		case node.ConstantNode:
			//f.ret[f.ir.Left()] = f.ir.Left().(node.ConstantNode).Data()
			return right, frame.DO
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//f.ret[f.ir.Left()] = f.p.heap.This(f.ir.Left().Object())
				return right, frame.DO
			}
			ret = frame.DO
			return seq, ret
		default:
			panic("wrong left")
		}
	}

	return left, frame.DO
}
