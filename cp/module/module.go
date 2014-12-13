package module

import (
	"cp/node"
	"cp/object"
)

type Module struct {
	Enter   node.Node
	Objects map[node.Node][]object.Object
	Nodes   []node.Node
}

func (m *Module) NodeByObject(obj object.Object) (ret node.Node) {
	if obj == nil {
		panic("obj must not be nil")
	}
	for i := 0; (i < len(m.Nodes)) && (ret == nil); i++ {
		node := m.Nodes[i]
		if node.Object() == obj {
			ret = node
		}
	}
	return ret
}
