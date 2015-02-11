package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"reflect"
)

func returnSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	a := rt2.NodeOf(f)

	left := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		switch a.Left().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//rt2.DataOf(f.Parent())[a.Object()] = a.Left().(node.ConstantNode).Data()
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = rt2.ThisScope(f).Provide(a.Left())(nil)
				rt2.RegOf(f.Parent())[context.RETURN] = rt2.ValueOf(f.Parent())[a.Object().Adr()]
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := rt2.ScopeFor(f, a.Left().Object().Adr())
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = sc.Select(a.Left().Object().Adr())
				rt2.RegOf(f.Parent())[context.RETURN] = rt2.ValueOf(f.Parent())[a.Object().Adr()]
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode, node.FieldNode:
			rt2.Push(rt2.New(a.Left()), f)
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = rt2.ValueOf(f)[a.Left().Adr()]
				rt2.RegOf(f.Parent())[context.RETURN] = rt2.ValueOf(f)[a.Left().Adr()]
				if rt2.RegOf(f.Parent())[context.RETURN] == nil {
					rt2.RegOf(f.Parent())[context.RETURN] = rt2.RegOf(f)[context.RETURN]
				}
				return frame.End()
			}
			ret = frame.LATER
		case node.IndexNode:
			return This(expectExpr(f, a.Left(), func(...IN) OUT {
				rt2.ValueOf(f.Parent())[a.Object().Adr()] = rt2.ValueOf(f)[a.Left().Adr()]
				rt2.RegOf(f.Parent())[context.RETURN] = rt2.ValueOf(f)[a.Left().Adr()]
				if rt2.RegOf(f.Parent())[context.RETURN] == nil {
					rt2.RegOf(f.Parent())[context.RETURN] = rt2.RegOf(f)[context.RETURN]
				}
				return End()
			}))
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong left")
		}
		return seq, ret
	}
	return left, frame.NOW
}
