package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
)

func caseSeq(f frame.Frame) (frame.Sequence, frame.WAIT) {
	n := rt2.NodeOf(f)
	var e int

	in := func(in ...IN) (out OUT) {
		cond := n.Right().(node.ElseNode)
		fmt.Println("case?", e, cond.Min(), cond.Max())
		if e < cond.Min() || e > cond.Max() { //case?
			if cond.Right() != nil {
				rt2.Push(rt2.New(cond.Right()), f)
			}
			out.do = Tail(STOP)
			out.next = LATER
		} else {
			for next := cond.Left(); next != nil && out.do == nil; next = next.Link() {
				var ok bool
				//				_c := next.Left()
				for _c := next.Left(); _c != nil && !ok; _c = _c.Link() {
					c := _c.(node.ConstantNode)
					fmt.Println("const", c.Data(), c.Min(), c.Max())
					if (c.Min() != nil) && (c.Max() != nil) {
						fmt.Println(e, *c.Max(), "..", *c.Min())
						ok = e >= *c.Min() && e <= *c.Max()
					} else {
						fmt.Println(e, c.Data())
						ok = int32Of(c.Data()) == int32(e)
					}
				}
				fmt.Println(ok)
				if ok {
					rt2.Push(rt2.New(next.Right()), f)
					out.do = Tail(STOP)
					out.next = LATER
				}
			}
			if out.do == nil && cond.Right() != nil {
				rt2.Push(rt2.New(cond.Right()), f)
				out.do = Tail(STOP)
				out.next = LATER
			}
		}
		assert.For(out.do != nil, 60)
		return out
	}

	return This(expectExpr(f, n.Left(), func(...IN) (out OUT) {
		_x := scope.GoTypeFrom(rt2.ValueOf(f)[n.Left().Adr()])
		switch x := _x.(type) {
		case nil:
			panic("nil")
		case int32:
			e = int(x)
			fmt.Println("case", e)
			out.do = in
			out.next = NOW
			return out
		default:
			panic(fmt.Sprintln("unsupported case expr", reflect.TypeOf(_x)))
		}
	}))
}
