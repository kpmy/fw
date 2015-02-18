package scope

import (
	"fw/cp"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
)

type PARAM struct {
	Objects object.Object
	Values  node.Node
	Frame   frame.Frame
	Tail    frame.Sequence
}

//менеджер зон видимости, зоны видимости динамические, создаются в момент входа в EnterNode
// pk, 20150112, инициализация параметров теперь происходит как и обычный frame.Sequence, с использованием стека
type Manager interface {
	context.ContextAware
	Select(cp.ID, ValueOf)
	Target(...Allocator) Allocator
	Provide(interface{}) Value
	String() string
	Exists(cp.ID) bool
}

type Allocator interface {
	Join(Manager)
}

type ScopeAllocator interface {
	Allocator
	Allocate(n node.Node, final bool)
	Dispose(n node.Node)
	Initialize(n node.Node, par PARAM) (frame.Sequence, frame.WAIT)
}

type HeapAllocator interface {
	Allocator
	Allocate(string, object.PointerType, ...interface{}) Value //указатель лежит в скоупе процедуры/модуля, а рекорд - в куче, поэтому нужно после создания экземпляра обновить указатель
	Dispose(id cp.ID)
}

var FindObjByName func(m Manager, name string) object.Object

var New func(role string) Manager
