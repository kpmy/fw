package std

import (
	"fmt"
	"fw/cp/node"
	"fw/cp/object"
	"fw/rt2/context"
	"fw/rt2/scope"
	"reflect"
	"ypk/assert"
)

type heap struct {
	d    context.Domain
	data *area
	next int
}

type desc struct {
	scope.ID
	link object.Object
}

func (d *desc) Object() object.Object {
	return d.link
}

func nh() scope.Manager {
	return &heap{data: &area{ready: true, root: nil, x: make(map[scope.ID]interface{})}}
}

func (h *heap) Allocate(n node.Node) scope.ValueFor {
	switch v := n.(type) {
	case node.VariableNode:
		switch t := v.Object().Complex().(type) {
		case object.PointerType:
			h.next++
			switch bt := t.Base().(type) {
			case object.RecordType:
				fake := object.New(object.VARIABLE)
				fake.SetComplex(bt)
				fake.SetType(object.COMPLEX)
				fake.SetName("@")
				r := &rec{link: fake}
				oid := scope.ID{Name: "@"}
				oid.Ref = new(int)
				*oid.Ref = h.next
				od := &desc{link: fake}
				od.ID = oid
				fake.SetRef(od)
				alloc(nil, h.data, oid, r)
				return func(interface{}) interface{} {
					return oid
				}
			default:
				panic(fmt.Sprintln("cannot allocate", reflect.TypeOf(t)))
			}
		default:
			panic(fmt.Sprintln("unsupported type", reflect.TypeOf(t)))
		}
	default:
		panic(fmt.Sprintln("unsupported node", reflect.TypeOf(v)))
	}
}

func (h *heap) Dispose(n node.Node) {
}

func (h *heap) Target(...scope.Allocator) scope.Allocator {
	return h
}

func (h *heap) Update(i scope.ID, val scope.ValueFor) {
	fmt.Println("update", i)
	var upd func(x interface{}) (ret interface{})
	upd = func(x interface{}) (ret interface{}) {
		fmt.Println(reflect.TypeOf(x))
		switch x := x.(type) {
		case value:
			old := x.get()
			tmp := val(old)
			assert.For(tmp != nil, 40) //если устанавливают значение NIL, значит делают что-то неверно
			//fmt.Println(tmp)
			x.set(tmp)
			ret = x
		case record:
			if i.Path == "" {
				//fmt.Println(i, depth)
				panic(0) //случай выбора всей записи целиком
			} else {
				z := x.getField(i.Path)
				ret = upd(z)
			}
		default:
			panic(0)
		}
		return ret
	}
	upd(h.data.get(i))
}

func (h *heap) Select(i scope.ID) interface{} {
	fmt.Println("heap select", i)
	type result struct {
		x interface{}
	}
	var res *result
	var sel func(interface{}) *result
	sel = func(x interface{}) (ret *result) {
		fmt.Println(x)
		switch x := x.(type) {

		case record:
			if i.Path == "" {
				ret = &result{x: x.(*rec).link}
			} else {
				z := x.getField(i.Path)
				ret = sel(z)
			}
		default:
			panic(0)
		}
		return ret
	}
	res = sel(h.data.get(i))
	assert.For(res != nil, 40)
	//fmt.Println(res.x)
	return res.x
}

func (h *heap) Init(d context.Domain) { h.d = d }

func (h *heap) Domain() context.Domain { return h.d }

func (h *heap) Handle(msg interface{}) {}
