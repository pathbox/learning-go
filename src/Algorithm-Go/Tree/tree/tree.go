package tree

import (
	"fmt"
	"math/rand"
)

type TNode struct {
	Left  *TNode
	Value int
	Right *TNode
}

// New returns a new, random binary tree holding the values k, 2k, ..., 10k
func New(k int) *TNode {
	var t *TNode
	for _, v := range rand.Perm(10) {
		t = insert(t, (1+v)*k) // 第一次得到一个node 只有value，没有Left和Right。第二次循环，上一次得到的node会被代入到insert中，得到一个新的node。新的node会和上一个node有连接，之后继续代入到insert
	} // 代入insert 能够知道value值，不知道Left和Right，通过和传入的t进行比较得到Left和Right关系
	return t
}

// 第一个node，传入的t参数是nil，所以返回一个没有Left和Right的node
// 之后传入node和v，会进行value的比较
func insert(t *TNode, v int) *TNode {
	if t == nil { // 递归结束，返回一个末尾node，这个node属于某个node的Left或Right
		return &TNode{nil, v, nil}
	}
	if v < t.Value { // 如果node值更大，v的node插入到node.Left。node的左子树和新的node继续进行比较，node会在左子树中找到一个合适的位置，插入
		t.Left = insert(t.Left, v)
	} else { // 如果node值更小，v的node插入到node.Right。node的右子树和新的node继续进行比较，node会在右子树中找到一个合适的位置，插入
		t.Right = insert(t.Right, v)
	}

	return t
}

func (t *TNode) String() string {
	if t == nil {
		return "()"
	}
	s := ""
	if t.Left != nil {
		s += t.Left.String() + " "
	}
	s += fmt.Sprint(t.Value)
	if t.Right != nil {
		s += " " + t.Right.String()
	}
	return "(" + s + ")"
}
