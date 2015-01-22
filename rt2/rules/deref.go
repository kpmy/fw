package rules

import (
	"fmt"
	//	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

func derefSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f).(node.DerefNode)
	sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
	fmt.Println("deref from ptr", n.Ptr())
	if n.Ptr() {
		switch l := n.Left().Object().(type) {
		case object.ParameterObject, object.VariableObject:
			sc.Select(l.Adr(), func(v scope.Value) {
				ptr, ok := v.(scope.Pointer)
				assert.For(ok, 60, reflect.TypeOf(v))
				rt2.ValueOf(f.Parent())[n.Adr()] = ptr.Get()
			})
			return frame.End()
		default:
			halt.As(100, l.Adr(), reflect.TypeOf(l))
		}
	} else {
		if n.Left().Object() != nil {
			switch l := n.Left().Object().(type) {
			case object.ParameterObject, object.VariableObject:
				rt2.ValueOf(f.Parent())[n.Adr()] = sc.Select(l.Adr())
				return frame.End()
			default:
				halt.As(100, l.Adr(), reflect.TypeOf(l))
			}
		} else {
			switch left := n.Left().(type) {
			case node.DerefNode:
				rt2.Push(rt2.New(left), f)
				rt2.Assert(f, func(f frame.Frame) (bool, int) {
					return rt2.ValueOf(f)[left.Adr()] != nil, 60
				})
				seq = Propose(func(...IN) OUT {
					rt2.ValueOf(f.Parent())[n.Adr()] = rt2.ValueOf(f)[left.Adr()]
					return End()
				})
				ret = LATER.wait()
			default:
				halt.As(100, reflect.TypeOf(left))
			}
		}
	}
	/* if ok {
		//		rt2.DataOf(f.Parent())[n] = scope.Designator(n.Left())
	} else {
		for z := n.Left(); !ok && z != nil; {
			switch z.(type) {
			case node.DerefNode:
				z = z.Left()
			default:
				data := sc.Select(z.Adr())
				assert.For(data != nil, 40)
				rt2.DataOf(f.Parent())[n] = data
				ok = true
			}
		}
	} */
	return seq, ret
}
