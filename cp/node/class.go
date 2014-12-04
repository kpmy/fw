package node

import "cp/object"

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
}

type ConstantNode interface {
	SetType(typ object.Type)
	SetData(data interface{})
}
