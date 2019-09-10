package memtable

import (
	"./internal"
	"./skiplist"
)

type Iterator struct {
	listIter *skiplist.Iterator
}

func (it *Iterator) Valid() bool {
	return it.listIter.Valid()
}

func (it *Iterator) InternalKey() *internal.InternalKey {
	return it.listIter.Key().(*internal.InternalKey)
}

func (it *Iterator) Next() {
	it.listIter.Next()
}

func (it *Iterator) Prev() {
	it.listIter.Prev()
}

func (it *Iterator) Seek(target interface{}) {
	it.listIter.Seek(target)
}

func (it *Iterator) SeekToLast() {
	it.listIter.SeekToLast()
}
