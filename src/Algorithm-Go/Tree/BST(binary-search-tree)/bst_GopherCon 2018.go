package bst

type Node struct {
	Key   int
	Left  *Node
	Right *Node
}

var tree = &Node{Key: 6}

func (n *Node) Search(key int) bool {
	// This is our base case. If n == nil, `key`
	// doesn't exist in our binary search tree.
	if n == nil {
		return false
	}

	if n.Key < key { // move right
		return n.Right.Search(key)
	} else if n.Key > key { // move left
		return n.Left.Search(key)
	}

	// n.Key == key, we found it!
	return true
}

func (n *Node) Insert(key int) {
	if n.Key < key {
		if n.Right == nil { // we found an empty spot, done!
			n.Right = &Node{Key: key}
		} else { // look right
			n.Right.Insert(key)
		}
	} else if n.Key > key {
		if n.Left == nil { // we found an empty spot, done!
			n.Left = &Node{Key: key}
		} else { // look left
			n.Left.Insert(key)
		}
	}
	// n.Key == key, don't need to do anything
}

/*
要保证整个性质，我们必须在删除的位置上，找一个合适的值来进行替换，使得BST上的每个节点都满足 当前节点的值大于左节点但是小于右节点

而替换策略就是：
1、当前删除位置，用左边子树的最大值的节点替换
2、或者是，用右边子树的最小值的节点替换
*/

func (n *Node) Delete(key int) *Node {
	// search for `key`
	if n.Key < key {
		n.Right = n.Right.Delete(key)
	} else if n.Key > key {
		n.Left = n.Left.Delete(key)
	} else { // n.Key == `key` 找到要删除的这个key的节点了
		if n.Left == nil { // just point to opposite node
			return n.Right // 如果key节点的左子树为nil，返回key节点的右子树
		} else if n.Right == nil { // just point to opposite node
			return n.Left // 如果key节点的右子树为nil，返回key节点的左子树
		}

		// if `n` has two children, you need to
		// find the next highest number that
		// should go in `n`'s position so that
		// the BST stays correct
		min := n.Right.Min() // 找到右子树的最小结点

		// we only update `n`'s key with min
		// instead of replacing n with the min
		// node so n's immediate children aren't orphaned
		n.Key = min // 换为当前节点
		n.Right = n.Right.Delete(min)
	}
	return n
}

func (n *Node) Min() int {
	if n.Left == nil {
		return n.Key
	}
	return n.Left.Min()
}

func (n *Node) Max() int {
	if n.Right == nil {
		return n.Key
	}
	return n.Right.Max()
}

// import "testing"

// func BenchmarkSearch(b *testing.B) {
//     tree := &Node{Key: 6}

//     for i := 0; i < b.N; i++ {
//         tree.Search(6)
//     }
// }
