package tree

import (
	"fmt"
	"strings"
)

type Node struct {
	element interface{}
	parent  *Node
	left    *Node
	right   *Node
}

// 判断是左子树还是右子树
type comparer func(interface{}, interface{}) bool

type Bst struct {
	compare comparer
	root    *Node
}

// 创建一个空的二叉树
func (b *Bst) New(compare comparer) *Bst {
	return &Bst{compare: compare}
}

//中序遍历
func inorder_tree_walk(tree *Node) {
	if tree == nil {
		return
	}

	inorder_tree_walk(tree.left)
	fmt.Println(tree.element)
	inorder_tree_walk(tree.right)
}

// 查找一个具有给定关键字的节点,运行时间最坏情况为这颗树的高h, O(h)
func (b *Bst) tree_search(tree *Node, element interface{}) *Node {
	if tree == nil {
		return nil
	}

	if tree.element == element {
		return tree
	}

	// 为真则走左，为假则走右
	if b.compare(element, tree.element) {
		return b.tree_search(tree.right, element)
	} else {
		return b.tree_search(tree.left, element)
	}
}

// 迭代法，避免递归的性能消耗， 你懂的
func (b *Bst) iterative_tree_search(tree *Node, element interface{}) *Node {
	if tree == nil {
		return nil
	}
	for element != tree.element {
		if b.compare(element, tree.element) {
			tree = tree.right
		} else {
			tree = tree.left
		}
	}
	return tree
}

func (b *Bst) tree_minimum(tree *Node) *Node {
	for tree != nil {
		tree = tree.left
	}
	return tree
}

// 最大关键元素, 显然最大关键元素和最小关键元素都是由二叉树的性质决定的，并且它们的代码也是对称的。
func (b *Bst) tree_maximum(tree *Node) *Node {
	for tree != nil {
		tree = tree.right
	}
	return tree
}

// 后继
func (b *Bst) tree_successor(tree *Node) *Node {
	if tree.right != nil {
		return b.tree_minimum(tree.right)
	}
	y := tree.parent
	for y != nil && tree == y.right {
		tree = y
		y = y.parent
	}
	return y
}

// 前继
func (b *Bst) tree_predecessor(tree *Node) *Node {
	if tree.left != nil {
		return b.tree_maximum(tree.left)
	}
	y := tree.parent
	for y != nil && tree == y.right {
		tree = y
		y = y.parent
	}
	return y
}

// 插入
func (b *Bst) tree_insert(tree *Node, element interface{}) *Node {
	y := &Node{}
	for tree != nil {
		// y为双亲节点
		y = tree
		// 举例子n1 < n2 为真则走左， 为假则走右
		if b.compare(element, tree.element) {
			tree = tree.left
		} else {
			tree = tree.right
		}
	}
	var node *Node = &Node{element, nil, nil}
	node.parent = y
	// 说明是一颗空二叉搜索树
	if y == nil {
		b.root = node
		return node
	}
	// 寻找左边还是右边需要插入
	if b.compare(element, y.element) {
		y.left = node
	} else {
		y.right = node
	}
	return y
}
