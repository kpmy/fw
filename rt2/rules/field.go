package rules

import (
	"reflect"
)

import (
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"ypk/halt"
)

func fieldSeq(in ...IN) (out OUT) {
	f := in[0].frame
	n := rt2.NodeOf(f)
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	switch l := n.Left().(type) {
	case node.VariableNode:
		sc.Select(l.Object().Adr(), func(v scope.Value) {
			rt2.ValueOf(f.Parent())[n.Adr()] = v.(scope.Record).Get(n.Object().Adr())
		})
		out = End()
	case node.FieldNode:
		rt2.Push(rt2.New(l), f)
		rt2.Assert(f, func(f frame.Frame) (bool, int) {
			_, ok := rt2.ValueOf(f)[l.Adr()].(scope.Record)
			return ok, 60
		})
		out.do = func(in ...IN) OUT {
			v := rt2.ValueOf(in[0].frame)[l.Adr()].(scope.Record)
			rt2.ValueOf(f.Parent())[n.Adr()] = v.Get(n.Object().Adr())
			return End()
		}
		out.next = LATER
	default:
		halt.As(100, reflect.TypeOf(l))
	}
	return out
}