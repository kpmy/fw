package node

import "cp/object"

type Class int

const (
	ENTER Class = iota
	ASSIGN
	VARIABLE
	DYADIC
	CONSTANT
)

type Node interface {
	SetLeft(n Node)
	SetRight(n Node)
	SetLink(n Node)
	SetObject(o object.Object)

	Left() Node
	Right() Node
	Link() Node
	Object() object.Object
}

func New(class Class) Node {
	switch class {
	case ENTER:
		return new(enterNode)
	case ASSIGN:
		return new(assignNode)
	case VARIABLE:
		return new(variableNode)
	case DYADIC:
		return new(dyadicNode)
	case CONSTANT:
		return new(constantNode)
	default:
		panic("no such class")
	}
}

type nodeFields struct {
	left, right, link Node
	obj               object.Object
}

func (nf *nodeFields) SetLeft(n Node) { nf.left = n }

func (nf *nodeFields) SetRight(n Node) { nf.right = n }

func (nf *nodeFields) SetLink(n Node) { nf.link = n }

func (nf *nodeFields) SetObject(o object.Object) { nf.obj = o }

func (nf *nodeFields) Left() Node { return nf.left }

func (nf *nodeFields) Right() Node { return nf.right }

func (nf *nodeFields) Link() Node { return nf.link }

func (nf *nodeFields) Object() object.Object { return nf.obj }

type enterNode struct {
	nodeFields
	enter Enter
}

func (e *enterNode) SetEnter(enter Enter) {
	e.enter = enter
}

type constantNode struct {
	nodeFields
	typ  object.Type
	data interface{}
}

func (c *constantNode) SetType(t object.Type) {
	c.typ = t
}

func (c *constantNode) SetData(data interface{}) {
	c.data = data
}

func (c *constantNode) Data() interface{} {
	return c.data
}

func (c *constantNode) Type() object.Type {
	return c.typ
}

type dyadicNode struct {
	nodeFields
	operation Operation
}

func (d *dyadicNode) SetOperation(op Operation) {
	d.operation = op
}

func (d *dyadicNode) Operation() Operation {
	return d.operation
}

type assignNode struct {
	nodeFields
}

func (a *assignNode) Self() AssignNode {
	return a
}

type variableNode struct {
	nodeFields
}

func (v *variableNode) Self() VariableNode {
	return v
}
