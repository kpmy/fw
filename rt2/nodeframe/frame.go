package nodeframe

import (
	"fmt"
	"fw/cp"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/rt2/scope"
	"fw/utils"
	"reflect"
	"ypk/assert"
)

var count int64

type NodeFrameUtils struct{}

func (fu NodeFrameUtils) New(n node.Node) (f frame.Frame) {
	assert.For(n != nil, 20)
	ok, _ := decision.AssertFor(n)
	//assert.For(ok, 21, "wrong node", reflect.TypeOf(n), code)
	if ok {
		f = new(nodeFrame)
		f.(*nodeFrame).ir = n
		f.(*nodeFrame).data = make(map[interface{}]interface{})
		f.(*nodeFrame).value = make(map[cp.ID]scope.Value)
		utils.PrintFrame("_", "NEW", reflect.TypeOf(n), n.Adr())
	}
	return f
}

func (fu NodeFrameUtils) Push(f, p frame.Frame) {
	assert.For(f != nil, 20)
	pp, _ := p.(*nodeFrame)
	pp.push(f)
}

func (fu NodeFrameUtils) NodeOf(f frame.Frame) node.Node {
	ff, ok := f.(*nodeFrame)
	if ok {
		assert.For(ff.ir != nil, 40)
		return ff.ir
	} else {
		return nil
	}
}

func (fu NodeFrameUtils) DataOf(f frame.Frame) map[interface{}]interface{} {
	return f.(*nodeFrame).data
}

func (fu NodeFrameUtils) ValueOf(f frame.Frame) map[cp.ID]scope.Value {
	return f.(*nodeFrame).value
}

func (fu NodeFrameUtils) ReplaceDomain(f frame.Frame, d context.Domain) {
	ff := f.(*nodeFrame)
	ff.domain = d
}

func (fu NodeFrameUtils) Assert(f frame.Frame, ok frame.Assert) {
	ff := f.(*nodeFrame)
	ff.assertion = ok
}

type nodeFrame struct {
	root      frame.Stack
	parent    frame.Frame
	ir        node.Node
	seq       frame.Sequence
	assertion frame.Assert
	domain    context.Domain
	data      map[interface{}]interface{}
	value     map[cp.ID]scope.Value
	num       int64
}

func done(f frame.Frame) {
	utils.PrintScope("== SCOPE ==")
	utils.PrintScope(f.Domain().Discover(context.SCOPE))
	utils.PrintScope("==  HEAP ==")
	utils.PrintScope(f.Domain().Discover(context.HEAP))
	utils.PrintScope("== CALL ==")
	utils.PrintScope(f.Domain().Discover(context.CALL))
	utils.PrintScope("^^^^^^^^^")
}

func (f *nodeFrame) Do() frame.WAIT {
	assert.For(f.seq != nil, 20)
	next, ret := f.seq(f)
	utils.PrintFrame(f.num, ret, reflect.TypeOf(f.ir), f.ir)
	utils.PrintFrame("data:", f.data, f.value)
	if next != nil {
		assert.For(ret != frame.STOP, 40)
		f.seq = next
	} else {
		assert.For(ret == frame.STOP || ret == frame.WRONG, 41)
		if ret == frame.WRONG {
			var fu NodeFrameUtils
			utils.PrintTrap()
			utils.PrintTrap("something WRONG")
			utils.PrintTrap(f.Domain().Discover(context.SCOPE))
			depth := 10
			f.root.ForEach(func(old frame.Frame) bool {
				n := fu.NodeOf(old)
				if n != nil {
					utils.PrintTrap(reflect.TypeOf(n), n.Adr())
					utils.PrintTrap("data: ", fu.DataOf(old), "value: ", fu.ValueOf(old))
				} else {
					utils.PrintTrap(reflect.TypeOf(old))
				}
				depth--
				return depth != 0
			})
			utils.PrintTrap(f.Domain().Discover(context.HEAP))
		}

	}
	defer done(f)
	return ret
}

func (f *nodeFrame) onPush() {
	f.num = count
	count++
	//	assert.For(count < 15, 40)
	utils.PrintFrame("_", "PUSH", reflect.TypeOf(f.ir))
	f.seq = decision.PrologueFor(f.ir)
}

func (f *nodeFrame) OnPop() {
	count--
	utils.PrintFrame("_", "POP", reflect.TypeOf(f.ir))
	f.seq = decision.EpilogueFor(f.ir)
	if f.seq != nil {
		_, end := f.seq(f)
		assert.For(end == frame.END || end == frame.STOP, 60, end)
	}
	if f.parent != nil {
		ff, ok := f.parent.(*nodeFrame)
		if ok && ff.assertion != nil {
			check := func(ok bool, msg ...interface{}) {
				if !ok {
					panic(fmt.Sprintln("assert for", reflect.TypeOf(ff.ir), ff.ir, fmt.Sprint(msg...)))
				} else {
					utils.PrintFrame("assert passed for", reflect.TypeOf(ff.ir), ff.ir)
				}
			}
			ff.assertion(ff, check)

			ff.assertion = nil
		}
	}
}

func (f *nodeFrame) push(n frame.Frame) {
	f.root.PushFor(n, f)
}

func (f *nodeFrame) OnPush(root frame.Stack, parent frame.Frame) {
	f.root = root
	f.parent = parent
	f.onPush()
}

func (f *nodeFrame) Parent() frame.Frame    { return f.parent }
func (f *nodeFrame) Root() frame.Stack      { return f.root }
func (f *nodeFrame) Domain() context.Domain { return f.domain }
func (f *nodeFrame) Init(d context.Domain) {
	assert.For(f.domain == nil, 20)
	assert.For(d != nil, 21)
	f.domain = d
}

func (f *nodeFrame) Handle(msg interface{}) {
	assert.For(msg != nil, 20)

}
