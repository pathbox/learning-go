package main

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

const DEFAULT_REPLICAS = 256

type HashRing []uint32 // hash ring 用slice 数组数据结构存储, 存的是每个虚拟节点的 key的值

// 实现 sort包接口
func (c HashRing) Len() int {
	return len(c)
}

func (c HashRing) Less(i, j int) bool {
	return c[i] < c[j]
}

func (c HashRing) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

type Node struct {
	ID       int
	IP       string
	Port     int
	HostName string
	Weight   int
}

func NewNode(id int, ip string, port int, name string, weight int) *Node {
	return &Node{
		ID:       id,
		IP:       ip,
		Port:     port,
		HostName: name,
		Weight:   weight,
	}
}

// 一致性结构
type Consistent struct {
	Nodes     map[uint32]Node // 虚拟IP节点存储
	numReps   int             // 用于构造虚拟节点的倍数因子
	Resources map[int]bool    // 真实IP节点存储
	ring      HashRing
	sync.RWMutex
}

func NewConsistent() *Consistent {
	return &Consistent{
		Nodes:     make(map[uint32]Node),
		numReps:   DEFAULT_REPLICAS,
		Resources: make(map[int]bool),
		ring:      HashRing{},
	}
}

func (c *Consistent) Add(node *Node) bool {
	c.Lock() // 加锁
	defer c.Unlock()

	if _, ok := c.Resources[node.ID]; ok { // Resources key是node.ID,如果这个ID存在，则已经存在Resources中了
		return false
	}

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)         // 根据这个node参数，得到一个字符串
		c.Nodes[c.hashStr(str)] = *(node) // c.Nodes 中存的是虚拟node节点数量，某个真实的节点对应count(c.numReps * node.Weight)个虚拟节点, 虚拟节点的key是通过IP组合字符串hash函数得到的一个hash uint32 数值
	}
	c.Resources[node.ID] = true // c.Resources 存真实node
	c.sortHashRing()
	return true
}

func (c *Consistent) sortHashRing() {
	c.ring = HashRing{}
	for k := range c.Nodes { // 将每个真实node的虚拟节点存到hashRing环中, k 是虚拟节点的key uint32 hash值
		c.ring = append(c.ring, k)
	}
	sort.Sort(c.ring)
}

func (c *Consistent) joinStr(i int, node *Node) string {
	return node.IP + "*" + strconv.Itoa(node.Weight) +
		"-" + strconv.Itoa(i) +
		"-" + strconv.Itoa(node.ID)
}

func (c *Consistent) hashStr(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key)) // hash函数 得到uint32hash值
}

func (c *Consistent) Get(key string) Node {
	c.RLock()
	defer c.RUnlock()

	hash := c.hashStr(key)    // 将 key字符串 hash函数得到一个 hash值 uint32数值
	i := c.search(hash)       // 在 c.ring 中根据这个hash值进行范围匹配得到一个 int索引值，这个值得范围是 0-len(c.ring)
	return c.Nodes[c.ring[i]] //c.ring 中存的是虚拟节点的hash值key，返回一个虚拟节点
}

// hash 是key字符串的hash uint32数值
func (c *Consistent) search(hash uint32) int {
	i := sort.Search(len(c.ring), func(i int) bool { return c.ring[i] >= hash }) // 通过二叉查找算法，找到满足 c.ring[i] >= hash 的最小的i值，如果没有找到满足条件的i，则返回最后的i值，这个i值在 0-len(c.ring)之间
	if i < len(c.ring) {
		if i == len(c.ring)-1 {
			return 0 // 环形列表，重新回到头部
		} else {
			return i
		}
	} else {
		return len(c.ring) - 1 // 返回 最后一位
	}
}

// 删除节点： 从Resource中删除实际节点，从c.Nodes中删除虚拟节点
func (c *Consistent) Remove(node *Node) {
	c.Lock()
	defer c.Unlock()

	if _, ok := c.Resources[node.ID]; !ok {
		return
	}

	delete(c.Resources, node.ID)

	count := c.numReps * node.Weight
	for i := 0; i < count; i++ {
		str := c.joinStr(i, node)
		delete(c.Nodes, c.hashStr(str))
	}
	c.sortHashRing()
}

func main() {
	cHashRing := NewConsistent()

	for i := 0; i < 10; i++ { // 增加10个节点
		si := fmt.Sprintf("%d", i)
		cHashRing.Add(NewNode(i, "172.18.1."+si, 8080, "host_"+si, 1))
	}

	for k, v := range cHashRing.Nodes {
		fmt.Println("Hash:", k, " IP:", v.IP)
	}

	ipMap := make(map[string]int)
	for i := 0; i < 1000; i++ {
		si := fmt.Sprintf("key%d", i) // 构造一个 key值,就是一个需要的字符串值
		k := cHashRing.Get(si)        // 根据这个key值，去hashRing中找到一个最近的虚拟节点key的值，根据这个key值找到一个虚拟节点，这个虚拟节点对应一个真实节点数据
		if _, ok := ipMap[k.IP]; ok {
			ipMap[k.IP] += 1
		} else {
			ipMap[k.IP] = 1
		}
	}
	for k, v := range ipMap {
		fmt.Println("Node IP:", k, " count:", v)
	}
}
