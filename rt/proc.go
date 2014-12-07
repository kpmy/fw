package rt

import (
	"cp/node"
	"errors"
	"fmt"
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

type Sequence interface {
	Do(f *frame) Wait
}

func NewProcessor() Processor {
	return new(procImpl).Init()
}

type frame struct {
	p      *procImpl
	parent *frame
	ir     node.Node
	seq    Sequence
	ret    map[node.Node]interface{}
}

func (f *frame) Do() (wait Wait) {
	if f.seq == nil {
		panic("no sequence")
	}
	return f.seq.Do(f)
}

func (f *frame) OnPush() {
	switch f.ir.(type) {
	case node.AssignNode:
		f.ret = make(map[node.Node]interface{}, 2)
		f.seq = new(assignSeq)
	case node.OperationNode:
		f.ret = make(map[node.Node]interface{}, 3)
		f.seq = new(opSeq)
	default:
		panic("unknown ir")
	}
}

func (f *frame) OnPop() {
	switch f.ir.(type) {
	case node.AssignNode:
		if f.ir.Link() != nil {
			f.p.stack.Push(NewFrame(f.p, f.ir.Link()))
		}
	case node.OperationNode:
		f.parent.ret[f.ir] = f.ret[f.ir]
	}
}

func (f *frame) push(t *frame) {
	t.parent = f
	f.p.stack.Push(t)
}

func NewFrame(p *procImpl, ir node.Node) Frame {
	f := new(frame)
	f.ir = ir
	f.p = p
	return f
}

type procImpl struct {
	stack Stack
	heap  Heap
	cycle int64
}

func (p *procImpl) Init() *procImpl {
	p.stack = NewStack()
	p.heap = NewHeap()
	return p
}

func (p *procImpl) ConnectTo(head node.Node) (err error) {
	if head != nil {
		switch head.(type) {
		// особый случай, после enter вправо, а не вниз
		case node.EnterNode:
			p.stack.Push(NewFrame(p, head.Right()))
		default:
			panic("oops")
		}
	} else {
		err = errors.New("not a head node")
	}
	return err
}

func (p *procImpl) Do() (res Result, err error) {
	if p.stack.Top() != nil {
		fmt.Println(p.cycle)
		p.cycle++
		f := p.stack.Top()
		//цикл дейкстры
		for {
			wait := f.Do()
			fmt.Println(wait)
			if wait == SKIP {
				break
			} else if wait == DO {
			} else if wait == WRONG {
				panic("something wrong")
			} else {
				if f == p.stack.Top() {
					p.stack.Pop()
				} else {
					panic("do not stop if not top on stack")
				}
				break
			}
		}
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
