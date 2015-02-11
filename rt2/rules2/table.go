//dynamicaly loading from outer space
package rules2

import (
	"fw/cp/module"
	"fw/rt2"
	"fw/rt2/context"
	"fw/rt2/decision"
	"fw/rt2/frame"
	"fw/rt2/frame/std"
	_ "fw/rt2/rules2/wrap"
	"fw/utils"
	"log"
	"time"
	"ypk/assert"
)

type flow struct {
	root   frame.Stack
	parent frame.Frame
	domain context.Domain
	fl     []frame.Frame
	cl     []frame.Frame
	this   int
}

func (f *flow) Do() (ret frame.WAIT) {
	const Z frame.WAIT = -1
	x := Z
	if f.this >= 0 {
		x = f.fl[f.this].Do()
	}
	switch x {
	case frame.NOW, frame.WRONG, frame.LATER, frame.BEGIN:
		ret = x
	case frame.END:
		old := f.Root().(*std.RootFrame).Drop()
		assert.For(old != nil, 40)
		f.cl = append(f.cl, old)
		ret = frame.LATER
	case frame.STOP, Z:
		f.this--
		if f.this >= 0 {
			ret = frame.LATER
		} else {
			if len(f.cl) > 0 {
				for _, old := range f.cl {
					n := rt2.NodeOf(old)
					rt2.Push(rt2.New(n), old.Parent())
				}
				f.cl = nil
				ret = frame.LATER
			} else {
				ret = frame.STOP
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
		log.Println("total steps", i)
		log.Println("spent", t1.Sub(t0))
	}
}

func init() {
	decision.Run = run
}
