package scope

import (
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/frame"
	"strconv"
)

const DEPTH = 16

type ID struct {
	Name  string
	Path  [DEPTH]string
	Index *int64
}

func (i ID) String() string {
	if i.Name != "" {
		ret := i.Name
		if i.Path[0] != "" {
			ret = ret + "." + i.Path[0]
		}
		if i.Index != nil {
			ret = ret + "[" + strconv.FormatInt(*i.Index, 10) + "]"
		}
		return ret
	} else {
		return "<empty id>"
	}
}

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
	Update(id ID, val ValueFor)
	Select(id ID) interface{}
	Target(...Allocator) Allocator
}

type Allocator interface{}

type ScopeAllocator interface {
	Allocator
	Allocate(n node.Node, final bool)
	Dispose(n node.Node)
	Initialize(n node.Node, par PARAM) (frame.Sequence, frame.WAIT)
}

type HeapAllocator interface {
	Allocator
}

//средство обновления значения
type ValueFor func(in interface{}) (out interface{})

var Designator func(n ...node.Node) ID
var FindObjByName func(m Manager, name string) object.Object

func This(i interface{}) Manager {
	return i.(Manager)
}

var NewStack func() Manager
var NewHeap func() Manager
