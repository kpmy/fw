package node

import (
	"fw/cp"
	"fw/cp/constant/enter"
	"fw/cp/constant/operation"
	"fw/cp/constant/statement"
	"fw/cp/object"
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

	cp.Id
}

type Statement interface {
	s() Statement
}

type Expression interface {
	e() Expression
}

type Designator interface {
	d() Designator
}

type EnterNode interface {
	Enter() enter.Enter
	SetEnter(enter enter.Enter)
	Node
	Statement
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
	Expression
}

// Self-designator for empty interfaces
type AssignNode interface {
	self() AssignNode
	SetStatement(statement.Statement)
	Statement() statement.Statement
	Node
	Statement
}

type VariableNode interface {
	self() VariableNode
	Node
	Designator
}

type CallNode interface {
	self() CallNode
	Node
	Statement
	Expression
}

type ProcedureNode interface {
	self() ProcedureNode
	Super(...string) bool
	Node
	Designator
}

type ParameterNode interface {
	Node
	self() ParameterNode
	Designator
}

type ReturnNode interface {
	Node
	self() ReturnNode
	Statement
}

type DyadicNode interface {
	OperationNode
	self() DyadicNode
	Expression
}

type MonadicNode interface {
	OperationNode
	Expression
	SetType(typ object.Type)
	Type() object.Type
	self() MonadicNode
	Complex(...object.ComplexType) object.ComplexType
}

type ConditionalNode interface {
	self() ConditionalNode
	Node
	Statement
}

type IfNode interface {
	self() IfNode
	Node
}

type WhileNode interface {
	self() WhileNode
	Node
	Statement
}

type RepeatNode interface {
	self() RepeatNode
	Node
	Statement
}

type ExitNode interface {
	self() ExitNode
	Node
	Statement
}

type LoopNode interface {
	self() LoopNode
	Node
	Statement
}

type DerefNode interface {
	self() DerefNode
	Node
	Ptr(...string) bool
	Designator
}

type FieldNode interface {
	self() FieldNode
	Node
	Designator
}

type IndexNode interface {
	self() IndexNode
	Node
	Designator
}

type TrapNode interface {
	self() TrapNode
	Node
	Statement
}

type WithNode interface {
	self() WithNode
	Node
	Statement
}

type GuardNode interface {
	self() GuardNode
	Node
	Type() object.ComplexType
	SetType(object.ComplexType)
	Designator
}

type CaseNode interface {
	self() CaseNode
	Node
	Statement
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

type RangeNode interface {
	self() RangeNode
	Node
}
