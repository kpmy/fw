package module

import (
	"fmt"
	"fw/cp"
	mod "fw/cp/module"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/xev"
	"os"
	"reflect"
	"ypk/assert"
	"ypk/halt"
)

type Loader func(*mod.Module)

type List interface {
	context.ContextAware
	AsList() []*mod.Module
	Load(string, ...Loader) (*mod.Module, error)
	Loaded(string) *mod.Module
	NewTypeCalc() TypeCalc
}

func New() List {
	return new(list).init()
}

type list struct {
	inner map[string]*mod.Module
	d     context.Domain
}

func (l *list) init() *list {
	l.inner = make(map[string]*mod.Module)
	return l
}

func (l *list) AsList() (ret []*mod.Module) {
	if len(l.inner) > 0 {
		ret = make([]*mod.Module, 0)
	}
	for _, v := range l.inner {
		ret = append(ret, v)
	}
	return ret
}

func (l *list) Domain() context.Domain {
	return l.d
}

func (l *list) Init(d context.Domain) {
	l.d = d
}

func (l *list) Handle(msg interface{}) {}

func (l *list) Load(name string, ldr ...Loader) (ret *mod.Module, err error) {
	assert.For(name != "", 20)
	fmt.Println("loading", name, "loaded", l.Loaded(name) != nil)
	ret = l.Loaded(name)
	var loader Loader = func(m *mod.Module) {}
	if len(ldr) > 0 {
		loader = ldr[0]
	}
	if ret == nil {
		path, _ := os.Getwd()
		ret = xev.Load(path, name+".oz")
		ret.Name = name
		for _, imp := range ret.Imports {
			fmt.Println(name, "imports", imp.Name, "loaded", l.Loaded(imp.Name) != nil)
			_, err = l.Load(imp.Name, loader)
		}
		if err == nil {
			ret.Init(func(t object.ComplexType) {
				fmt.Println(t.Qualident())
			})
			l.inner[name] = ret
			loader(ret)
			fmt.Println("loaded", name)
		}
	}
	return ret, err
}

func (l *list) Loaded(name string) *mod.Module {
	assert.For(name != "", 20)
	return l.inner[name]
}

func (l *list) NewTypeCalc() TypeCalc {
	return &tc{ml: l}
}

func ModuleDomain(d context.Domain, name string) context.Domain {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	return uni.Discover(name).(context.Domain)
}

func Module(d context.Domain, name string) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	return uni.Discover(context.MOD).(List).Loaded(name)
}

func DomainModule(d context.Domain) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	name := uni.Id(d)
	assert.For(name != "", 40)
	ml := uni.Discover(context.MOD).(List)
	return ml.Loaded(name)
}

func ModuleOfNode(d context.Domain, x node.Node) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, n := range m.Nodes {
			if n == x {
				return m
			}
		}
	}
	return nil
}

func ModuleOfObject(d context.Domain, x object.Object) *mod.Module {
	assert.For(x != nil, 20)
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, v := range m.Objects {
			for _, o := range v {
				if o == x {
					return m
				}
			}
		}
		for _, i := range m.Imports {
			for _, o := range i.Objects {
				if o == x {
					return m
				}
			}
		}
	}
	return nil
}

func ModuleOfType(d context.Domain, x object.ComplexType) *mod.Module {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	for _, m := range ml.AsList() {
		for _, v := range m.Types {
			for _, o := range v {
				if o.Adr() == x.Adr() { //сравнение по конкретному типу
					return m
				}
			}
		}
	}
	return nil
}

func MapImportType(d context.Domain, imp string, t object.ComplexType) object.ComplexType {
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	m := ml.Loaded(imp)
	for _, v := range m.Types[m.Enter] {
		if v.Equals(t) {
			return v
		}
	}
	return nil
}

func MapImportObject(d context.Domain, t object.Object) object.Object {
	if t.Imp() == "" {
		return t
	}
	imp := t.Imp()
	uni := d.Discover(context.UNIVERSE).(context.Domain)
	ml := uni.Discover(context.MOD).(List)
	m := ml.Loaded(imp)
	for _, v := range m.Objects[m.Enter] {
		if v.Name() == t.Name() {
			if v.Type() == t.Type() {
				return v
			}
		}
	}
	panic(0)
}

type TypeCalc interface {
	ConnectTo(interface{})
	MethodList() map[int][]Method
	ForeignBase() (*mod.Module, object.ComplexType)
}

type Method struct {
	Enter node.EnterNode
	Obj   object.Object
	Mod   *mod.Module
}

type tc struct {
	ml  List
	m   *mod.Module
	typ object.ComplexType
	TypeCalc
}

type Inherited interface {
	Complex(...object.ComplexType) object.ComplexType
}

func (c *tc) ConnectTo(x interface{}) {
	switch t := x.(type) {
	case object.ComplexType:
		c.typ = t
	case object.TypeObject:
		c.typ = t.Complex()
	default:
		halt.As(100, reflect.TypeOf(t))
	}
	c.m = ModuleOfType(c.ml.Domain(), c.typ)
	assert.For(c.m != nil, 60, c.typ.Qualident(), c.typ)
}

func (c *tc) MethodList() (ret map[int][]Method) {
	type mid struct {
		id   cp.ID
		name string
	}
	ret = make(map[int][]Method, 0)
	tmp := make(map[mid]object.Object, 0)
	depth := -1
	var deep func(*mod.Module, object.ComplexType)
	list := func(m *mod.Module, t object.ComplexType) {
		ol := m.Objects[m.Enter]
		for _, _po := range ol {
			switch po := _po.(type) {
			case object.ProcedureObject:
				var et node.EnterNode
				proc := m.NodeByObject(po)
				local := false
				for i := range proc {
					if e, ok := proc[i].(node.EnterNode); ok {
						local = true
						et = e
					}
				}
				if local && po.Link() != nil {
					for pt := po.Link().Complex(); pt != nil; {
						if t.Equals(pt) && tmp[mid{id: po.Adr(), name: po.Name()}] == nil {
							//fmt.Println("method", m.Name, pt.Qualident(), po.Name(), po.Adr(), local)
							tmp[mid{id: po.Adr(), name: po.Name()}] = po
							ret[depth] = append(ret[depth], Method{Enter: et, Obj: po, Mod: m})
							break
						} else if t.Equals(pt) {
							//fmt.Println("double ", pt.Qualident(), po.Name(), po.Adr())
						}
						if _, ok := pt.(Inherited); ok {
							pt = pt.(Inherited).Complex()
						} else {
							pt = nil
						}

					}
				}

			}
		}
	}
	foreign := func(t object.ComplexType) {
		for _, n := range c.m.Imports {
			for _, _it := range n.Objects {
				switch it := _it.(type) {
				case object.TypeObject:
					if it.Complex().Adr() == t.Adr() {
						nm := c.ml.Loaded(n.Name)
						nt := nm.TypeByName(nm.Enter, it.Name())
						deep(nm, nt)
					}
				}
			}
		}
	}
	deep = func(m *mod.Module, x object.ComplexType) {
		depth++
		tmp = make(map[mid]object.Object, 0)
		for t := x; t != nil; {
			list(m, t)
			z := t.(Inherited).Complex()
			if z != nil {
				t = z
			} else {
				foreign(t)
				t = nil
			}
		}
	}
	deep(c.m, c.typ)
	return
}

func (c *tc) ForeignBase() (*mod.Module, object.ComplexType) {
	for _, n := range c.m.Imports {
		for _, _it := range n.Objects {
			switch it := _it.(type) {
			case object.TypeObject:
				if it.Complex().Adr() == c.typ.Adr() {
					nm := c.ml.Loaded(n.Name)
					nt := nm.TypeByName(nm.Enter, it.Name())
					return nm, nt
				}
			}
		}
	}
	return nil, nil
}

func (c *tc) String() (ret string) {
	foreign := func(t object.ComplexType) {
		for _, n := range c.m.Imports {
			for _, _it := range n.Objects {
				switch it := _it.(type) {
				case object.TypeObject:
					if it.Complex().Adr() == t.Adr() {
						nm := c.ml.Loaded(n.Name)
						nt := nm.TypeByName(nm.Enter, it.Name())
						other := c.ml.NewTypeCalc()
						other.ConnectTo(nt)
						ret = fmt.Sprintln(ret, "foreign", other)
					}
				}
			}
		}
	}
	for t := c.typ; t != nil; {
		ret = fmt.Sprintln(ret, t.Qualident())
		ol := c.m.Objects[c.m.Enter]
		for _, _po := range ol {
			switch po := _po.(type) {
			case object.ProcedureObject:
				proc := c.m.NodeByObject(po)
				local := false
				for i := range proc {
					if _, ok := proc[i].(node.EnterNode); ok {
						local = true
					}
				}
				if po.Link() != nil {
					pt := po.Link().Complex()
					var pb object.ComplexType
					if _, ok := pt.(Inherited); ok {
						pb = pt.(Inherited).Complex()
					}
					if t.Equals(pt) || t.Equals(pb) {
						ret = fmt.Sprintln(ret, po.Name(), local)
					}
				}
			}
		}
		switch z := t.(type) {
		case object.PointerType:
			if z.Complex() != nil {
				t = z.Complex()
			} else {
				foreign(t)
				t = nil
			}
		case object.RecordType:
			if z.Complex() != nil {
				t = z.Complex()
			} else {
				foreign(t)
				t = nil
			}
		default:
			halt.As(0, reflect.TypeOf(t))
		}
	}
	return
}
