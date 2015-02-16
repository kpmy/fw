package rt2

import (
	"fw/cp"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
	"ypk/assert"
)

var utils nodeframe.NodeFrameUtils

//func DataOf(f frame.Frame) map[interface{}]interface{} { panic(100) }
func RegOf(f frame.Frame) map[interface{}]interface{} { return utils.DataOf(f) }
func ValueOf(f frame.Frame) map[cp.ID]scope.Value     { return utils.ValueOf(f) }
func NodeOf(f frame.Frame) node.Node                  { return utils.NodeOf(f) }
func Push(f, p frame.Frame)                           { utils.Push(f, p) }
func New(n node.Node) frame.Frame                     { return utils.New(n) }

func ThisScope(f frame.Frame) scope.Manager {
	return f.Domain().Discover(context.SCOPE).(scope.Manager)
}

func ScopeFor(f frame.Frame, id cp.ID, fn ...scope.ValueOf) (ret scope.Manager) {
	ret = f.Domain().Global().Discover(context.SCOPE).(scope.Manager)
	assert.For(ret != nil, 60, id)
	if len(fn) == 1 {
		ret.Select(id, fn[0])
	}
	return
}

func ReplaceDomain(f frame.Frame, d context.Domain) { utils.ReplaceDomain(f, d) }
func Assert(f frame.Frame, ok frame.Assert)         { utils.Assert(f, ok) }
