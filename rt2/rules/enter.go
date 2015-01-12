package rules

import (
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
)

func enterSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)
	body := fu.NodeOf(f).Right()
	tail := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		if body == nil {
			//случай пустого тела процедуры/секции BEGIN
			return frame.End()
		} else if f.Parent() != nil {
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
	sm := scope.This(f.Domain().Discover(context.SCOPE))
	//fmt.Println(n.Object())
	if n.Object() != nil {
		par, ok := fu.DataOf(f)[n.Object()].(node.Node)
		//fmt.Println(fu.DataOf(f)[n.Object()])
		//fmt.Println(ok)
		if ok {
			sm.Target().(scope.ScopeAllocator).Allocate(n, false)
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				return sm.Target().(scope.ScopeAllocator).Initialize(n, scope.PARAM{Objects: n.Object().Link(), Values: par, Frame: f, Tail: tail})
			}
		} else {
			sm.Target().(scope.ScopeAllocator).Allocate(n, true)
			seq = tail
		}
	} else {
		sm.Target().(scope.ScopeAllocator).Allocate(n, true)
		seq = tail
	}
	return seq, frame.NOW
}
