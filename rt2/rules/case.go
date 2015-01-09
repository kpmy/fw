package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"reflect"
	"ypk/assert"
)

func caseSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	var e int

	in := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		cond := n.Right().(node.ElseNode)
		if e < cond.Min() || e > cond.Max() { //case?
			seq = frame.Tail(frame.STOP)
			ret = frame.NOW
		} else {
			for next := cond.Left(); next != nil && seq == nil; next = next.Link() {
				var ok bool
				for _c := next.Left(); _c != nil && !ok; _c = _c.Right() {
					c := _c.(node.ConstantNode)
					if (c.Min() != nil) && (c.Max() != nil) {
						//	fmt.Println(e, *c.Max(), *c.Min())
						ok = e >= *c.Min() && e <= *c.Max()
					} else {
						//	fmt.Println(e, c.Data())
						ok = int32Of(c.Data()) == int32(e)
					}
				}
				//fmt.Println(ok)
				if ok {
					rt2.Push(rt2.New(next.Right()), f)
					seq = frame.Tail(frame.STOP)
					ret = frame.LATER
				}
			}
			if seq == nil && cond.Right() != nil {
				rt2.Push(rt2.New(cond.Right()), f)
				seq = frame.Tail(frame.STOP)
				ret = frame.LATER
			}
		}
		assert.For(seq != nil, 60)
		return seq, ret
	}

	return expectExpr(f, n.Left(), func(f frame.Frame) (frame.Sequence, frame.WAIT) {
		_x := rt2.DataOf(f)[n.Left()]
		switch x := _x.(type) {
		case nil:
			panic("nil")
		case int32:
			e = int(x)
			return in, frame.NOW
		default:
			panic(fmt.Sprintln("unsupported case expr", reflect.TypeOf(_x)))
		}
	})
}
