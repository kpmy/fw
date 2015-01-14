package std

import (
	"container/list"
	"fmt"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"ypk/assert"
)

type FlowFrame interface {
}

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

func (f *RootFrame) PushFor(fr, parent frame.Frame) {
	_, ok := fr.(*RootFrame)
	if ok {
		panic("impossibru")
	}
	f.inner.PushFront(fr)
	if parent == nil {
		domain := f.Domain().(context.Factory).New()
		domain.Attach(context.SCOPE, scope.NewStack())
		fr.Init(domain)
	} else {
		fr.Init(parent.Domain())
	}
	fr.OnPush(f, parent)
}

func (f *RootFrame) Pop() {
	if f.inner.Front() != nil {
		elem := f.inner.Front()
		frame := elem.Value.(frame.Frame)
		frame.OnPop()
		f.inner.Remove(elem)
	} else {
		panic("it's empty stack")
	}
}

func (f *RootFrame) Top() (fr frame.Frame) {
	elem := f.inner.Front()
	if elem != nil {
		fr = elem.Value.(frame.Frame)
	}
	return fr
}

func (f *RootFrame) Drop() (fr frame.Frame) {
	elem := f.inner.Front()
	if elem != nil {
		fr = elem.Value.(frame.Frame)
		f.inner.Remove(elem)
	}
	return fr
}

func (f *RootFrame) Do() (res frame.WAIT) {
	var (
		trapped bool
	)
	if f.Top() != nil {
		x := f.Top()
		//цикл дейкстры
		for {
			wait := x.Do()
			//fmt.Println(wait)
			if wait == frame.LATER || wait == frame.BEGIN || wait == frame.END {
				break
			} else if wait == frame.NOW {
			} else if wait == frame.WRONG {
				trapped = true
				fmt.Println("it's a trap")
				break
				//panic("something wrong") do nothing, it's a trap
			} else if wait == frame.STOP {
				if x == f.Top() {
					f.Pop()
				} else {
					panic("do not stop if not top on stack")
				}
				break
			} else {
				panic("wrong wait code")
			}
		}
	}
	if f.Top() != nil && !trapped {
		res = frame.NOW
	} else {
		res = frame.STOP
	}
	return res
}

func (f *RootFrame) ForEach(run func(x frame.Frame) bool) {
	e := f.inner.Front()
	ok := true
	for (e != nil) && ok {
		ok = run(e.Value.(frame.Frame))
		e = e.Next()
	}
}

func (f *RootFrame) OnPush(a frame.Stack, b frame.Frame) {}
func (f *RootFrame) OnPop()                              {}
func (f *RootFrame) Parent() frame.Frame                 { return nil }
func (f *RootFrame) Root() frame.Stack                   { return nil }
func (f *RootFrame) Domain() context.Domain              { return f.domain }
func (f *RootFrame) Init(d context.Domain) {
	assert.For(f.domain == nil, 20)
	assert.For(d != nil, 21)
	f.domain = d
}

func (f *RootFrame) Handle(msg interface{}) {}
