package main

import (
	"fmt"
	"sync"
)

type NodeQueue struct {
	nodes []Node
	lock  sync.RWMutex
}

type Node struct {
	value int
}

type Graph struct {
	nodes []*Node
	edges map[Node][]*Node // 使用map：邻接表表示的无向图
	lock  sync.RWMutex
}

func (g *Graph) BFS(f func(node *Node)) {
	g.lock.RLock()
	defer g.lock.RUnlock()

	q := NewNodeQueue()
	// 取图的第一个节点入队列
	head := g.nodes[0]
	q.Enqueue(*head)

	// 标识节点是否已经被访问过
	visited := make(map[*Node]bool)
	visited[head] = true

	// 遍历所有节点直到队列为空
	for {
		if q.IsEmpty() {
			break
		}
		node := q.Dequeue()
		visited[node] = true
		nexts := g.edges[*node]
		// 将所有未访问过的邻接节点入队列
		for _, next := range nexts {
			// 如果节点已经被访问过
			if visited[next] {
				continue
			}
			q.Enqueue(*next)
			visited[next] = true
		}
		// 对每个正在遍历的节点执行回调
		if f != nil {
			f(node)
		}
	}
}

func NewNodeQueue() *NodeQueue {
	q := NodeQueue{}
	q.lock.Lock()
	defer q.lock.Unlock()
	q.nodes = []Node{}
	return &q
}

func (q *NodeQueue) Enqueue(n Node) {
	q.lock.Lock()
	defer q.lock.Unlock()
	q.nodes = append(q.nodes, n)
}

func (q *NodeQueue) Dequeue() *Node {
	q.lock.Lock()
	defer q.lock.Unlock()
	node := q.nodes[0]
	q.nodes = q.nodes[1:]
	return &node
}

func (q *NodeQueue) IsEmpty() bool {
	q.lock.RLock()
	defer q.lock.RUnlock()
	return len(q.nodes) == 0
}

func (g *Graph) AddNode(n *Node) {
	g.lock.Lock()
	defer g.lock.Unlock()
	g.nodes = append(g.nodes, n)
}

func (g *Graph) AddEdge(u, v *Node) {
	g.lock.Lock()
	defer g.lock.Unlock()

	// 首次建立图
	if g.edges == nil {
		g.edges = make(map[Node][]*Node)
	}
	g.edges[*u] = append(g.edges[*u], v) // 建立 u->v 的边
	g.edges[*v] = append(g.edges[*v], u) // 无向图，建立 v -> u的边
}

func main() {
	g := Graph{}
	n1, n2, n3, n4, n5 := Node{1}, Node{2}, Node{3}, Node{4}, Node{5}

	g.AddNode(&n1)
	g.AddNode(&n2)
	g.AddNode(&n3)
	g.AddNode(&n4)
	g.AddNode(&n5)

	g.AddEdge(&n1, &n2)
	g.AddEdge(&n1, &n5)
	g.AddEdge(&n2, &n3)
	g.AddEdge(&n2, &n4)
	g.AddEdge(&n2, &n5)
	g.AddEdge(&n3, &n4)
	g.AddEdge(&n4, &n5)

	g.BFS(func(node *Node) {
		fmt.Printf("[Current Traverse Node]: %v\n", node)
	})
}
