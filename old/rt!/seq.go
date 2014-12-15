package rt

import (
	"cp/constant/operation"
	"cp/node"
	"cp/object"
	"cp/statement"
)

type assignSeq struct {
	step func() Wait
}

func (s *assignSeq) Do(f *frame) (ret Wait) {
	a := f.ir.(node.AssignNode)
	assign := func() Wait {
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
				case node.VariableNode:
					a := f.p.heap.This(f.ir.Right().Object())
					b := f.p.heap.This(f.ret[f.ir.Left()].(object.Object))
					b.Set(a)
					return STOP
				default:
					panic("wrong right")
				}
			}
			return DO
		default:
			panic("wrong left")
		}
	}
	if s.step == nil {
		ret = DO
		switch a.Statement() {
		case statement.ASSIGN:
			s.step = assign
		default:
			panic("unknown assign subclass")
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
		case operation.PLUS:
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

type callSeq struct {
	step func() Wait
}

func (s *callSeq) Do(f *frame) Wait {
	if s.step == nil {
		switch f.ir.Left().(type) {
		case node.ProcedureNode:
			proc := f.p.thisMod.NodeByObject(f.ir.Left().Object())
			f.push(NewFrame(f.p, proc).(*frame))
			s.step = func() Wait {
				return STOP
			}
			return SKIP
		default:
			panic("unknown call left")
		}
	} else {
		return s.step()
	}
}

type enterSeq struct {
	step func() Wait
}

func (e *enterSeq) Do(f *frame) Wait {
	if e.step == nil {
		//for f.ir == EnterNode entering to .Right()
		f.push(NewFrame(f.p, f.ir.Right()).(*frame))
		e.step = func() Wait {
			return STOP
		}
		return SKIP
	} else {
		return e.step()
	}
}
