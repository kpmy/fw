package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/module"
	"fw/rt2/scope"
	"fw/utils"
)

func enterSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	body := rt2.NodeOf(f).Right()
	tail := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		if body == nil {
			//случай пустого тела процедуры/секции BEGIN
			return frame.End()
		} else if f.Parent() != nil {
			//Вход в процедуру не несет значимых действий и просто заменяет себя в цепочке родителей на своего родителя
			//При вызове фрейма с другим доменом это мешает, надо убрать
			//Через DataOf(f.Parent()) может передаваться результат выполнения
			//rt2.Push(rt2.New(body), f.Parent())
			rt2.Push(rt2.New(body), f)
			return frame.Tail(frame.STOP), frame.LATER
		} else {
			//Особый случай, вход в модуль, секция BEGIN
			rt2.Push(rt2.New(body), f)
			//fmt.Println("begin", module.DomainModule(f.Domain()).Name)
			//Выход из модуля, секция CLOSE
			next := n.Link()
			if next != nil {
				seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					utils.PrintFrame("end", module.DomainModule(f.Domain()).Name)
					f.Root().PushFor(rt2.New(next), f)
					return frame.Tail(frame.STOP), frame.END
				}
			} else {
				seq = frame.Tail(frame.STOP)
			}
			return seq, frame.BEGIN

		}
	}
	sm := f.Domain().Discover(context.SCOPE).(scope.Manager)
	//fmt.Println(n.Object())
	if n.Object() != nil {
		par, ok := rt2.RegOf(f)[n.Object()].(node.Node)
		//fmt.Println(rt2.DataOf(f)[n.Object()])
		//fmt.Println(ok)
		if ok {
			sm.Target().(scope.ScopeAllocator).Allocate(n, false)
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				return sm.Target().(scope.ScopeAllocator).Initialize(n,
					scope.PARAM{
						Objects: n.Object().Link(),
						Values:  par,
						Frame:   f,
						Tail:    tail})
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
