package node

import (
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/object"
	"fw/cp/statement"
)

type EnterNode interface {
	Enter() enter.Enter
	SetEnter(enter enter.Enter)
	Node
}

type OperationNode interface {
	SetOperation(op operation.Operation)
	Operation() operation.Operation
	Node
}

type ConstantNode interface {
	SetType(typ object.Type)
	SetData(data interface{})
	Data() interface{}
	Type() object.Type
	Node
}

// Self-designator for empty interfaces
type AssignNode interface {
	self() AssignNode
	SetStatement(statement.Statement)
	Statement() statement.Statement
	Node
}

type VariableNode interface {
	self() VariableNode
	Node
}

type CallNode interface {
	self() CallNode
	Node
}

type ProcedureNode interface {
	self() ProcedureNode
	Node
}

type ParameterNode interface {
	Node
	self() ParameterNode
}

type ReturnNode interface {
	Node
	self() ReturnNode
}

type DyadicNode interface {
	OperationNode
	self() DyadicNode
}

type MonadicNode interface {
	OperationNode
	SetType(typ object.Type)
	Type() object.Type
	self() MonadicNode
}

type ConditionalNode interface {
	self() ConditionalNode
}

type IfNode interface {
	self() IfNode
}

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
