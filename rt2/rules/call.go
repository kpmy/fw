package rules

import (
	"cp/node"
	"rt2/context"
	"rt2/frame"
	mod "rt2/module"
	"rt2/nodeframe"
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
	switch n.Left().(type) {
	case node.ProcedureNode:
		m := mod.DomainModule(f.Domain())
		proc := m.NodeByObject(n.Left().Object())
		nf := fu.New(proc)
		f.Root().Push(nf)
		//передаем ссылку на цепочку значений параметров в данные фрейма входа в процедуру
		if (n.Right() != nil) && (proc.Object() != nil) {
			dd := make(map[interface{}]interface{})
			dd[proc.Object()] = n.Right()
			m := new(frame.SetDataMsg)
			m.Data = dd
			nf.(context.ContextAware).Handle(m)
		}
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			return frame.End()
		}
		ret = frame.SKIP
	default:
		panic("unknown call left")
	}
	return seq, ret
}
