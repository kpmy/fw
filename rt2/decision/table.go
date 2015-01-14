package decision

import (
	"fw/cp/module"
	"fw/cp/node"
	"fw/rt2/context"
	"fw/rt2/frame"
)

var (
	Run         func(global context.Domain, init []*module.Module)
	PrologueFor func(n node.Node) frame.Sequence
	EpilogueFor func(n node.Node) frame.Sequence
)
