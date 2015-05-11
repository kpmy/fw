package items

import (
	"github.com/kpmy/ypk/assert"
)

type ID struct {
	In   Data
	This Key
}

type Key interface {
	EqualTo(Key) int
	Hash() int
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
	Get(Key) Item
	Remove(Key)

	Hold(Key)
	Link(Key, ID)
	Put(Key, Item)

	Exists(Key) bool
	ForEach(func(Value) bool)

	Begin()
	End()
	Drop()
}

func New() Data {
	return &tree{root: make([]*node, 0)}
}

type tree struct {
	root []*node
}

type node struct {
	data map[int]Value
}

func (t *tree) Begin() {
	n := &node{data: make(map[int]Value)}
	t.root = append(t.root, n)
}

func (t *tree) top() (ret *node) {
	if len(t.root) > 0 {
		ret = t.root[len(t.root)-1]
	}
	return
}

func (t *tree) End() {
	var this *node
	if len(t.root) > 0 {
		this = t.root[len(t.root)-1]
	}
	assert.For(this != nil, 20)
}

func (t *tree) Drop() {
	tmp := make([]*node, 0)
	for i := 0; i < len(t.root)-1; i++ {
		tmp = append(tmp, t.root[i])
	}
	t.root = tmp
}

func (t *tree) Set(k Key, i Item) {
	assert.For(i != nil, 20)
	x := t.top()
	x.data[k.Hash()] = i
	i.KeyOf(k)
}

func (t *tree) Get(k Key) (ret Item) {
	for i := len(t.root) - 1; i >= 0 && ret == nil; i-- {
		tmp := t.root[i].data[k.Hash()]
		switch this := tmp.(type) {
		case Item:
			ret = tmp.(Item)
		case Link:
			ret = this.To().In.Get(this.To().This)
		}
	}
	return
}

func (t *tree) Remove(k Key) {
	var tmp Value
	for i := len(t.root) - 1; i >= 0 && tmp != nil; i-- {
		tmp = t.root[i].data[k.Hash()]
		if tmp != nil {
			delete(t.root[i].data, k.Hash())
		}
	}
}

type dummy struct {
	k Key
}

func (d *dummy) KeyOf(...Key) Key { return d.k }

func (t *tree) Hold(key Key) {
	n := t.top()
	n.data[key.Hash()] = &dummy{k: key}
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

func (t *tree) Link(key Key, to ID) {
	l := &link{k: key, id: to}
	n := t.top()
	n.data[key.Hash()] = l
}

func (t *tree) Put(k Key, i Item) {
	t.Set(k, i)
}

func (t *tree) Exists(k Key) (ret bool) {
	for i := len(t.root) - 1; i >= 0 && !ret; i-- {
		ret = t.root[i].data[k.Hash()] != nil
	}
	return
}

func (t *tree) ForEach(f func(Value) bool) {
	ok := false
	for i := len(t.root) - 1; i >= 0 && !ok; i-- {
		for _, v := range t.root[i].data {
			ok = f(v)
		}
	}
}
