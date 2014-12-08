package node

import (
	"cp/object"
	"cp/statement"
)

type Enter int

const (
	MODULE Enter = iota
)

type Operation int

const (
	PLUS Operation = iota
)

type EnterNode interface {
	SetEnter(enter Enter)
}

type OperationNode interface {
	SetOperation(op Operation)
	Operation() Operation
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
