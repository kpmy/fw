package scope

import (
	"cp/node"
	"cp/object"
	"rt2/context"
)

//менеджер зон видимости, зоны видимости динамические, создаются в момент входа в EnterNode
type Manager interface {
	context.ContextAware
	Update(o object.Object, val Value)
	Select(o object.Object) interface{}
	Allocate(n node.Node)
	Dispose(n node.Node)
}

//средство обновления значения
type Value func(in interface{}) (out interface{})
