package wrap

import (
	"fw/cp/node"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/rt2/rules2/wrap/eval"
	"reflect"
	"ypk/halt"
)

func init() {
	decision.PrologueFor = prologue
	decision.EpilogueFor = epilogue
	decision.AssertFor = test
}

func This(o eval.OUT) (seq frame.Sequence, ret frame.WAIT) {
	ret = o.Next.Wait()
	if ret != frame.STOP {
		seq = Propose(o.Do)
	}
	return seq, ret
}

func Propose(a eval.Do) frame.Sequence {
	return func(fr frame.Frame) (frame.Sequence, frame.WAIT) {
		return This(a(eval.IN{}))
	}
}

func test(n node.Node) (bool, int) {
	switch n.(type) {
	case node.ConstantNode:
		return false, -1
	default:
		return true, 0
	}
	panic(0)
}

func prologue(n node.Node) frame.Sequence {
	/*	//fmt.Println(reflect.TypeOf(n))
		switch next := n.(type) {
		case node.EnterNode:
			var tail frame.Sequence
			tail = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				q := f.Root().Queue()
				if q != nil {
					f.Root().PushFor(q, nil)
					return tail, frame.NOW
				} else {
					return enterSeq, frame.NOW
				}
			}
			return tail
		case node.AssignNode:
			return assignSeq
		case node.OperationNode:
			switch n.(type) {
			case node.DyadicNode:
				return dopSeq
			case node.MonadicNode:
				return mopSeq
			default:
				panic("no such op")
			}
		case node.CallNode:
			return callSeq
		case node.ReturnNode:
			return returnSeq
		case node.ConditionalNode:
			return ifSeq
		case node.IfNode:
			return ifExpr
		case node.WhileNode:
			return whileSeq
		case node.RepeatNode:
			return repeatSeq
		case node.LoopNode:
			return loopSeq
		case node.ExitNode:
			return exitSeq
		case node.DerefNode:
			return derefSeq
		case node.InitNode:
			return frame.Tail(frame.STOP)
		case node.IndexNode:
			return indexSeq
		case node.FieldNode:
			return Propose(fieldSeq)
		case node.TrapNode:
			return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				switch code := next.Left().(type) {
				case node.ConstantNode:
					log.Println("TRAP:", traps.This(code.Data()))
					return frame.Tail(frame.WRONG), frame.NOW
				default:
					panic(fmt.Sprintln("unsupported code", reflect.TypeOf(code)))
				}
			}
		case node.WithNode:
			return withSeq
		case node.GuardNode:
			return guardSeq
		case node.CaseNode:
			return caseSeq
		case node.RangeNode:
			return rangeSeq
		case node.CompNode:
			return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				right := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					if next.Right() != nil {
						rt2.Push(rt2.New(next.Right()), f)
						return frame.Tail(frame.STOP), frame.LATER
					}
					return frame.End()
				}
				left := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
					if next.Left() != nil {
						rt2.Push(rt2.New(next.Left()), f)
						return right, frame.LATER
					}
					return right, frame.NOW
				}
				return left, frame.NOW
			}
		default:
			panic(fmt.Sprintln("unknown node", reflect.TypeOf(n), n.Adr()))
		}*/
	switch n.(type) {
	case node.Statement:
		return Propose(eval.BeginStatement)
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	panic(0)
}

func epilogue(n node.Node) frame.Sequence {
	/*	switch e := n.(type) {
		case node.AssignNode, node.InitNode, node.CallNode, node.ConditionalNode, node.WhileNode,
			node.RepeatNode, node.ExitNode, node.WithNode, node.CaseNode, node.CompNode:
			return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				next := n.Link()
				//fmt.Println("from", reflect.TypeOf(n))
				//fmt.Println("next", reflect.TypeOf(next))
				if next != nil {
					nf := rt2.New(next)
					if nf != nil {
						f.Root().PushFor(nf, f.Parent())
					}
				}
				if _, ok := n.(node.CallNode); ok {
					if f.Parent() != nil {
						par := rt2.RegOf(f.Parent())
						for k, v := range rt2.RegOf(f) {
							par[k] = v
						}
						val := rt2.ValueOf(f.Parent())
						for k, v := range rt2.ValueOf(f) {
							val[k] = v
						}
					}
				}
				return frame.End()
			}
		case node.EnterNode:
			return func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
				//fmt.Println(rt_module.DomainModule(f.Domain()).Name)
				if e.Enter() == enter.PROCEDURE {
					rt2.ThisScope(f).Target().(scope.ScopeAllocator).Dispose(n)
				}
				//возвращаем результаты вызова функции
				if f.Parent() != nil {
					par := rt2.RegOf(f.Parent())
					for k, v := range rt2.RegOf(f) {
						par[k] = v
					}
					val := rt2.ValueOf(f.Parent())
					for k, v := range rt2.ValueOf(f) {
						val[k] = v
					}
				}
				return frame.End()
			}
		case node.OperationNode, node.ReturnNode, node.IfNode, node.LoopNode,
			node.DerefNode, node.IndexNode, node.TrapNode, node.GuardNode, node.RangeNode, node.FieldNode:
			return nil
		default:
			fmt.Println(reflect.TypeOf(n))
			panic("unhandled epilogue")
		}*/
	switch n.(type) {
	case node.Statement:
		return Propose(eval.EndStatement)
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	panic(0)
}
