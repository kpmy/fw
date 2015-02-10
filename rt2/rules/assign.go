package rules

import (
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"reflect"
	"ypk/halt"
)

func inc_dec_seq(f frame.Frame, code operation.Operation) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	a := node.New(constant.ASSIGN, cp.Some()).(node.AssignNode)
	a.SetStatement(statement.ASSIGN)
	a.SetLeft(n.Left())
	op := node.New(constant.DYADIC, cp.Some()).(node.OperationNode)
	op.SetOperation(code)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	a.SetRight(op)
	rt2.Push(rt2.New(a), f)
	/*seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(n.Left().Object().Adr(), scope.Simple(rt2.ValueOf(f)[op.Adr()]))
		return frame.End()
	}
	ret = frame.LATER */
	return frame.Tail(frame.STOP), frame.LATER
}

/*func decSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)
	op := node.New(constant.DYADIC, cp.Some()).(node.OperationNode)
	op.SetOperation(operation.MINUS)
	op.SetLeft(n.Left())
	op.SetRight(n.Right())
	rt2.Push(rt2.New(op), f)
	seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		sc.Update(n.Left().Object().Adr(), scope.Simple(rt2.ValueOf(f)[op.Adr()]))
		return frame.End()
	}
	ret = frame.LATER
	return seq, ret
}
*/
func assignSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	a := rt2.NodeOf(f)

	var left scope.Value
	var rightId cp.ID

	right := func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
		vleft := left.(scope.Variable)
		switch r := a.Right().(type) {
		case node.ConstantNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := rt2.ThisScope(f)
				vleft.Set(sc.Provide(a.Right())(nil)) //scope.Simple(a.Right().(node.ConstantNode).Data()))
				return frame.End()
			}
			ret = frame.NOW
		case node.VariableNode, node.ParameterNode:
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				sc := rt2.ScopeFor(f, a.Right().Object().Adr())
				vleft.Set(sc.Select(a.Right().Object().Adr()))
				return frame.End()
			}
			ret = frame.NOW
		case node.OperationNode, node.CallNode, node.DerefNode, node.FieldNode:
			rt2.Push(rt2.New(a.Right()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[a.Right().Adr()] != nil || rt2.RegOf(f)["RETURN"] != nil, 61
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
				val := rt2.ValueOf(f)[a.Right().Adr()]
				if val == nil {
					val = rt2.RegOf(f)["RETURN"].(scope.Value)
				}
				vleft.Set(val)
				return frame.End()
			}
			ret = frame.LATER
		case node.IndexNode:
			rightId = r.Adr()
			rt2.Push(rt2.New(r), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[rightId] != nil, 62
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {

				right := rt2.ValueOf(f)[r.Adr()]
				switch z := r.Left().(type) {
				case node.VariableNode, node.ParameterNode:
					sc := rt2.ScopeFor(f, z.Object().Adr())
					arr := sc.Select(z.Object().Adr()).(scope.Array)
					right = arr.Get(right)
					vleft.Set(right)
					return frame.End()
				case node.DerefNode:
					return This(expectExpr(f, z, func(in ...IN) (out OUT) {
						arr := rt2.ValueOf(f)[z.Adr()].(scope.Array)
						right = arr.Get(right)
						vleft.Set(right)
						return End()
					}))
				default:
					halt.As(100, reflect.TypeOf(z), z)
				}
				panic(0)
			}
			ret = frame.LATER
		case node.ProcedureNode:
			sc := rt2.ThisScope(f)
			vleft.Set(sc.Provide(a.Right().Object())(nil))
			return frame.End()
		default:
			fmt.Println(reflect.TypeOf(a.Right()))
			panic("wrong right")
		}
		return seq, ret
	}

	switch a.(node.AssignNode).Statement() {
	case statement.ASSIGN:
		switch l := a.Left().(type) {
		case node.VariableNode, node.ParameterNode:
			sc := rt2.ScopeFor(f, a.Left().Object().Adr())
			left = sc.Select(a.Left().Object().Adr())
			seq, ret = right(f)
		case node.FieldNode:
			rt2.Push(rt2.New(l), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 63
			})
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				left = rt2.ValueOf(f)[l.Adr()]
				return right(f)
			}
			ret = frame.LATER
		case node.IndexNode:
			rt2.Push(rt2.New(a.Left()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				return rt2.ValueOf(f)[l.Adr()] != nil, 64
			})
			seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				left = rt2.ValueOf(f)[l.Adr()]
				switch z := l.Left().(type) {
				case node.VariableNode, node.ParameterNode:
					sc := rt2.ScopeFor(f, l.Left().Object().Adr())
					arr := sc.Select(l.Left().Object().Adr()).(scope.Array)
					left = arr.Get(left)
					return right(f)
				case node.DerefNode:
					return This(expectExpr(f, z, func(in ...IN) (out OUT) {
						arr := rt2.ValueOf(f)[z.Adr()].(scope.Array)
						left = arr.Get(left)
						out.do = Expose(right)
						out.next = NOW
						return
					}))
				default:
					halt.As(100, reflect.TypeOf(z), z)
				}
				panic(0)
			}
			ret = frame.LATER
		case node.DerefNode:
			//			rt2.DataOf(f)[a.Left()] = scope.ID{}
			rt2.Push(rt2.New(a.Left()), f)
			rt2.Assert(f, func(f frame.Frame) (bool, int) {
				ok := rt2.ValueOf(f)[l.Adr()] != nil
				return ok, 65
			})
			seq = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				left = rt2.ValueOf(f)[a.Left().Adr()]
				return right(f)
			}
			ret = frame.LATER
		default:
			fmt.Println(reflect.TypeOf(a.Left()))
			panic("wrong left")
		}
	case statement.INC, statement.INCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode, node.FieldNode:
			seq, ret = inc_dec_seq(f, operation.PLUS)
		default:
			panic(fmt.Sprintln("wrong left", reflect.TypeOf(a.Left())))
		}
	case statement.DEC, statement.EXCL:
		switch a.Left().(type) {
		case node.VariableNode, node.ParameterNode, node.FieldNode:
			seq, ret = inc_dec_seq(f, operation.MINUS)
		default:
			panic(fmt.Sprintln("wrong left", reflect.TypeOf(a.Left())))
		}
	case statement.NEW:
		heap := f.Domain().Discover(context.HEAP).(scope.Manager).Target().(scope.HeapAllocator)
		if a.Right() != nil {
			seq, ret = This(expectExpr(f, a.Right(), func(in ...IN) (out OUT) {
				//fmt.Println("NEW", rt2.ValueOf(f)[a.Right().Adr()], "here")
				switch z := a.Left().(type) {
				case node.VariableNode:
					sc := rt2.ScopeFor(f, a.Left().Object().Adr())
					fn := heap.Allocate(a.Left(), rt2.ValueOf(f)[a.Right().Adr()])
					sc.Update(a.Left().Object().Adr(), fn)
					return End()
				case node.FieldNode:
					fn := heap.Allocate(a.Left(), rt2.ValueOf(f)[a.Right().Adr()])
					rt2.Push(rt2.New(z), in[0].frame)
					rt2.Assert(f, func(f frame.Frame) (bool, int) {
						return rt2.ValueOf(f)[z.Adr()] != nil, 65
					})
					out.do = func(in ...IN) OUT {
						field := rt2.ValueOf(in[0].frame)[z.Adr()].(scope.Variable)
						field.Set(fn(nil))
						return End()
					}
					out.next = LATER
					return
				default:
					halt.As(100, reflect.TypeOf(z))
				}
				panic(0)
			}))
		} else {
			//fmt.Println("NEW here", a.Left().Adr())
			fn := heap.Allocate(a.Left())
			sc := rt2.ScopeFor(f, a.Left().Object().Adr())
			sc.Update(a.Left().Object().Adr(), fn)
			return frame.End()
		}
	default:
		panic(fmt.Sprintln("wrong statement", a.(node.AssignNode).Statement()))
	}
	return seq, ret
}
