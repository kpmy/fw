//dynamicaly loading from outer space
package rules

import (
	"fmt"
	"fw/cp/module"
	"fw/cp/node"
	"fw/cp/traps"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/rt2/frame/std"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"
	"time"
	"ypk/assert"
)

func prologue(n node.Node) frame.Sequence {
	//fmt.Println(reflect.TypeOf(n))
	switch next := n.(type) {
	case node.EnterNode:
		var tail frame.Sequence
		tail = func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			q := f.Root().Queue()
			if q != nil {
				f.Root().PushFor(q, nil)
				return tail, frame.NOW
			} else {
				return enterSeq, frame.NOW
			}
		}
		return tail
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
	case node.DerefNode:
		return derefSeq
	case node.InitNode:
		return frame.Tail(frame.STOP)
	case node.IndexNode:
		return indexSeq
	case node.FieldNode:
		return Propose(fieldSeq)
	case node.TrapNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			switch code := next.Left().(type) {
			case node.ConstantNode:
				utils.PrintTrap("TRAP:", traps.This(code.Data()))
				return frame.Tail(frame.WRONG), frame.NOW
			default:
				panic(fmt.Sprintln("unsupported code", reflect.TypeOf(code)))
			}
		}
	case node.WithNode:
		return withSeq
	case node.GuardNode:
		return guardSeq
	case node.CaseNode:
		return caseSeq
	case node.RangeNode:
		return rangeSeq
	case node.CompNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			right := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				if next.Right() != nil {
					rt2.Push(rt2.New(next.Right()), f)
					return frame.Tail(frame.STOP), frame.LATER
				}
				return frame.End()
			}
			left := func(f frame.Frame) (frame.Sequence, frame.WAIT) {
				if next.Left() != nil {
					rt2.Push(rt2.New(next.Left()), f)
					return right, frame.LATER
				}
				return right, frame.NOW
			}
			return left, frame.NOW
		}
	default:
		panic(fmt.Sprintln("unknown node", reflect.TypeOf(n)))
	}
}

func epilogue(n node.Node) frame.Sequence {
	switch n.(type) {
	case node.AssignNode, node.InitNode, node.CallNode, node.ConditionalNode, node.WhileNode,
		node.RepeatNode, node.ExitNode, node.WithNode, node.CaseNode, node.CompNode:
		return func(f frame.Frame) (frame.Sequence, frame.WAIT) {
			next := n.Link()
			//fmt.Println("from", reflect.TypeOf(n))
			//fmt.Println("next", reflect.TypeOf(next))
			if next != nil {
				f.Root().PushFor(rt2.New(next), f.Parent())
			}
			if _, ok := n.(node.CallNode); ok {
				if f.Parent() != nil {
					par := rt2.RegOf(f.Parent())
					for k, v := range rt2.RegOf(f) {
						par[k] = v
					}
					val := rt2.ValueOf(f.Parent())
					for k, v := range rt2.ValueOf(f) {
						val[k] = v
					}
				}
			}
			return frame.End()
		}
	case node.EnterNode:
		return func(f frame.Frame) (seq frame.Sequence, ret frame.WAIT) {
			//fmt.Println(rt_module.DomainModule(f.Domain()).Name)
			sm := f.Domain().Discover(context.SCOPE).(scope.Manager)
			sm.Target().(scope.ScopeAllocator).Dispose(n)
			//возвращаем результаты вызова функции
			if f.Parent() != nil {
				par := rt2.RegOf(f.Parent())
				for k, v := range rt2.RegOf(f) {
					par[k] = v
				}
				val := rt2.ValueOf(f.Parent())
				for k, v := range rt2.ValueOf(f) {
					val[k] = v
				}
			}
			return frame.End()
		}
	case node.OperationNode, node.ReturnNode, node.IfNode, node.LoopNode,
		node.DerefNode, node.IndexNode, node.TrapNode, node.GuardNode, node.RangeNode, node.FieldNode:
		return nil
	default:
		fmt.Println(reflect.TypeOf(n))
		panic("unhandled epilogue")
	}
}

type flow struct {
	root   frame.Stack
	parent frame.Frame
	domain context.Domain
	fl     []frame.Frame
	cl     []frame.Frame
	this   int
}

func (f *flow) Do() (ret frame.WAIT) {
	const Z WAIT = -1
	x := Z
	if f.this >= 0 {
		x = waiting(f.fl[f.this].Do())
	}
	switch x {
	case NOW, WRONG, LATER, BEGIN:
		ret = WAIT.wait(x)
	case END:
		old := f.Root().(*std.RootFrame).Drop()
		assert.For(old != nil, 40)
		f.cl = append(f.cl, old)
		ret = WAIT.wait(LATER)
	case STOP, Z:
		f.this--
		if f.this >= 0 {
			ret = WAIT.wait(LATER)
		} else {
			if len(f.cl) > 0 {
				for _, old := range f.cl {
					n := rt2.NodeOf(old)
					rt2.Push(rt2.New(n), old.Parent())
				}
				f.cl = nil
				ret = WAIT.wait(LATER)
			} else {
				ret = WAIT.wait(STOP)
			}
		}
	}
	utils.PrintFrame(">", ret)
	return ret
}

func (f *flow) OnPush(root frame.Stack, parent frame.Frame) {
	f.root = root
	f.parent = parent
	//fmt.Println("flow control pushed")
	f.this = len(f.fl) - 1
}

func (f *flow) OnPop() {
	//fmt.Println("flow control poped")
}

func (f *flow) Parent() frame.Frame    { return f.parent }
func (f *flow) Root() frame.Stack      { return f.root }
func (f *flow) Domain() context.Domain { return f.domain }
func (f *flow) Init(d context.Domain) {
	assert.For(f.domain == nil, 20)
	assert.For(d != nil, 21)
	f.domain = d
}

func (f *flow) Handle(msg interface{}) {
	assert.For(msg != nil, 20)
}

func (f *flow) grow(global context.Domain, m *module.Module) {
	utils.PrintScope("queue", m.Name)
	nf := rt2.New(m.Enter)
	f.root.PushFor(nf, nil)
	f.fl = append(f.fl, nf)
	global.Attach(m.Name, nf.Domain())
}

func run(global context.Domain, init []*module.Module) {
	{
		fl := &flow{root: std.NewRoot()}
		global.Attach(context.STACK, fl.root.(context.ContextAware))
		global.Attach(context.MT, fl)
		for i := len(init) - 1; i >= 0; i-- {
			fl.grow(global, init[i])
		}
		fl.root.PushFor(fl, nil)
		i := 0
		t0 := time.Now()
		for x := frame.NOW; x == frame.NOW; x = fl.root.(frame.Frame).Do() {
			utils.PrintFrame("STEP", i)
			//assert.For(i < 1000, 40)
			i++
		}
		t1 := time.Now()
		fmt.Println("total steps", i)
		fmt.Println("spent", t1.Sub(t0))
	}
}

func init() {
	decision.PrologueFor = prologue
	decision.EpilogueFor = epilogue
	decision.Run = run
}
