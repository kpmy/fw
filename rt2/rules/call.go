package rules

import (
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
	mod "fw/rt2/module"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
)

/**
Для CallNode
	.Left() указывает на процедуру
	.Left().Object() указывает на список внутренних объектов, в т.ч. переменных
	.Object() указывает первый элемент из списка входных параметров/переменных,
	то же что и.Left().Object().Link(), далее .Link() указывает на следующие за ним входные параметры
	.Right() указывает на узлы, которые передаются в параметры
*/
func callSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)

	call := func(proc node.Node) {
		nf := fu.New(proc)
		fu.Push(nf, f)
		//передаем ссылку на цепочку значений параметров в данные фрейма входа в процедуру
		if (n.Right() != nil) && (proc.Object() != nil) {
			fu.DataOf(nf)[proc.Object()] = n.Right()
		}
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			var fu nodeframe.FrameUtils
			fu.DataOf(f.Parent())[n] = fu.DataOf(f)[n.Left().Object()]
			return frame.End()
		}
		ret = frame.LATER
	}

	switch n.Left().(type) {
	case node.ProcedureNode:
		m := mod.DomainModule(f.Domain())
		proc := m.NodeByObject(n.Left().Object())
		call(proc)
	case node.VariableNode:
		m := mod.DomainModule(f.Domain())
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		obj := sc.Select(n.Left().Object()).(object.Object)
		proc := m.NodeByObject(obj)
		call(proc)
	default:
		panic("unknown call left")
	}
	return seq, ret
}
