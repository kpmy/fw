package nodeframe

import (
	"cp/node"
	"rt2/context"
	"rt2/decision"
	"rt2/frame"
	"ypk/assert"
)

type FrameUtils struct{}

func (fu FrameUtils) New(n node.Node) (f frame.Frame) {
	assert.For(n != nil, 20)
	f = new(nodeFrame)
	f.(*nodeFrame).ir = n
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

type nodeFrame struct {
	root   frame.Stack
	parent frame.Frame
	ir     node.Node
	seq    frame.Sequence
	domain context.Domain
}

func (f *nodeFrame) Do() frame.WAIT {
	assert.For(f.seq != nil, 20)
	next, ret := f.seq(f)
	if next != nil {
		assert.For(ret != frame.STOP, 40)
		f.seq = next
	} else {
		assert.For(ret == frame.STOP, 41)
	}
	return ret
}

func (f *nodeFrame) onPush() {
	f.seq = decision.PrologueFor(f.ir)
}

func (f *nodeFrame) OnPop() {
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
