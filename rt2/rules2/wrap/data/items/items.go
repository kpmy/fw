package items

import (
	"container/list"
	"fmt"
	"ypk/assert"
	"ypk/halt"
)

type Opts int

const INIT Opts = iota

type ID struct {
	In   Data
	This Key
}

type Key interface {
	EqualTo(Key) int
}

type Value interface {
	KeyOf(...Key) Key
}

type Link interface {
	To(...ID) ID
	Value
}

type Item interface {
	Data(...interface{}) interface{}
	Value
	Copy(Item)
}

type Data interface {
	Set(Key, Item)
	Get(Key, ...Opts) Item
	Hold(Key)
	Begin()
	End()
	Drop()
}

func New() Data {
	return &data{x: list.New()}
}

type dummy struct {
	k Key
}

func (d *dummy) String() string {
	return fmt.Sprint(d.k)
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

func (d *data) Get(k Key, opts ...Opts) (ret Item) {
	if len(opts) == 0 {
		d.check()
	} else {
		switch opts[0] {
		case INIT: //do nothing
		default:
			halt.As(100, fmt.Sprint(opts))
		}
	}
	for x, e := d.find(k, nil); x != nil && ret == nil; {
		switch v := x.(type) {
		case nil: //do nothing
		case Item:
			ret = v
		case Link:
			to := v.To()
			if to.In == d {
				x, e = d.find(to.This, e)
			} else {
				ret = to.In.Get(to.This)
			}
		}
	}
	assert.For(ret != nil, 60, k)
	return
}

func (d *data) Hold(key Key) {
	assert.For(key != nil, 20)
	d.x.PushFront(&dummy{k: key})
}

type link struct {
	k  Key
	id ID
}

func (l *link) KeyOf(k ...Key) Key {
	if len(k) == 1 {
		l.k = k[0]
	}
	return l.k
}

func (l *link) To(id ...ID) ID {
	if len(id) == 1 {
		l.id = id[0]
	}
	return l.id
}

func NewLink(to ID) Link {
	return &link{id: to}
}

type limit struct{}
type limit_key struct{}

func (l *limit_key) EqualTo(Key) int { return -1 }
func (l *limit) KeyOf(...Key) Key    { return &limit_key{} }

func (d *data) check() {
	t := d.x.Front()
	if t != nil {
		_, ok := t.Value.(*limit)
		assert.For(ok, 30)
	}
}

func (d *data) Begin() {
	d.check()
}

func (d *data) End() {
	for x := d.x.Front(); x != nil; {
		d, ok := x.Value.(*dummy)
		assert.For(!ok, 40, "missing value for item ", d)
		x = x.Next()
		if x != nil {
			if _, ok := x.Value.(*limit); ok {
				x = nil
			}
		}
	}
	d.x.PushFront(&limit{})
}
func (d *data) Drop() {
	d.check()
	for x := d.x.Front(); x != nil; {
		d.x.Remove(x)
		x = d.x.Front()
		if _, ok := x.Value.(*limit); ok {
			x = nil
		}
	}
}
