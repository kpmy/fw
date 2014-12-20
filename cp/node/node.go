package node

import (
	"cp/constant"
	"cp/object"
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
