package scope

import (
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
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

//менеджер зон видимости, зоны видимости динамические, создаются в момент входа в EnterNode
type Manager interface {
	context.ContextAware
	Update(id ID, val ValueFor)
	Select(id ID) interface{}
	Allocate(n node.Node, final bool)
	Dispose(n node.Node)
	Initialize(n node.Node, o object.Object, val node.Node)
}

//средство обновления значения
type ValueFor func(in interface{}) (out interface{})

var Designator func(n node.Node) ID
var FindObjByName func(m Manager, name string) object.Object

func This(i interface{}) Manager {
	return i.(Manager)
}

var New func() Manager
