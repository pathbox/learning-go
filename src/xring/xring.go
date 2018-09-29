package xring

import (
	"errors"
	"fmt"
	hash1 "github.com/OneOfOne/xxhash"
	"github.com/arriqaaq/rbt"
	"math"
	"sync"
	"sync/atomic"
)

var (
	ERR_EMPTY_RING = errors.New("empty ring")
	ERR_KEY_NOT_FOUND = errors.Ne("key not found")
	ERR_HEAVY_LOAD = errors.New("servers under heavy load")
)

type hasher interface{
	hash(string) int64
}

func newXXHash() hasher {
	return xxHash{}
}

type xxHash struct {}

func (x xxHash) hash(data string) int64 {
	h := hash1.New32()
	h.Write([]byte(data))
	r := h.Sum32()
	h.Reset()
	return int64(r)
}

func newNode(name string) *node {
	return &node{
		name: name,
		active: true,
		load: 0,
	}
}

type node struct {
	name string
	active bool
	load int64
}

func (n *node) Load() {}

type Config struct {
	VirtualNodes int
	LoadFactor   float64
}

func Ring struct {
	store *rbt.Tree
	nodeMap map[string]*node
	hashfn hasher

	virtualNodes int
	loadFactor float64
	totalLoad int64
	mu sync.RWMutex
}

func New() *Ring {
	r := &Ring{
		store: rbt.NewTree(),
		nodeMap: make(map[string]*node),
		hashfn: newXXHash(),
	}
	return r
}

func NewRing(nodes []string, cnf *Config) *Ring {
	r := &Ring{
		store: rbt.NewTree(),
		nodeMap: make(map[string]*node),
		virtualNodes: cnf.VirtualNodes,
		loadFactor: cnf.LoadFactor,
		hashfn: newXXHash(),
	}
	if r.loadFactor <= 0 {
		r.loadFactor = 1
	}

	for _, node := range nodes {
		r.Add(node)
	}

	return r
}

func (r *Ring) hash(val string) int64 {
	return r.hashfn.hash(val)
}

func (r *Ring) Add(node string) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _,ok := r.nodeMap[node]; ok{
		return
	}
	r.nodeMap[node] = newNode(node)
	hashKey := r.hash(node)
	r.store.Insert(hashKey, node)

	for i := 0; i < r.virtualNodes; i++{
		vNodeKey :=fmt.Sprintf("%s-%d",node,i)
		r.nodeMap[vNodeKey] = newNode(vNodeKey)
		hashKey := r.hash(vNodeKey)
		r.store.Insert(hashKey, node)
	}
}

func (r *Ring) Remove(node string){
	r.mu.Lock()
	defer r.mu.Unlock()

	if _,ok := r.nodeMap[node]; !ok {
		return
	}

	hashKey := r.hash(node)
	r.store.Delete(hashKey)

	for i := 0; i < r.virtualNodes; i++{
		vNodeKey := fmt.Sprintf("%s-%d", node, i)
		hashKey := r.hash(vNodeKey)
		r.store.Delete(hashKey)
		delete(r.nodeMap, vNodeKey)
	}
	delete(r.nodeMap, node)
}

func (r *Ring) Get(key string) (string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	if r.store.Size() == 0 {
		return "", ERR_EMPTY_RING
	}

	var q *rbt.Node
	hashKey := r.hash(key)
	q = r.store.Nearest(hashKey)

	var count int
	for {
		// Avoid infinite recursion if all the servers are under heavy load
		// and if user forgot to call the Done() method
		if count >= r.store.Size() {
			return "", ERR_HEAVY_LOAD
		}

		if hashKey > q.GetKey(){
			g := rbt.FindSuccessor(q)
			if g != nil {
				q =g
			} else {
				q = r.store.Root()
			}
		}

		if r.loadOK(q.GetValue()) {
			break
		}
		count++

		h := rbt.FindSuccessor(q)
		if h == nil {
			q =r.store.Minimum()
		} else {
			q = h
		}
	}
	atomic.AddInt64(&r.nodeMap[q.GetValue()].load, 1)
	atomic.AddInt64(&r.totalLoad,1)
	return q.GetValue(), nil
}

func (r *Ring) loadOK(node string) bool {
	// a safety check if someone performed r.Done more than needed
	if r.totalLoad < 0 {
		r.totalLoad = 0
	}

	var avgLoadPerNode float64
	avgLoadPerNode = float64(int(r.totalLoad+1) / (len(r.nodeMap)))
	if avgLoadPerNode == 0 {
		avgLoadPerNode = 1
	}
	avgLoadPerNode = math.Ceil(avgLoadPerNode * r.loadFactor)

	vnode, ok := r.nodeMap[node]
	if !ok {
		fmt.Printf("given host(%s) not in loadsMap\n", vnode.name)
		return false
	}

	if float64(vnode.load)+1 <= avgLoadPerNode {
		return true
	}

	return false
}

func (r *Ring) Done(node string) {
	if _, ok := r.nodeMap[node]; !ok {
		return
	}
	atomic.AddInt64(&r.nodeMap[node].load, -1)
	atomic.AddInt64(&r.totalLoad, -1)
}