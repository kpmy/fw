package nodeframe

import (
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/utils"
	"reflect"
	"ypk/assert"
)

var count int64

type FrameUtils struct{}

func (fu FrameUtils) New(n node.Node) (f frame.Frame) {
	assert.For(n != nil, 20)
	f = new(nodeFrame)
	f.(*nodeFrame).ir = n
	f.(*nodeFrame).data = make(map[interface{}]interface{})
	//utils.Println("_", "NEW", reflect.TypeOf(n))
	return f
}

func (fu FrameUtils) Push(f, p frame.Frame) {
	assert.For(f != nil, 20)
	pp, _ := p.(*nodeFrame)
	pp.push(f)

}

func (fu FrameUtils) NodeOf(f frame.Frame) node.Node {
	ff, _ := f.(*nodeFrame)
	assert.For(ff.ir != nil, 40)
	return ff.ir
}

func (fu FrameUtils) DataOf(f frame.Frame) map[interface{}]interface{} {
	return f.(*nodeFrame).data
}

type nodeFrame struct {
	root   frame.Stack
	parent frame.Frame
	ir     node.Node
	seq    frame.Sequence
	domain context.Domain
	data   map[interface{}]interface{}
	num    int64
}

func (f *nodeFrame) Do() frame.WAIT {
	assert.For(f.seq != nil, 20)
	next, ret := f.seq(f)
	utils.Println(f.num, ret, reflect.TypeOf(f.ir))
	if next != nil {
		assert.For(ret != frame.STOP, 40)
		f.seq = next
	} else {
		assert.For(ret == frame.STOP, 41)
	}
	return ret
}

func (f *nodeFrame) onPush() {
	f.num = count
	count++
	utils.Println("_", "PUSH", reflect.TypeOf(f.ir))
	f.seq = decision.PrologueFor(f.ir)
}

func (f *nodeFrame) OnPop() {
	count--
	utils.Println("_", "POP", reflect.TypeOf(f.ir))
	f.seq = decision.EpilogueFor(f.ir)
	if f.seq != nil {
		_, _ = f.seq(f)
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
