package frame

import (
	"container/list"
	"fmt"
	"rt2/context"
	"ypk/assert"
)

type RootFrame struct {
	inner  list.List
	domain context.Domain
}

func (f *RootFrame) init() *RootFrame {
	f.inner = *list.New()
	return f
}

func NewRoot() *RootFrame {
	return new(RootFrame).init()
}

func (f *RootFrame) Push(frame Frame) {
	f.PushFor(frame, nil)
}

func (f *RootFrame) PushFor(frame, parent Frame) {
	_, ok := frame.(*RootFrame)
	if ok {
		panic("impossibru")
	}
	f.inner.PushFront(frame)
	frame.Init(f.Domain())
	frame.OnPush(f, parent)
}

func (f *RootFrame) Pop() {
	if f.inner.Front() != nil {
		elem := f.inner.Front()
		frame := elem.Value.(Frame)
		frame.OnPop()
		f.inner.Remove(elem)
	} else {
		panic("it's empty stack")
	}
}

func (f *RootFrame) Top() (frame Frame) {
	elem := f.inner.Front()
	if elem != nil {
		frame = elem.Value.(Frame)
	}
	return frame
}

func (f *RootFrame) Do() (res WAIT) {
	if f.Top() != nil {
		x := f.Top()
		//цикл дейкстры
		for {
			wait := x.Do()
			fmt.Println(wait)
			if wait == SKIP {
				break
			} else if wait == DO {
			} else if wait == WRONG {
				panic("something wrong")
			} else {
				if x == f.Top() {
					f.Pop()
				} else {
					panic("do not stop if not top on stack")
				}
				break
			}
		}
	}
	if f.Top() != nil {
		res = DO
	} else {
		res = STOP
	}
	return res
}

func (f *RootFrame) OnPush(a Stack, b Frame) {}
func (f *RootFrame) OnPop()                  {}
func (f *RootFrame) Parent() Frame           { return nil }
func (f *RootFrame) Root() Stack             { return nil }
func (f *RootFrame) Domain() context.Domain  { return f.domain }
func (f *RootFrame) Init(d context.Domain) {
	assert.For(f.domain == nil, 20)
	assert.For(d != nil, 21)
	f.domain = d
}
func (f *RootFrame) Handle(msg interface{}) {}

func (w WAIT) String() string {
	switch w {
	case DO:
		return "DO"
	case SKIP:
		return "SKIP"
	case STOP:
		return "STOP"
	case WRONG:
		return "WRONG"
	default:
		panic("wrong wait value")
	}
}
