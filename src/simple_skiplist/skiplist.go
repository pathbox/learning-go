package skiplist

import (
	"math/rand"
	"sync"
	"time"
)

// 结点结构
type SkipListNode struct {
	key  int
	data interface{}
	next []*SkipListNode
}

// 跳表结构
type SkipList struct {
	head   *SkipListNode
	tail   *SkipListNode
	length int           // 数据总量
	level  int           // 层数
	mut    *sync.RWMutex // 读写锁
	rand   *rand.Rand    // 随机函数 随机数生成器，用于生成随机层数，随机生成的层数要满足P=0.5的几何分布，大致可以理解为：掷硬币连续出现正面的次数，就是我们要的层数
}

func (list *SkipList) randomLevel() int {
	level := 1
	for ; level < list.level && list.rand.Uint32()&0x1 == 1; level++ {
	}
	return level
}

func NewSkipList(level int) *SkipList {
	list := &SkipList{}
	if level <= 0 {
		level = 32
	}
	list.level = level
	list.head = &SkipListNode{next: make([]*SkipListNode, level, level)}
	list.tail = &SkipListNode{}
	list.mut = &sync.RWMutex{}
	list.rand = rand.New(rand.NewSource(time.Now().UnixNano()))
	for index := range list.head.next {
		list.head.next[index] = list.tail
	}
	return list
}

// 我们需要插入一个3，并且调用randomLevel得到的层数为3，那么插入3需要如下几部：

//     1、沿着LEVEL3的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——HEAD和层数——3。

//     2、沿着LEVEL2的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——&1和层数——2。

//     3、沿着LEVEL1的链查找第一个比3大的节点或者TAIL节点，记录下该节点的前一个节点——&2和层数——1。

//     4、生成一个新的节点newNode，key赋值为3，将newNode插入HEAD、&1、&2之后，即HEAD.next[3]=&3，&1.next[2]=&3，&2.next[1]=&3。

//     5、给newNode的next赋值，即步骤4中HEAD.next[3]、&1.next[2]、2.next[1]原本的值。

//     注意：为了易于理解，上述步骤中所有索引均从1开始，而代码中则从0开始，所以代码中均有索引=层数-1的关系。

func (list *SkipList) Add(key int, data interface{}) {
	list.mut.Lock()
	defer list.mut.Unlock()
	// 确定插入深度
	level := list.randomLevel()
	// 查找插入部位
	update := make([]*SkipListNode, level, level)
	node := list.head
	for index := level - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				update[index] = node // 找到一个插入部位
				break
			} else if node1.key == key {
				node1.data = data
				return
			} else {
				node = node1
			}
		}
	}

	// 执行插入
	newNode := &SkipListNode{key, data, make([]*SkipListNode, level, level)}
	for index, node := range update {
		node.next[index], newNode.next[index] = newNode, node.next[index]
	}
	list.length++
}

func (list *SkipList) Remove(key int) bool {
	list.mut.Lock()
	defer list.mut.Unlock()
	node := list.head
	remove := make([]*SkipListNode, list.level, list.level)
	var target *SkipListNode
	for index := len(node.next) - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				remove[index] = node //找到啦
				target = node1
				break
			} else {
				node = node1
			}
		}
	}
	//2.执行删除
	if target != nil {
		for index, node1 := range remove {
			if node1 != nil {
				node1.next[index] = target.next[index]
			}
		}
		list.length--
		return true
	}
	return false
}

func (list *SkipList) Find(key int) interface{} {
	list.mut.RLock()
	defer list.mut.Unlock()

	node := list.head

	for index := len(node.next) - 1; index >= 0; index-- {
		for {
			node1 := node.next[index]
			if node1 == list.tail || node1.key > key {
				break
			} else if node1.key == key {
				return node1.data
			} else {
				node = node1
			}
		}
	}
	return nil
}

func (list *SkipList) Length() int {
	list.mut.RLock()
	defer list.mut.RUnlock()
	return list.length
}
