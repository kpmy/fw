package rules

import (
	"encoding/json"
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	cpm "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

/**
Для CallNode
	.Left() указывает на процедуру или на переменную процедурного типа
	.Left().Object() указывает на список внутренних объектов, в т.ч. переменных
	.Object() указывает первый элемент из списка входных параметров/переменных,
	то же что и.Left().Object().Link(), далее .Link() указывает на следующие за ним входные параметры
	.Right() указывает на узлы, которые передаются в параметры
*/

var sys map[string]func(f frame.Frame, par node.Node) (frame.Sequence, frame.WAIT)

type Msg struct {
	Type    string
	Command string
	Data    string
}

func callHandler(f frame.Frame, obj object.Object, data interface{}) {
	//n := rt2.Utils.NodeOf(f)
	//fmt.Println("call handler", obj)
	if obj == nil {
		return
	}
	m := rtm.DomainModule(f.Domain())
	cn := node.New(constant.CALL, int(cp.SomeAdr()))
	ol := m.NodeByObject(obj)
	assert.For(len(ol) <= 1, 40)
	cn.SetLeft(ol[0])
	cc := node.New(constant.CONSTANT, int(cp.SomeAdr())).(node.ConstantNode)
	cc.SetData(data)
	cc.SetType(object.SHORTSTRING)
	cn.SetRight(cc)
	rt2.Push(rt2.New(cn), f)
}

func process(f frame.Frame, par node.Node) (frame.Sequence, frame.WAIT) {
	assert.For(par != nil, 20)
	sm := f.Domain().Discover(context.SCOPE).(scope.Manager)
	do := func(val string) {
		if val != "" {
			msg := &Msg{}
			if err := json.Unmarshal([]byte(val), msg); err == nil {
				switch msg.Type {
				case "log":
					fmt.Print(msg.Data)
					callHandler(f, scope.FindObjByName(sm, "go_handler"), `{"type":"log"}`)
				case "core":
					switch msg.Command {
					case "load":
						//fmt.Println("try to load", msg.Data)
						glob := f.Domain().Discover(context.UNIVERSE).(context.Domain)
						modList := glob.Discover(context.MOD).(rtm.List)
						fl := glob.Discover(context.MT).(*flow)
						_, err := modList.Load(msg.Data, func(m *cpm.Module) {
							fl.grow(glob, m)
						})
						assert.For(err == nil, 60)
					default:
						halt.As(100, msg.Command)
					}
				default:
					panic(40)
				}
			} else {
				fmt.Println(val, "not a json")
			}
		}
	}
	var val string
	switch p := par.(type) {
	case node.ConstantNode:
		val = par.(node.ConstantNode).Data().(string)
		do(val)
		return frame.Tail(frame.STOP), frame.LATER
	case node.VariableNode, node.ParameterNode:
		val = scope.GoTypeFrom(sm.Select(p.Object().Adr())).(string)
		do(val)
		return frame.Tail(frame.STOP), frame.LATER
	case node.DerefNode:
		rt2.Push(rt2.New(p), f)
		return This(expectExpr(f, p, func(...IN) (out OUT) {
			v := rt2.ValueOf(f)[p.Adr()]
			assert.For(v != nil, 60)
			val = scope.GoTypeFrom(v).(string)
			do(val)
			out.do = Tail(STOP)
			out.next = LATER
			return out
		}))
	default:
		halt.As(100, "unsupported param", reflect.TypeOf(p))
	}
	panic(0)
}

func init() {
	sys = make(map[string]func(f frame.Frame, par node.Node) (frame.Sequence, frame.WAIT))
	sys["go_process"] = process
}

func syscall(f frame.Frame) (frame.Sequence, frame.WAIT) {
	n := rt2.NodeOf(f)
	name := n.Left().Object().Name()
	return sys[name](f, n.Right())
}

func callSeq(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
	n := rt2.NodeOf(f)

	call := func(proc node.Node, d context.Domain) {
		nf := rt2.New(proc)
		rt2.Push(nf, f)
		if d != nil {
			rt2.ReplaceDomain(nf, d)
		}
		//передаем ссылку на цепочку значений параметров в данные фрейма входа в процедуру
		if (n.Right() != nil) && (proc.Object() != nil) {
			rt2.RegOf(nf)[proc.Object()] = n.Right()
		} else {
			//fmt.Println("no data for call")
		}
		seq = func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			//			rt2.DataOf(f.Parent())[n] = rt2.DataOf(f)[n.Left().Object()]
			if f.Parent() != nil {
				rt2.ValueOf(f.Parent())[n.Adr()] = rt2.ValueOf(f)[n.Left().Object().Adr()]
			}
			return frame.End()
		}
		ret = frame.LATER
	}

	switch p := n.Left().(type) {
	case node.EnterNode:
		call(p, nil)
	case node.ProcedureNode:
		m := rtm.DomainModule(f.Domain())
		ml := f.Domain().Discover(context.UNIVERSE).(context.Domain).Discover(context.MOD).(rtm.List)
		if p.Super() {
			fmt.Println("supercall, stop for now")
			seq = Propose(Tail(STOP))
			ret = frame.NOW
		} else {
			if imp := m.ImportOf(n.Left().Object()); imp == "" || imp == m.Name {
				switch p.Object().Mode() {
				case object.LOCAL_PROC, object.EXTERNAL_PROC:
					proc := m.NodeByObject(n.Left().Object())
					assert.For(proc != nil, 40)
					call(proc[0], nil)
				case object.TYPE_PROC:
					var proc []node.Node
					//sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
					return This(expectExpr(f, n.Right(), func(...IN) (out OUT) {
						v := rt2.ValueOf(f)[n.Right().Adr()]
						var dm context.Domain
						var fn object.ProcedureObject
						_, c := scope.Ops.TypeOf(v)
						mod := rtm.ModuleOfType(f.Domain(), c)
						dm = f.Domain().Discover(context.UNIVERSE).(context.Domain).Discover(mod.Name).(context.Domain)
						ol := mod.Objects[mod.Enter]
						for _, _po := range ol {
							switch po := _po.(type) {
							case object.ProcedureObject:
								if po.Name() == p.Object().Name() && po.Link() != nil && po.Link().Complex() == c {
									fn = po
								}

							}
						}
						assert.For(fn != nil, 40)
						proc = mod.NodeByObject(fn)
						assert.For(proc != nil, 40)
						call(proc[0], dm)
						out.do = Tail(STOP)
						out.next = LATER
						return out
					}))

				default:
					halt.As(100, "wrong proc mode ", p.Object().Mode())
				}

				//fmt.Println(len(proc), len(n.Left().Object().Ref()))
				//fmt.Println("proc refs", proc)

			} else {
				m := ml.Loaded(imp)
				pl := m.ObjectByName(m.Enter, n.Left().Object().Name())
				var proc object.ProcedureObject
				var nl []node.Node
				for _, n := range pl {
					if n.Mode() == p.Object().Mode() {
						proc = n.(object.ProcedureObject)
					}
				}
				//utils.PrintFrame("proc refs", len(proc))
				switch proc.Mode() {
				case object.LOCAL_PROC, object.EXTERNAL_PROC:
					nl = m.NodeByObject(proc)
					utils.PrintFrame("foreign call", len(nl), "proc refs", proc)
					call(nl[0], f.Domain().Discover(context.UNIVERSE).(context.Domain).Discover(imp).(context.Domain))
				default:
					halt.As(100, "wrong proc mode ", p.Object().Mode())
				}
			}
		}
	case node.VariableNode:
		m := rtm.DomainModule(f.Domain())
		sc := f.Domain().Discover(context.SCOPE).(scope.Manager)
		obj := scope.GoTypeFrom(sc.Select(n.Left().Object().Adr()))

		if obj, ok := obj.(object.Object); ok {
			proc := m.NodeByObject(obj)
			call(proc[0], nil)
		} else {
			name := n.Left().Object().Name()
			switch {
			case name == "go_process":
				return syscall(f)
			default:
				panic(fmt.Sprintln("unknown sysproc variable", name))
			}
		}

	default:
		halt.As(100, "unknown call left: ", reflect.TypeOf(p))
	}
	return seq, ret
}
