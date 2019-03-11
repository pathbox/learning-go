package rbt

import (
	"fmt"
	"math"
)

/*
 A binary search tree is a red-black tree if it satisfies the following
 red-black properties:
1. Every node is either red or black.
2. Every leaf (NIL) is black.
3. If a node is red, then both its children are black.
4. Every simple path from a node to a descendant leaf contains the same
   number of black nodes.
*/

const (
	INT_MAX = math.MaxInt64
)

func abs(i int64) int64 {
	if i < 0 {
		return -i
	}
	return i
}

type Color int

const (
	RED   = Color(0)
	BLACK = Color(1)
)

func NewTree() *Tree {
	return &Tree{}
}

type Tree struct {
	root *Node
	size int
}

type Node struct {
	key    int64
	value  string
	color  Color
	parent *Node
	left   *Node
	right  *Node
}

func (t *Tree) Root() *Node {
	return t.root
}

// Inorder traversal
func (t *Tree) Traverse() {
	fn := func(n *Node) {
		fmt.Println(n.key)
	}
	t.root.traverse(fn)
}

func (t *Tree) Insert(key int64, value string) {
	x := newNode(key, value)
	// normal BST insertion
	t.insert(x)
}

func (t *Tree) insert(item *Node) {
	var y *Node
	x := t.root

	for x != nil { //  不断循环，直到最右子树或最左子树，最右子树或最左子树的左右是nil，所以就停止循环了
		y = x
		if item.key < x.key { // 小往左走，大往右走
			// insert value into the left node
			x = x.left
		} else if item.key > x.key {
			// insert value into the right node
			x = x.right
		} else {
			// value exists
			return
		}
	}

	t.size++
	item.parent = y // parent
	item.color = RED

	if y == nil {
		item.color = BLACK
		t.root = item
	} else if item.key < y.key {
		y.left = item
	} else {
		y.right = item
	}

	// Checking RBT conditions and repairing the node
	t.insertRepairNode(item)
}

/*
What happens next depends on the color of other nearby nodes. There are several cases of red–black tree insertion to handle:
    N is the root node, i.e., first node of red–black tree
    N's parent (P) is black
    P is red (so it can't be the root of the tree) and N's uncle (U) is red
    P is red and U is black
*/
func (t *Tree) insertRepairNode(x *Node) {
	// N's parent (P) is not black
	var y *Node
	for x != t.root && x.parent.color == RED {
		if x.parent == x.grandparent().left {
			y = x.grandparent().right
			if y != nil && y.color == RED {
				x.parent.color = BLACK
				y.color = BLACK
				x.grandparent().color = RED
				x = x.grandparent()
			} else {
				if x == x.parent.right {
					x = x.parent
					t.leftRotate(x)
				}
				x.parent.color = BLACK
				x.grandparent().color = RED
				t.rightRotate(x.grandparent())
			}
		} else {
			y = x.grandparent().left
			if y != nil && y.color == RED {
				x.parent.color = BLACK
				y.color = BLACK
				x.grandparent().color = RED
				x = x.grandparent()
			} else {
				if x == x.parent.left {
					x = x.parent
					t.rightRotate(x)
				}
				x.parent.color = BLACK
				x.grandparent().color = RED
				t.leftRotate(x.grandparent())
			}
		}
	}
	// N is the root node, i.e., first node of red–black tree
	t.root.color = BLACK

}
func (t *Tree) leftRotate(x *Node) {
	// Default node inserted will be a red node
	y := x.right
	x.right = y.left
	if y.left != nil {
		y.left.parent = x
	}
	y.parent = x.parent

	// thie is root
	if x.parent == nil {
		t.root = y
	} else {
		if x == x.parent.left {
			x.parent.left = y
		} else {
			x.parent.right = y
		}
	}
	y.left = x
	x.parent = y
}

func (t *Tree) rightRotate(x *Node) {
	y := x.left
	x.left = y.right
	if y.right != nil {
		y.right.parent = x
	}
	y.parent = x.parent

	// this is root
	if x.parent == nil {
		t.root = y
	} else {
		if x == x.parent.right {
			x.parent.right = y
		} else {
			x.parent.left = y
		}
	}
	y.right = x
	x.parent = y

}

func (t *Tree) replace(a, b *Node) {
	if a.parent == nil {
		t.root = b
	} else if a == a.parent.left {
		a.parent.left = b
	} else {
		a.parent.right = b
	}
	if b != nil {
		b.parent = a.parent
	}
}

func (t *Tree) Search(key int64) *Node {
	x := t.root

	if x == nil {
		return nil
	}

	for x != nil {
		switch {
		case key == x.key:
			return x
		case key < x.key:
			x = x.left
		case key > x.key:
			x = x.right
		}
	}

	return nil
}

func (t *Tree) Delete(key int64) {
	z := t.Search(key)
	if z == nil {
		return
	}
	t.delete(z)
}

func (t *Tree) delete(z *Node) *Node {
	// fmt.Printf("del: %+v\n", z)

	var x, y *Node
	y = z

	if z.left == nil {
		x = z.right
		t.replace(z, z.right)
	} else if z.right == nil {
		x = z.left
		t.replace(z, z.left)

	} else {
		y = z.successor()
		if y.left != nil {
			x = y.left
		} else {
			x = y.right
		}
		x.parent = y.parent

		if y.parent == nil {
			t.root = x
		} else {
			if y == y.parent.left {
				y.parent.left = x
			} else {
				y.parent.right = x
			}
		}
	}

	if y.color == BLACK {
		t.deleteRepairNode(x)
	}
	t.size--

	return y
}

func (t *Tree) deleteRepairNode(x *Node) {
	if x == nil {
		return
	}
	var w *Node
	for x != t.root && x.color == BLACK {
		if x == x.parent.left {
			w = x.sibling()
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.leftRotate(x.parent)
				w = x.parent.right
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.right.color == BLACK {
					w.left.color = BLACK
					w.color = RED
					t.rightRotate(w)
					w = x.parent.right
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.right.color = BLACK
				t.leftRotate(x.parent)
				x = t.root
			}
		} else {
			w = x.sibling()
			if w.color == RED {
				w.color = BLACK
				x.parent.color = RED
				t.rightRotate(x.parent)
				w = x.parent.left
			}
			if w.left.color == BLACK && w.right.color == BLACK {
				w.color = RED
				x = x.parent
			} else {
				if w.left.color == BLACK {
					w.right.color = BLACK
					w.color = RED
					t.leftRotate(w)
					w = x.parent.left
				}
				w.color = x.parent.color
				x.parent.color = BLACK
				w.left.color = BLACK
				t.rightRotate(x.parent)
				x = t.root
			}

		}
	}
	x.color = BLACK
}

// Preorder prints the tree in pre order
func (t *Tree) Preorder() {
	fmt.Println("preorder begin!")
	if t.root != nil {
		t.root.preorder()
	}
	fmt.Println("preorder end!")
}

func (t *Tree) Size() int {
	return t.size
}

func (t *Tree) Nearest(key int64) *Node {
	return nearestNode(t.root, key)
}

func (t *Tree) Minimum() *Node {
	if t.root != nil {
		return t.root.minimum()
	}
	return nil
}

func nearestNode(root *Node, key int64) *Node {
	if root == nil {
		return nil
	}
	var minDiff int64
	var minDiffKey *Node

	minDiff = INT_MAX
	n := root
	for n != nil {
		if n.key == key {
			minDiffKey = n
			return minDiffKey
		}
		newDiff := abs(n.key - key)
		if minDiff > newDiff {
			minDiff = newDiff
			minDiffKey = n
		}
		if key < n.key {
			n = n.left
		} else {
			n = n.right
		}
	}

	return minDiffKey
}

func newNode(key int64, value string) *Node {
	return &Node{
		key:   key,
		value: value,
	}
}

type Node struct {
	key    int64
	value  string
	color  Color
	parent *Node
	left   *Node
	right  *Node
}

func (n *Node) GetKey() int64 {
	return n.key
}

func (n *Node) GetValue() string {
	return n.value
}

func (n *Node) father() *Node {
	return n.parent
}

func (n *Node) grandparent() *Node {
	g := n.father()
	// No father means no granparent
	if g == nil {
		return nil
	}
	return g.parent
}

func (n *Node) sibling() *Node {
	p := n.father()
	// No parent means no brother
	if p == nil {
		return nil
	}
	if n == p.left {
		return p.right
	}
	return p.left
}

func (n *Node) uncle() *Node {
	p := n.father()
	g := n.grandparent()
	// No grandparent means no uncle
	if g == nil {
		return nil
	}
	return p.sibling()
}

func (n *Node) successor() *Node {
	if n.right != nil {
		return n.right.minimum()
	}
	y := n.parent
	for y != nil && n == y.right {
		n = y
		y = y.parent
	}
	return y
}

func (n *Node) predecessor() *Node {
	if n.left != nil {
		return n.left.maximum()
	}
	y := n.parent
	for y != nil && n == y.left {
		n = y
		y = y.parent
	}
	return y
}

func (n *Node) minimum() *Node {
	for n.left != nil {
		n = n.left
	}
	return n
}

func (n *Node) maximum() *Node {
	for n.right != nil {
		n = n.right
	}
	return n
}

func (n *Node) traverse(fn func(*Node)) {
	if n == nil {
		return
	}
	n.left.traverse(fn)
	fn(n)
	n.right.traverse(fn)
}

func (n *Node) preorder() {
	fmt.Printf("(%v %v)", n.key, n.value)
	if n.parent == nil {
		fmt.Printf("nil")
	} else {
		fmt.Printf("whose parent is %v", n.parent.key)
	}
	if n.color == RED {
		fmt.Println(" and color RED")
	} else {
		fmt.Println(" and color BLACK")
	}
	if n.left != nil {
		fmt.Printf("%v's left child is ", n.key)
		n.left.preorder()
	}
	if n.right != nil {
		fmt.Printf("%v's right child is ", n.key)
		n.right.preorder()
	}
}

func FindSuccessor(n *Node) *Node {
	if n.right != nil {
		return n.right.minimum()
	}
	y := n.parent
	for y != nil && n == y.right {
		n = y
		y = y.parent
	}
	return y
}

func FindPredecessor(n *Node) *Node {
	if n.left != nil {
		return n.left.maximum()
	}
	y := n.parent
	for y != nil && n == y.left {
		n = y
		y = y.parent
	}
	return y
}
