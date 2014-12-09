package node

import (
	"cp/constant/enter"
	"cp/constant/operation"
	"cp/object"
	"cp/statement"
)

type EnterNode interface {
	SetEnter(enter enter.Enter)
}

type OperationNode interface {
	SetOperation(op operation.Operation)
	Operation() operation.Operation
}

type ConstantNode interface {
	SetType(typ object.Type)
	SetData(data interface{})
	Data() interface{}
	Type() object.Type
}

// Self-designator for empty interfaces
type AssignNode interface {
	Self() AssignNode
	SetStatement(statement.Statement)
	Statement() statement.Statement
}

type VariableNode interface {
	Self() VariableNode
}

type CallNode interface {
	Self() CallNode
}

type ProcedureNode interface {
	Self() ProcedureNode
}

type enterNode struct {
	nodeFields
	enter enter.Enter
}

func (e *enterNode) SetEnter(enter enter.Enter) {
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
	operation operation.Operation
}

func (d *dyadicNode) SetOperation(op operation.Operation) {
	d.operation = op
}

func (d *dyadicNode) Operation() operation.Operation {
	return d.operation
}

type assignNode struct {
	nodeFields
	stat statement.Statement
}

func (a *assignNode) Self() AssignNode {
	return a
}

func (a *assignNode) SetStatement(s statement.Statement) {
	a.stat = s
}

func (a *assignNode) Statement() statement.Statement {
	return a.stat
}

type variableNode struct {
	nodeFields
}

func (v *variableNode) Self() VariableNode {
	return v
}

type callNode struct {
	nodeFields
}

func (v *callNode) Self() CallNode {
	return v
}

type procedureNode struct {
	nodeFields
}

func (v *procedureNode) Self() ProcedureNode {
	return v
}
