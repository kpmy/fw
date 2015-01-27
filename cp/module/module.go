package module

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"ypk/assert"
)

type Import struct {
	Name    string
	Objects []object.Object
}

type Module struct {
	Name    string
	Enter   node.Node
	Objects map[node.Node][]object.Object
	Nodes   []node.Node
	Types   map[node.Node][]object.ComplexType
	Imports map[string]Import
}

type named interface {
	Name() string
}

func (m *Module) ObjectByName(scope node.Node, name string) (rl []object.Object) {
	assert.For(name != "", 20)
	find := func(v []object.Object) (ret []object.Object) {
		for _, o := range v {
			if o.Name() == name {
				ret = append(ret, o)
			}
			//fmt.Println(o.Name(), name, o.Name() == name)
		}
		return ret
	}
	if scope == nil {
		for _, v := range m.Objects {
			rl = append(rl, find(v)...)
		}
	} else {
		rl = find(m.Objects[scope])
	}

	return rl
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

func (m *Module) ImportOf(obj object.Object) string {
	contains := func(v []object.Object) bool {
		for _, o := range v {
			if o == obj {
				return true
			}
		}
		return false
	}
	for _, v := range m.Imports {
		if contains(v.Objects) {
			return v.Name
		}
	}
	return ""
}

func (m *Module) NodeByObject(obj object.Object) (ret []node.Node) {
	assert.For(obj != nil, 20)
	for i := 0; (i < len(m.Nodes)) && (ret == nil); i++ {
		node := m.Nodes[i]
		if node.Object() == obj {
			ret = append(ret, node)
		}
	}
	return ret
}
