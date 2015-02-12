package wrap

import (
	"fw/cp/node"
	"fw/rt2"
	"fw/rt2/context"
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
		var key interface{}
		if fr.Parent() != nil {
			key = rt2.RegOf(fr.Parent())[context.KEY]
		}
		in := eval.IN{IR: rt2.NodeOf(fr), Frame: fr, Parent: fr.Parent(), Key: key}
		return This(a(in))
	}
}

func test(n node.Node) (bool, int) {
	switch n.(type) {
	default:
		return true, 0
	}
}

func prologue(n node.Node) frame.Sequence {
	switch n.(type) {
	case node.Statement:
		return Propose(eval.BeginStatement)
	case node.Expression:
		return Propose(eval.BeginExpression)
	case node.Designator:
		return Propose(eval.BeginDesignator)
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	panic(0)
}

func epilogue(n node.Node) frame.Sequence {
	switch n.(type) {
	case node.Statement:
		return Propose(eval.EndStatement)
	case node.ConstantNode, node.VariableNode, node.DyadicNode: //do nothing
	default:
		halt.As(100, reflect.TypeOf(n))
	}
	return nil
}
