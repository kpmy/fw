package rt

import (
	"cp/node"
	"errors"
	"fmt"
	"reflect"
)

type Result int

const (
	OK Result = iota
	END
	ERROR
)

type Processor interface {
	ConnectTo(head node.Node) error
	Do() (Result, error)
}

func NewProcessor() Processor {
	return new(procImpl).Init()
}

type frame struct {
	Frame
	prologue, epilogue interface{}
	ir, ret            node.Node
}

type procImpl struct {
	stack Stack
	heap  Heap
}

func (p *procImpl) Init() *procImpl {
	p.stack = NewStack()
	p.heap = NewHeap()
	return p
}

func (p *procImpl) ConnectTo(head node.Node) (err error) {
	if head != nil {
		switch head.(type) {
		case node.EnterNode:
			f := new(frame)
			f.ir = head.Right()
			f.prologue = prologue
			f.epilogue = prologue
			p.stack.Push(f)
		default:
			panic("oops")
		}
	} else {
		err = errors.New("not a head node")
	}
	return err
}

func prologue() {

}

func (p *procImpl) doExpression() {
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

func (p *procImpl) Do() (res Result, err error) {
	if p.stack.Top() != nil {
		p.doExpression()
	} else {
		err = errors.New("no program")
	}
	if p.stack.Top() != nil {
		res = OK
	} else {
		res = END
		fmt.Println(p.heap)
	}
	return res, err
}
