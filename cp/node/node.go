package node

import (
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/object"
	"fw/cp/statement"
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
	Min() *int
	Max() *int
	SetMin(int)
	SetMax(int)
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
	Node
}

type IfNode interface {
	self() IfNode
	Node
}

type WhileNode interface {
	self() WhileNode
	Node
}

type RepeatNode interface {
	self() RepeatNode
	Node
}

type ExitNode interface {
	self() ExitNode
	Node
}

type LoopNode interface {
	self() LoopNode
	Node
}

type DerefNode interface {
	self() DerefNode
	Node
}

type FieldNode interface {
	self() FieldNode
	Node
}

type IndexNode interface {
	self() IndexNode
	Node
}

type TrapNode interface {
	self() TrapNode
	Node
}

type WithNode interface {
	self() WithNode
	Node
}

type GuardNode interface {
	self() GuardNode
	Node
	Type() object.ComplexType
	SetType(object.ComplexType)
}

type CaseNode interface {
	self() CaseNode
	Node
}

type ElseNode interface {
	Node
	Min(...int) int
	Max(...int) int
	//	SetMin(int)
	//	SetMax(int)

}

type DoNode interface {
	self() DoNode
	Node
}
