//dynamicaly loading from outer space
package rules

import (
	"cp/node"
	"fmt"
	"reflect"
	"rt2/decision"
	"rt2/frame"
	"rt2/nodeframe"
	"ypk/assert"
)

func prologue(n node.Node) frame.Sequence {
	var fu nodeframe.FrameUtils
	fmt.Println(reflect.TypeOf(n))
	switch n.(type) {
	case node.EnterNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			node := fu.NodeOf(f).Right()
			assert.For(node != nil, 40)
			f.Root().Push(fu.New(node))
			return frame.Tail(frame.STOP), frame.SKIP
		}
	case node.AssignNode:
		return assignSeq
	case node.OperationNode:
		return opSeq
	case node.CallNode:
		return callSeq
	default:
		panic("unknown node")
	}
}

func epilogue(n node.Node) frame.Sequence {
	var fu nodeframe.FrameUtils
	switch n.(type) {
	case node.AssignNode, node.CallNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			next := n.Link()
			if next != nil {
				f.Root().Push(fu.New(next))
			}
			return frame.End()
		}
	case node.OperationNode, node.EnterNode:
		return nil //do nothing
	default:
		fmt.Println(reflect.TypeOf(n))
		panic("unhandled epilogue")
	}
}

func init() {
	decision.PrologueFor = prologue
	decision.EpilogueFor = epilogue
}
