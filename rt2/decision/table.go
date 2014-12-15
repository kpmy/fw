package decision

import (
	"cp/node"
	"rt2/frame"
)

var (
	PrologueFor func(n node.Node) frame.Sequence
	EpilogueFor func(n node.Node) frame.Sequence
)
