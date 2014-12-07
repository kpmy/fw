package rt

import (
	"cp/node"
	"cp/object"
)

type assignSeq struct {
	step func() Wait
}

func (s *assignSeq) Do(f *frame) (ret Wait) {
	_ = f.ir.(node.AssignNode)
	if s.step == nil {
		ret = DO
		s.step = func() Wait {
			switch f.ir.Left().(type) {
			case node.VariableNode:
				f.ret[f.ir.Left()] = f.ir.Left().Object()
				s.step = func() Wait {
					switch f.ir.Right().(type) {
					case node.ConstantNode:
						f.ret[f.ir.Right()] = f.ir.Right().(node.ConstantNode).Data()
						s.step = func() Wait {
							a := f.p.heap.This(f.ret[f.ir.Left()].(object.Object))
							a.Set(f.ret[f.ir.Right()])
							return STOP
						}
						return DO
					case node.OperationNode:
						nf := NewFrame(f.p, f.ir.Right()).(*frame)
						f.push(nf)
						s.step = func() Wait {
							a := f.p.heap.This(f.ret[f.ir.Left()].(object.Object))
							a.Set(f.ret[f.ir.Right()])
							return STOP
						}
						return SKIP
					default:
						panic("wrong right")
					}
				}
				return DO
			default:
				panic("wrong left")
			}
		}
	} else {
		ret = s.step()
	}
	return ret
}

type opSeq struct {
	step func() Wait
}

func (s *opSeq) Do(f *frame) (ret Wait) {
	_ = f.ir.(node.OperationNode)
	op := func() Wait {
		switch f.ir.(node.OperationNode).Operation() {
		case node.PLUS:
			a := f.ret[f.ir.Left()]
			b := f.ret[f.ir.Right()]
			f.ret[f.ir] = Sum(a, b)
			return STOP
		default:
			panic("unknown operation")
		}
	}
	right := func() Wait {
		switch f.ir.Right().(type) {
		case node.ConstantNode:
			f.ret[f.ir.Right()] = f.ir.Right().(node.ConstantNode).Data()
			s.step = op
			return DO
		case node.VariableNode:
			s.step = func() Wait {
				f.ret[f.ir.Right()] = f.p.heap.This(f.ir.Right().Object())
				s.step = op
				return DO
			}
			return SKIP
		default:
			panic("wrong right")
		}
	}
	left := func() Wait {
		switch f.ir.Left().(type) {
		case node.ConstantNode:
			f.ret[f.ir.Left()] = f.ir.Left().(node.ConstantNode).Data()
			s.step = right
			return DO
		case node.VariableNode:
			s.step = func() Wait {
				f.ret[f.ir.Left()] = f.p.heap.This(f.ir.Left().Object())
				s.step = right
				return DO
			}
			return SKIP
		default:
			panic("wrong left")
		}
	}
	if s.step == nil {
		s.step = left
		ret = DO
	} else {
		ret = s.step()
	}
	return ret
}
