package eval

import (
	"fw/cp/constant/statement"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

func doEnter(in IN) OUT {
	var next Do
	tail := func(IN) (out OUT) {
		body := in.IR.Right()
		switch {
		case body == nil:
			return End()
		case body != nil && in.Parent != nil:
			panic(0)
		case body != nil && in.Parent == nil: //секция BEGIN
			rt2.Push(rt2.New(body), in.Frame)
			end := in.IR.Link()
			if end != nil { //секция CLOSE
				out.Do = func(in IN) OUT {
					in.Frame.Root().PushFor(rt2.New(end), in.Frame)
					return OUT{Do: Tail(STOP), Next: LATER}
				}
			} else {
				out.Do = Tail(STOP)
			}
			out.Next = BEGIN
		}
		return
	}
	sm := rt2.ThisScope(in.Frame)
	if in.IR.Object() != nil { //параметры процедуры
		panic(0)
	} else {
		sm.Target().(scope.ScopeAllocator).Allocate(in.IR, true)
		next = tail
	}
	return Now(next)
}

func doAssign(in IN) (out OUT) {
	const (
		right = "assign:right"
		left  = "assign:left"
	)
	a := in.IR.(node.AssignNode)
	switch a.Statement() {
	case statement.ASSIGN:
		out = GetExpression(in, right, a.Right(), func(in IN) OUT {
			id := KeyOf(in, right)
			val := rt2.ValueOf(in.Frame)[id]
			assert.For(val != nil, 40, id)
			return GetDesignator(in, left, a.Left(), func(in IN) OUT {
				id := KeyOf(in, left)
				v, ok := rt2.ValueOf(in.Frame)[id].(scope.Variable)
				assert.For(ok, 41, reflect.TypeOf(v))
				v.Set(val)
				return End()
			})
		})
	default:
		halt.As(100, "unsupported assign statement", a.Statement())
	}
	return
}

func doCall(in IN) OUT {
	panic(0)
}
