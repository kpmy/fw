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
	"reflect"
)

func prologue(n node.Node) frame.Sequence {
	//fmt.Println(reflect.TypeOf(n))
	switch n.(type) {
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
	default:
		panic(fmt.Sprintln("unknown node", reflect.TypeOf(n)))
	}
}

func epilogue(n node.Node) frame.Sequence {
	var fu nodeframe.FrameUtils
	switch n.(type) {
	case node.AssignNode, node.CallNode, node.ConditionalNode, node.WhileNode, node.RepeatNode, node.ExitNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			next := n.Link()
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
	case node.OperationNode, node.ReturnNode, node.IfNode, node.LoopNode:
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
