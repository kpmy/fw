package scope

import (
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
)

//менеджер зон видимости, зоны видимости динамические, создаются в момент входа в EnterNode
type Manager interface {
	context.ContextAware
	UpdateObj(o object.Object, val ValueFor)
	SelectObj(o object.Object) interface{}
	UpdateNode(n node.Node, val ValueFor)
	SelectNode(n node.Node) interface{}
	FindObjByName(name string) object.Object
	Allocate(n node.Node)
	Dispose(n node.Node)
	Initialize(n node.Node, o object.Object, val node.Node)
}

//средство обновления значения
type ValueFor func(in interface{}) (out interface{})
