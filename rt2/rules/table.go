//dynamicaly loading from outer space
package rules

import (
	"cp/node"
	"fmt"
	"reflect"
	"rt2/context"
	"rt2/decision"
	"rt2/frame"
	"rt2/nodeframe"
	"rt2/scope"
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
	default:
		panic(fmt.Sprintln("unknown node", reflect.TypeOf(n)))
	}
}

func epilogue(n node.Node) frame.Sequence {
	var fu nodeframe.FrameUtils
	switch n.(type) {
	case node.AssignNode:
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
	case node.OperationNode:
		return nil //do nothing
	case node.CallNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			next := n.Link()
			if next != nil {
				f.Root().PushFor(fu.New(next), f.Parent())
			}
			return frame.End()
		}
	case node.ReturnNode:
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
