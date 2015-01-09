package rt2

import (
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
	"fw/rt2/nodeframe"
	"fw/rt2/scope"
)

var utils nodeframe.FrameUtils

func DataOf(f frame.Frame) map[interface{}]interface{} { return utils.DataOf(f) }
func NodeOf(f frame.Frame) node.Node                   { return utils.NodeOf(f) }
func Push(f, p frame.Frame)                            { utils.Push(f, p) }
func New(n node.Node) frame.Frame                      { return utils.New(n) }
func ScopeOf(f frame.Frame) scope.Manager {
	return f.Domain().Discover(context.SCOPE).(scope.Manager)
}
