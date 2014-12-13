package module

import (
	"cp/node"
	"cp/object"
)

type Module struct {
	Enter   node.Node
	Objects []object.Object
}
