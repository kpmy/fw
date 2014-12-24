package node

import (
	"fw/cp/constant"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/object"
	"fw/cp/statement"
)

func New(class constant.Class) Node {
	switch class {
	case constant.ENTER:
		return new(enterNode)
	case constant.ASSIGN:
		return new(assignNode)
	case constant.VARIABLE:
		return new(variableNode)
	case constant.DYADIC:
		return new(dyadicNode)
	case constant.CONSTANT:
		return new(constantNode)
	case constant.CALL:
		return new(callNode)
	case constant.PROCEDURE:
		return new(procedureNode)
	case constant.PARAMETER:
		return new(parameterNode)
	case constant.RETURN:
		return new(returnNode)
	case constant.MONADIC:
		return new(monadicNode)
	case constant.CONDITIONAL:
		return new(conditionalNode)
	case constant.IF:
		return new(ifNode)
	case constant.REPEAT:
		return new(repeatNode)
	case constant.WHILE:
		return new(whileNode)
	case constant.EXIT:
		return new(exitNode)
	case constant.LOOP:
		return new(loopNode)
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
	enter enter.Enter
}

func (e *enterNode) SetEnter(enter enter.Enter) { e.enter = enter }
func (e *enterNode) Enter() enter.Enter         { return e.enter }

type constantNode struct {
	nodeFields
	typ  object.Type
	data interface{}
}

func (c *constantNode) SetType(t object.Type) { c.typ = t }

func (c *constantNode) SetData(data interface{}) { c.data = data }

func (c *constantNode) Data() interface{} { return c.data }

func (c *constantNode) Type() object.Type { return c.typ }

type dyadicNode struct {
	nodeFields
	operation operation.Operation
}

func (d *dyadicNode) SetOperation(op operation.Operation) { d.operation = op }

func (d *dyadicNode) Operation() operation.Operation { return d.operation }

func (d *dyadicNode) self() DyadicNode { return d }

type assignNode struct {
	nodeFields
	stat statement.Statement
}

func (a *assignNode) self() AssignNode { return a }

func (a *assignNode) SetStatement(s statement.Statement) { a.stat = s }

func (a *assignNode) Statement() statement.Statement { return a.stat }

type variableNode struct {
	nodeFields
}

func (v *variableNode) self() VariableNode { return v }

type callNode struct {
	nodeFields
}

func (v *callNode) self() CallNode { return v }

type procedureNode struct {
	nodeFields
}

func (v *procedureNode) self() ProcedureNode { return v }

type parameterNode struct {
	nodeFields
}

func (v *parameterNode) self() ParameterNode { return v }

type returnNode struct {
	nodeFields
}

func (v *returnNode) self() ReturnNode { return v }

type monadicNode struct {
	nodeFields
	operation operation.Operation
	typ       object.Type
}

func (v *monadicNode) self() MonadicNode { return v }

func (v *monadicNode) SetOperation(op operation.Operation) { v.operation = op }

func (v *monadicNode) Operation() operation.Operation { return v.operation }

func (v *monadicNode) SetType(t object.Type) { v.typ = t }
func (v *monadicNode) Type() object.Type     { return v.typ }

type conditionalNode struct {
	nodeFields
}

func (v *conditionalNode) self() ConditionalNode { return v }

type ifNode struct {
	nodeFields
}

func (v *ifNode) self() IfNode { return v }

type whileNode struct {
	nodeFields
}

func (v *whileNode) self() WhileNode { return v }

type repeatNode struct {
	nodeFields
}

func (v *repeatNode) self() RepeatNode { return v }

type exitNode struct {
	nodeFields
}

func (v *exitNode) self() ExitNode { return v }

type loopNode struct {
	nodeFields
}

func (v *loopNode) self() LoopNode { return v }
