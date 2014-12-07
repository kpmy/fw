package rt

import (
	"cp/node"
	"fmt"
	"reflect"
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
							fmt.Println(reflect.TypeOf(f.ret[f.ir.Right()]))
							fmt.Println("присвоение константы")
							return STOP
						}
						return DO
					case node.OperationNode:
						nf := NewFrame(f.p, f.ir.Right()).(*frame)
						f.push(nf)
						s.step = func() Wait {
							fmt.Println(reflect.TypeOf(f.ret[f.ir.Right()]))
							fmt.Println("присвоение результата операции")
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
			f.ret[f.ir] = 0 //a + b
			fmt.Println("сложение")
			fmt.Println(reflect.TypeOf(a), reflect.TypeOf(b))
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
				f.ret[f.ir.Right()] = f.p.heap.ThisVariable(f.ir.Right().Object())
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
				f.ret[f.ir.Left()] = f.p.heap.ThisVariable(f.ir.Left().Object())
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

/*
func (p *procImpl) doStat() {
	f := p.stack.Top().(*frame)
	fmt.Println(reflect.TypeOf(f.ir))
	switch f.ir.(type) {
	//assign works like .left := .right
	case node.AssignNode:
		if f.prologue != nil {
			f.ret = f.ir.Link()
			switch f.ir.Left().(type) {
			case node.VariableNode: //nothing to do
			default:
				panic("left is not variable")
			}
			switch f.ir.Right().(type) {
			case node.ConstantNode:
				x := p.heap.ThisVariable(f.ir.Left().Object())
				_ = x
				*x.(*int) = f.ir.Right().(node.ConstantNode).Data().(int)
				fmt.Println(x)
				x = p.heap.ThisVariable(f.ir.Left().Object())
				fmt.Println(*x.(*int))
			case node.OperationNode:
				nf := new(frame)
				nf.ir = f.ir.Right()
				nf.ret = f.ir
				nf.prologue = prologue
				nf.epilogue = node.New(node.CONSTANT)
				f.epilogue = nf.epilogue
				p.stack.Push(nf)
			default:
				panic("unknown right assign")
			}
			f.prologue = nil
		} else if f.epilogue != nil {
			switch f.epilogue.(type) {
			case node.ConstantNode:
				x := p.heap.ThisVariable(f.ir.Left().Object())
				_ = x
				*x.(*int) = f.epilogue.(node.ConstantNode).Data().(int)
				fmt.Println(*x.(*int))
			default:
				fmt.Println("no custom epilogue")
			}
			p.stack.Pop()
			if f.ret != nil {
				nf := new(frame)
				nf.ir = f.ret
				nf.prologue = prologue
				nf.epilogue = prologue
				p.stack.Push(nf)
			}

		}
	case node.OperationNode:
		switch f.ir.(node.OperationNode).Operation() {
		case node.PLUS:
			var a, b int
			switch f.ir.Left().(type) {
			case node.ConstantNode:
				a = f.ir.Left().(node.ConstantNode).Data().(int)
			case node.VariableNode:
				x := p.heap.ThisVariable(f.ir.Left().Object())
				a = *x.(*int)
			default:
				panic("unknown left operand")
			}
			switch f.ir.Right().(type) {
			case node.ConstantNode:
				b = f.ir.Right().(node.ConstantNode).Data().(int)
			case node.VariableNode:
				x := p.heap.ThisVariable(f.ir.Right().Object())
				b = *x.(*int)
			default:
				panic("unknown right operand")
			}
			switch f.epilogue.(type) {
			case node.ConstantNode:
				fmt.Println(a, b)
				f.epilogue.(node.ConstantNode).SetData(a + b)
				p.stack.Pop()
			default:
				panic("unknown epilogue")
			}
		default:
			panic("unknown operation")
		}
	default:
		panic("ooops")
	}
}
*/
