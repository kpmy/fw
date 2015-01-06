package module

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"ypk/assert"
)

type Module struct {
	Enter   node.Node
	Objects map[node.Node][]object.Object
	Nodes   []node.Node
	Types   map[node.Node][]object.ComplexType
}

type named interface {
	Name() string
}

func (m *Module) TypeByName(scope node.Node, name string) (ret object.ComplexType) {
	assert.For(name != "", 20)
	for _, typ := range m.Types[scope] {
		fmt.Print(typ)
		if v, ok := typ.(named); ok && v.Name() == name {
			ret = typ
			break //стыд какой
		}
	}
	return ret
}

func (m *Module) NodeByObject(obj object.Object) (ret node.Node) {
	assert.For(obj != nil, 20)
	for i := 0; (i < len(m.Nodes)) && (ret == nil); i++ {
		node := m.Nodes[i]
		if node.Object() == obj {
			ret = node
		}
	}
	return ret
}
