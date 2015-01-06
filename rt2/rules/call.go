package rules

import (
	"encoding/json"
	"fmt"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	mod "fw/rt2/module"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"ypk/assert"
)

/**
Для CallNode
	.Left() указывает на процедуру или на переменную процедурного типа
	.Left().Object() указывает на список внутренних объектов, в т.ч. переменных
	.Object() указывает первый элемент из списка входных параметров/переменных,
	то же что и.Left().Object().Link(), далее .Link() указывает на следующие за ним входные параметры
	.Right() указывает на узлы, которые передаются в параметры
*/

var sys map[string]func(f frame.Frame, par node.Node)

type Msg struct {
	Type string
	Data string
}

func callHandler(f frame.Frame, obj object.Object, data interface{}) {
	//n := rt2.Utils.NodeOf(f)
	fmt.Println("call handler", obj)
	if obj == nil {
		return
	}
	m := mod.DomainModule(f.Domain())
	cn := node.New(constant.CALL)
	cn.SetLeft(m.NodeByObject(obj))
	cc := node.New(constant.CONSTANT).(node.ConstantNode)
	cc.SetData(data)
	cn.SetRight(cc)
	rt2.Utils.Push(rt2.Utils.New(cn), f)
}

func process(f frame.Frame, par node.Node) {
	assert.For(par != nil, 20)
	sm := scope.This(f.Domain().Discover(context.SCOPE))
	switch par.(type) {
	case node.ConstantNode:
		msg := &Msg{}
		val := par.(node.ConstantNode).Data().(string)
		if err := json.Unmarshal([]byte(val), msg); err == nil {
			switch msg.Type {
			case "log":
				fmt.Println(msg.Data)
				callHandler(f, scope.FindObjByName(sm, "go_handler"), `{"type":"log"}`)
			default:
				panic(40)
			}
		}
	default:
		panic(fmt.Sprintln("unsupported param"))
	}
}

func init() {
	sys = make(map[string]func(f frame.Frame, par node.Node))
	sys["go_process"] = process
}

func syscall(f frame.Frame) {
	n := rt2.Utils.NodeOf(f)
	name := n.Left().Object().Name()
	sys[name](f, n.Right())
}

func callSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	var fu nodeframe.FrameUtils
	n := fu.NodeOf(f)

	call := func(proc node.Node) {
		nf := fu.New(proc)
		fu.Push(nf, f)
		//передаем ссылку на цепочку значений параметров в данные фрейма входа в процедуру
		if (n.Right() != nil) && (proc.Object() != nil) {
			fu.DataOf(nf)[proc.Object()] = n.Right()
		} else {
			//fmt.Println("no data for call")
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
		obj := sc.Select(scope.Designator(n.Left()))

		if obj, ok := obj.(object.Object); ok {
			proc := m.NodeByObject(obj)
			call(proc)
		} else {
			name := n.Left().Object().Name()
			switch {
			case name == "go_process":
				syscall(f)
				return frame.Tail(frame.STOP), frame.LATER
			default:
				panic(fmt.Sprintln("unknown sysproc variable", name))
			}
		}

	default:
		panic("unknown call left")
	}
	return seq, ret
}
