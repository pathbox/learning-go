package skiplist

// 调表的时间复杂度：O(logn)  底层是一个链表或双向链表，上部分有多层级的索引链表
const MaxLevel = 32

const p = 0.5 // 向右或向下的概率(每个节点晋升到下一级索引的概率)，一般来说0.5是一个合适的选择
// 当节点数量足够大实话，建立的索引也就足够分享，就越接近严格的每两个节点中有一个晋升的效果

type Node struct {
	value  uint32
	levels []*Level // 索引节点
}

type Level struct {
	next *Node
}

type SkipList struct {
	header *Node
	length uint32 // 原始链表长度，表头节点不计入
	height uint32 // 最高的节点的层数
}

func NewSkipList() *SkipList {
	return &SkipList{
		header: NewNode(MaxLevel, 0),
		length: 0,
		height: 1,
	}
}

func NewNode(level, value uint32) *Node {
	node := new(Node)
	node.value = value
	node.levels = make([]*Level, level)

	for i := 0; i < len(node.levels); i++ {
		node.levels[i] = new(Level)
	}
	return node
}

// 这样就有： 1/2的概率level为1 1/4的概率level为2 ...
func (s1 *SkipList) randomLevel() int {
	level := 1
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for r.Float64() < p && level < MaxLevel { // 随机一个0-1的数，如果小于晋升概率p，且总层数不大于最大层数
		level++
	}
	return level
}

func (s1 *SkipList) Add(value uint32) bool {
	if value <= 0 {
		return false
	}
	update := make([]*Node, MaxLevel)
	tmp := s1.header
	// 每一次循环都是一次寻找有序单链表的插入过程
	for i := int(s1.height) - 1; i >= 0; i-- {
		// 每次循环不充值tmp，直接从上一层确认的节点开始向下一层查找
		for tmp.levels[i].next != nil && tmp.levels[i].next.value < value {
			tmp = tmp.levels[i].next
		}
		// 避免插入相同元素
		if tmp.levels[i].next != nil && tmp.levels[i].next.value == value {
			return false
		}

		update[i] = tmp
	}

	level := s1.randomLevel()

	node := NewNode(uint32(level), value)

	if uint32(level) > s1.height {
		s1.height = uint32(level)
	}

	for i := 0; i < level; i++ {
		// 说明新节点层数超过了跳表当前的最高层数，此时将头节点对应层数的后继节点设置为新节点
		if update[i] == nil {
			s1.header.levels[i].next = node
			continue
		}

		node.levels[i].next = update[i].levels[i].next
		update[i].levels[i].next = node
	}
	s1.length++
	return true

	/* 从头节点的最高层开始查询，每次循环都可以理解为一次寻找有序单链表插入位置的过程。

	找到在这层索引的插入位置，存入 update 数组中。

	遍历完一层后，直接使用这一层查到的节点，即代码中的 tmp 开始遍历下一层索引。

	重复1-3步直到结束。

	获取新节点的层数，以确定从哪一层开始插入。如果层数大于跳表当前的最高层数，修改当前最高层数。

	遍历 update 数组，但只遍历到新节点的最大层数。

	增加跳表长度，返回 true 表示新增成功。
	update 长度为 5

	那么会从 3 层开始向下遍历，在二级索引这层找到 9 应该插入的位置——1 和 10 之间，update[2] 记录包含 1 的节点。

	在一级索引这层找到 9 应该插入的位置——7 和 10 之间，update[1] 记录了包含 7 的节点。

	在原链表这层找到 9 应该插入的位置——8 和 10 之间，update[0] 记录了包含 8 的节点。
	*/

}

func (s1 *SkipList) Find(value uint32) *Node {
	var node *Node
	tmp := s1.header
	for i := int(s1.height) - 1; i >= 0; i-- {
		for tmp.levels[i].next != nil && tmp.levels[i].next.value <= value {
			tmp = tmp.levels[i].next
		}
		if tmp.value == value {
			node = tmp
			break
		}
	}
	return node
}
