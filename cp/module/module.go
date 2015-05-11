package module

import (
	"fw/cp"
	"fw/cp/constant/enter"
	"fw/cp/node"
	"fw/cp/object"
	"github.com/kpmy/ypk/assert"
	"github.com/kpmy/ypk/halt"
	"reflect"
	"strconv"
)

type Import struct {
	Name    string
	Objects []object.Object
	Types   []object.ComplexType
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
		//fmt.Print(typ)
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

func (m *Module) ImportOfType(obj object.ComplexType) string {
	contains := func(v []object.ComplexType) bool {
		for _, o := range v {
			if o.Adr() == obj.Adr() {
				return true
			}
		}
		return false
	}
	for _, v := range m.Imports {
		if contains(v.Types) {
			return v.Name
		}
	}
	return ""
}

func (m *Module) NodeByObject(obj object.Object) (ret []node.Node) {
	assert.For(obj != nil, 20)
	for i := 0; (i < len(m.Nodes)) && (ret == nil); i++ {
		node := m.Nodes[i]
		if node.Object() != nil {
			if node.Object().Adr(0, 0) == obj.Adr() {
				ret = append(ret, node)
			}
		}
	}
	return ret
}

func (m *Module) Init(inittd ...func(t object.ComplexType)) {
	typeName := func(id cp.ID) string {
		for _, s := range m.Objects {
			for _, o := range s {
				switch t := o.(type) {
				case object.TypeObject:
					if o.Complex().Adr() == id {
						return t.Name()
					}
				}
			}
		}
		for _, i := range m.Imports {
			for _, o := range i.Objects {
				switch t := o.(type) {
				case object.TypeObject:
					if o.Complex().Adr() == id {
						return ">" + i.Name + "." + t.Name()
					}
				}
			}
		}
		return strconv.Itoa(int(id))
	}

	//fmt.Println("init", m.Name)
	for k, s := range m.Types {
		q := ""
		switch e := k.(type) {
		case node.EnterNode:
			switch e.Enter() {
			case enter.MODULE:
				q = m.Name
			case enter.PROCEDURE:
				q = m.Name + "..." + e.Object().Name()
			}
		default:
			halt.As(100, reflect.TypeOf(e))
		}
		for _, t := range s {
			t.Qualident(q + "." + typeName(t.Adr()))
		}
	}
	for _, s := range m.Imports {
		q := s.Name
		for _, t := range s.Types {
			t.Qualident(q + "." + typeName(t.Adr()))
		}
		for _, o := range s.Objects {
			o.Imp(s.Name)
		}
	}
	if len(inittd) > 0 {
		for _, s := range m.Types {
			for _, t := range s {
				inittd[0](t)
				switch rec := t.(type) {
				case object.RecordType:
					for x := rec.Link(); x != nil; x = x.Link() {
						switch f := x.(type) {
						case object.FieldObject:
							f.TypeOf(rec)
						}
					}
				}
			}
		}
	}
}
