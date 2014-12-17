package scope

import (
	"cp/node"
	"cp/object"
	"rt2/context"
)

//менеджер зон видимости
type Manager interface {
	context.ContextAware
	Calculate(n node.Node) Area
	Allocate(n node.Node)
	Dispose(n node.Node)
}

//зона видимости
type Area interface {
	Get(o object.Object) Object
	Set(o Object)
}

//объект зоны видимости
type Object interface{}
