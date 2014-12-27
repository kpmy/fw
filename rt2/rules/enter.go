package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"ypk/assert"
)

func enterSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	body := fu.NodeOf(f).Right()
	assert.For(body != nil, 40)
	sm := scope.This(f.Domain().Discover(context.SCOPE))
	sm.Allocate(n)
	fmt.Println(n.Object())
	if n.Object() != nil {
		par, ok := fu.DataOf(f)[n.Object()].(node.Node)
		//fmt.Println(fu.DataOf(f)[n.Object()])
		//fmt.Println(ok)
		if ok {
			sm.Initialize(n, n.Object().Link(), par)
		}
	}
	if f.Parent() != nil {
		//Вход в процедуру не несет значимых действий и просто заменяет себя в цепочке родителей на своего родителя
		fu.Push(fu.New(body), f.Parent())
		return frame.Tail(frame.STOP), frame.LATER
	} else {
		//Особый случай, вход в модуль, секция BEGIN
		fu.Push(fu.New(body), f)

		//Выход из модуля, секция CLOSE
		next := n.Link()
		if next != nil {
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				f.Root().PushFor(fu.New(next), f)
				return frame.Tail(frame.STOP), frame.LATER
			}
		} else {
			seq = frame.Tail(frame.STOP)
		}
		return seq, frame.LATER

	}
}
