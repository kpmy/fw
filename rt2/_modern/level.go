package data

import (
	"fw/cp"
	"fw/cp/node"
	"fw/rt2/scope"
)

type level struct {
	root   node.Node
	k      map[cp.ID]int
	v      map[int]scope.Variable
	r      map[int]scope.Ref
	l      map[int]*level
	next   int
	ready  bool
	nested bool
}
