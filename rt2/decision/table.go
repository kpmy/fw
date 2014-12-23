package decision

import (
	"fw/cp/node"
	"fw/rt2/frame"
)

var (
	PrologueFor func(n node.Node) frame.Sequence
	EpilogueFor func(n node.Node) frame.Sequence
)
