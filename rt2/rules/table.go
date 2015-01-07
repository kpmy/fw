//dynamicaly loading from outer space
package rules

import (
	"fmt"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"
)

func prologue(n node.Node) frame.Sequence {
	//fmt.Println(reflect.TypeOf(n))
	switch next := n.(type) {
	case node.EnterNode:
		return enterSeq
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
	case node.TrapNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			switch code := next.Left().(type) {
			case node.ConstantNode:
				utils.PrintTrap("TRAP:", code.Data())
				return frame.Tail(frame.WRONG), frame.NOW
			default:
				panic(fmt.Sprintln("unsupported code", reflect.TypeOf(code)))
			}
		}
	case node.WithNode:
		return withSeq
	case node.GuardNode:
		return guardSeq
	default:
		panic(fmt.Sprintln("unknown node", reflect.TypeOf(n)))
	}
}

func epilogue(n node.Node) frame.Sequence {
	var fu nodeframe.FrameUtils
	switch n.(type) {
	case node.AssignNode, node.InitNode, node.CallNode, node.ConditionalNode, node.WhileNode,
		node.RepeatNode, node.ExitNode, node.WithNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			next := n.Link()
			//fmt.Println("from", reflect.TypeOf(n))
			//fmt.Println("next", reflect.TypeOf(next))
			if next != nil {
				f.Root().PushFor(fu.New(next), f.Parent())
			}
			return frame.End()
		}
	case node.EnterNode:
		return func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			sm := scope.This(f.Domain().Discover(context.SCOPE))
			sm.Dispose(n)
			return frame.End()
		}
	case node.OperationNode, node.ReturnNode, node.IfNode, node.LoopNode,
		node.DerefNode, node.IndexNode, node.TrapNode, node.GuardNode:
		return nil
	default:
		fmt.Println(reflect.TypeOf(n))
		panic("unhandled epilogue")
	}
}

func init() {
	decision.PrologueFor = prologue
	decision.EpilogueFor = epilogue
}
