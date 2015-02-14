package items

import (
	"container/list"
	"ypk/assert"
	"ypk/halt"
)

type Key interface {
	EqualTo(Key) int
}

type Value interface {
	KeyOf(...Key) Key
}

type Link interface {
	To() Key
	Value
}

type Item interface {
	Data(...interface{}) interface{}
	Value
	Copy(Item)
}

type Data interface {
	Set(Key, Item)
	Link(Key, Key)
	Get(Key) Item
	Limit()
	Drop()
}

func New() Data {
	return &data{x: list.New()}
}

type data struct {
	x *list.List
}

func (d *data) find(k Key, from *list.Element) (ret Value, elem *list.Element) {
	if from == nil {
		from = d.x.Front()
	}
	for x := from; x != nil && ret == nil; x = x.Next() {
		if v, ok := x.Value.(Value); ok {
			if z := v.KeyOf().EqualTo(k); z == 0 {
				ret = v
			}
		}
	}
	return
}

func (d *data) Set(k Key, v Item) {
	assert.For(v != nil, 20)
	assert.For(v.KeyOf() == nil, 21)
	x, _ := d.find(k, nil)
	if x == nil {
		v.KeyOf(k)
		d.x.PushFront(v)
	} else {
		halt.As(123)
	}
}

func (d *data) Get(k Key) (ret Item) {
	for x, e := d.find(k, nil); x != nil && ret == nil; {
		switch v := x.(type) {
		case nil: //do nothing
		case Item:
			ret = v
		case Link:
			x, e = d.find(v.To(), e)
		}
	}
	return
}

type link struct {
	k, t Key
}

func (l *link) KeyOf(k ...Key) Key {
	if len(k) == 1 {
		l.k = k[0]
	}
	return l.k
}

func (l *link) To() Key {
	return l.t
}

func (d *data) Link(key Key, to Key) {
	v, _ := d.find(key, nil)
	if v == nil {
		d.x.PushFront(&link{k: key, t: to})
	} else {
		halt.As(123)
	}
}

type limit struct{}
type limit_key struct{}

func (l *limit_key) EqualTo(Key) int { return -1 }
func (l *limit) KeyOf(...Key) Key    { return &limit_key{} }

func (d *data) Limit() { d.x.PushFront(&limit{}) }
func (d *data) Drop() {
	panic(0)
}
