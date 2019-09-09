package skiplist

type Iterator struct {
	list *SkipList
	node *Node
}

func (it *Iterator) Valid() bool {
	return it.node != nil
}

func (it *Iterator) Key() interface{} {
	return it.node.key
}

func (it *Iterator) Next() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.node.getNext(0)
}

// 往前移动一个节点
// Advances to the previous p	osition.
func (it *Iterator) Prev() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.findLessThan(it.node.key) // prev 比当前的key value小
	if it.node == it.list.head {
		it.node = nil
	}
}

// 找到目标key
func (it *Iterator) Seek(target interface{}) {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node, _ = it.list.findGreaterOrEqual(target)
}

func (it *Iterator) SeekToFirst() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.head.getNext(0)
}

// Position at the last entry in list.
// Final state of iterator is Valid() iff list is not empty.
func (it *Iterator) SeekToLast() {
	it.list.mu.RLock()
	defer it.list.mu.RUnlock()

	it.node = it.list.findlast()
	if it.node == it.list.head {
		it.node = nil
	}
}
