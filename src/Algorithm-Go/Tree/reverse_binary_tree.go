package main

import (
	"fmt"

	// "golang.org/x/tour/tree"
	tree "./tree"
)

//  递归遍历这棵树的每一个node节点，将node的左儿子节点和右儿子节点互换，之后继续递归下去，直到叶子节点或同时没有左儿子节点和右儿子节点的节点位置
func Reverse(root *tree.TNode) {
	if root == nil {
		return
	}

	if root.Left == nil && root.Right == nil { // 判断节点是否同时有左儿子节点和右儿子节点
		return
	}

	root.Left, root.Right = root.Right, root.Left // 将node的左节点和右节点互换位置
	Reverse(root.Left)                            // 递归当前node的左子树
	Reverse(root.Right)                           // 递归当前node的右子树
}

func main() {
	root := tree.New(1)
	fmt.Println("Before reverse", root.String())
	Reverse(root)

	fmt.Println("After reverse", root.String())

}
