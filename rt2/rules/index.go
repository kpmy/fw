package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"reflect"
)

func indexSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	i := rt2.NodeOf(f)

	switch i.Right().(type) {
	case node.ConstantNode:
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.ValueOf(f.Parent())[i.Adr()] = rt2.ThisScope(f).Provide(i.Right())(nil)
			return frame.End()
		}
		ret = frame.NOW
	case node.VariableNode:
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			sc := rt2.ScopeFor(f, i.Right().Object().Adr())
			rt2.ValueOf(f.Parent())[i.Adr()] = sc.Select(i.Right().Object().Adr())
			return frame.End()
		}
		ret = frame.NOW
	case node.OperationNode, node.CallNode:
		rt2.Push(rt2.New(i.Right()), f)
		rt2.Assert(f, func(f frame.Frame) (bool, int) {
			return rt2.ValueOf(f)[i.Right().Adr()] != nil, 60
		})
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			rt2.ValueOf(f.Parent())[i.Adr()] = rt2.ValueOf(f)[i.Right().Adr()]
			return frame.End()
		}
		ret = frame.LATER
	default:
		panic(fmt.Sprintln("unsupported type", reflect.TypeOf(i.Right())))
	}
	return seq, ret
}
