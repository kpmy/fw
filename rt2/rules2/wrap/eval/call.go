package eval

import (
	"encoding/json"
	"fmt"
	"fw/cp"
	"fw/cp/constant"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2"
	"fw/rt2/frame"
	rtm "fw/rt2/module"
	"fw/rt2/scope"
	"math"
	"reflect"
	"ypk/assert"
	"ypk/halt"
	"ypk/mathe"
)

var sys map[string]func(IN, node.Node) OUT

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
	cn := node.New(constant.CALL, cp.Some())
	ol := m.NodeByObject(obj)
	assert.For(len(ol) <= 1, 40)
	cn.SetLeft(ol[0])
	cc := node.New(constant.CONSTANT, cp.Some()).(node.ConstantNode)
	cc.SetData(data)
	cc.SetType(object.SHORTSTRING)
	cn.SetRight(cc)
	rt2.Push(rt2.New(cn), f)
}

func go_process(in IN, par node.Node) OUT {
	assert.For(par != nil, 20)
	sm := rt2.ThisScope(in.Frame)
	do := func(val string) {
		if val != "" {
			msg := &Msg{}
			if err := json.Unmarshal([]byte(val), msg); err == nil {
				switch msg.Type {
				case "log":
					fmt.Print(msg.Data)
					callHandler(in.Frame, scope.FindObjByName(sm, "go_handler"), `{"type":"log"}`)
				case "core":
					switch msg.Command {
					case "load":
						LoadMod(in.Frame, msg.Data)
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
	const left = "sys:left"
	return GetExpression(in, left, par, func(IN) OUT {
		val := rt2.ValueOf(in.Frame)[KeyOf(in, left)]
		assert.For(val != nil, 20)
		do(scope.GoTypeFrom(val).(string))
		return Later(Tail(STOP))
	})
}

func go_math(in IN, par node.Node) OUT {
	const (
		LN   = 1.0
		MANT = 2.0
		EXP  = 3.0
	)
	assert.For(par != nil, 20)
	sm := rt2.ThisScope(in.Frame)
	res := math.NaN()
	switch p := par.(type) {
	case node.VariableNode:
		val := sm.Select(p.Object().Adr())
		rv, ok := scope.GoTypeFrom(val).([]float64)
		assert.For(ok && (len(rv) > 1), 100, rv)
		switch rv[0] {
		case LN:
			res = math.Log(rv[1])
		case MANT:
			res, _ = mathe.Me(rv[1])
		case EXP:
			_, res = mathe.Me(rv[1])
		default:
			halt.As(100, rv[0])
		}
	default:
		halt.As(100, reflect.TypeOf(p))
	}
	id := cp.ID(cp.Some())
	rt2.RegOf(in.Parent)[in.Key] = id
	rt2.ValueOf(in.Parent)[id] = scope.TypeFromGo(res)
	return End()
}

func init() {
	sys = make(map[string]func(IN, node.Node) OUT)
	sys["go_process"] = go_process
	sys["go_math"] = go_math
}

var LoadMod func(frame.Frame, string)

func syscall(in IN) OUT {
	n := in.IR
	name := n.Left().Object().Name()
	return sys[name](in, n.Right())
}
