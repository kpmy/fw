package scope

import (
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
)

//менеджер зон видимости, зоны видимости динамические, создаются в момент входа в EnterNode
type Manager interface {
	context.ContextAware
	Update(o object.Object, val ValueFor)
	Select(o object.Object) interface{}
	FindObjByName(name string) object.Object
	Allocate(n node.Node)
	Dispose(n node.Node)
	Initialize(n node.Node, o object.Object, val node.Node)
}

//средство обновления значения
type ValueFor func(in interface{}) (out interface{})
