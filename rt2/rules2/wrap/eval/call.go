package eval

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
	"math"
	"reflect"
	"ypk/assert"
	"ypk/halt"
	"ypk/mathe"
)

var sys map[string]func(f frame.Frame, par node.Node) OUT

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

func go_process(f frame.Frame, par node.Node) OUT {
	assert.For(par != nil, 20)
	sm := rt2.ThisScope(f)
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
						panic(0)
						//fmt.Println("try to load", msg.Data)
						glob := f.Domain().Discover(context.UNIVERSE).(context.Domain)
						modList := glob.Discover(context.MOD).(rtm.List)
						//						fl := glob.Discover(context.MT).(*flow)
						ml := make([]*cpm.Module, 0)
						_, err := modList.Load(msg.Data, func(m *cpm.Module) {
							ml = append(ml, m)
						})
						for i := len(ml) - 1; i >= 0; i-- {
							//							fl.grow(glob, ml[i])
						}
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
		return Later(Tail(STOP))
	case node.VariableNode, node.ParameterNode:
		val = scope.GoTypeFrom(sm.Select(p.Object().Adr())).(string)
		do(val)
		return Later(Tail(STOP))
	case node.DerefNode:
		panic(0)
		rt2.Push(rt2.New(p), f)
		/*		return This(expectExpr(f, p, func(...IN) (out OUT) {
				v := rt2.ValueOf(f)[p.Adr()]
				assert.For(v != nil, 60)
				val = scope.GoTypeFrom(v).(string)
				do(val)
				out.do = Tail(STOP)
				out.next = LATER
				return out
			})) */
	default:
		halt.As(100, "unsupported param", reflect.TypeOf(p))
	}
	panic(0)
}

func go_math(f frame.Frame, par node.Node) OUT {
	const (
		LN   = 1.0
		MANT = 2.0
		EXP  = 3.0
	)
	assert.For(par != nil, 20)
	sm := rt2.ThisScope(f)
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
	rt2.RegOf(f.Parent())[context.RETURN] = scope.TypeFromGo(res)
	return End()
}

func init() {
	sys = make(map[string]func(f frame.Frame, par node.Node) OUT)
	sys["go_process"] = go_process
	sys["go_math"] = go_math
}

func syscall(f frame.Frame) OUT {
	n := rt2.NodeOf(f)
	name := n.Left().Object().Name()
	return sys[name](f, n.Right())
}
