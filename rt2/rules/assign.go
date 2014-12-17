package rules

import (
	"cp/node"
	"cp/statement"
	"rt2/context"
	"rt2/frame"
	"rt2/nodeframe"
	"rt2/scope"
)

func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	a := fu.NodeOf(f)
	switch a.(node.AssignNode).Statement() {
	case statement.ASSIGN:
		m := new(frame.SetDataMsg)
		m.Data = make([]interface{}, 1)
		f.(context.ContextAware).Handle(m)
		switch a.Left().(type) {
		case node.VariableNode:
			switch a.Right().(type) {
			case node.ConstantNode:
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
					sc.Update(a.Left().Object(), func(interface{}) interface{} {
						return a.Right().(node.ConstantNode).Data()
					})
					return frame.End()
				}
				ret = frame.DO
			case node.VariableNode:
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
					sc.Update(a.Left().Object(), func(interface{}) interface{} {
						return sc.Select(a.Right().Object())
					})
					return frame.End()
				}
				ret = frame.DO
			case node.OperationNode:
				fu.Push(fu.New(a.Right()), f)
				seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
					sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
					sc.Update(a.Left().Object(), func(interface{}) interface{} {
						return fu.DataOf(f)[0]
					})
					return frame.End()
				}
				ret = frame.SKIP
			default:
				panic("wrong right")
			}
		default:
			panic("wrong left")
		}
	default:
		panic("wrong statement")
	}
	return seq, ret
}
