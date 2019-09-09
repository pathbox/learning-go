package memtable

import (
	"./internal"
	"./skiplist"
)

type MemTable struct {
	table       *skiplist.SkipList
	memoryUsage uint64
}

func New() *MemTable {
	var memTable MemTable
	memTable.table = skiplist.New(internal.InternalKeyComparator)
	return &memTable
}

func (memTable *MemTable) NewIterator() *Iterator {
	return &Iterator{listIter: memTable.table.NewIterator()}
}

func (memTable *MemTable) Add(seq uint64, valueType internal.ValueType, key, value []byte) {
	internalKey := internal.NewInternalKey(seq, valueType, key, value)
	memTable.memoryUsage += uint64(16 + len(key) + len(value))
	memTable.table.Insert(internalKey)
}

func (memTable *MemTable) Get(key []byte) ([]byte, error) {
	lookupKey := internal.LookupKey(key)
	it := memTable.table.NewIterator()
	it.Seek(lookupKey)
	if it.Valid() {
		internalKey := it.Key().(*internal.InternalKey)
		if internal.UserKeyComparator(key, internalKey.UserKey) == 0 {
			if internalKey.Type == internal.TypeValue {
				return internalKey.UserValue, nil
			} else {
				return nil, internal.ErrDeletion
			}
		}
	}
	return nil, internal.ErrNotFound
}

func (memTable *MemTable) ApproximateMemoryUsage() uint64 {
	return memTable.memoryUsage
}
