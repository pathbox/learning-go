package tree

import (
	"container/list"
	"fmt"
)

type Node struct {
	Val   int
	Left  *Node
	Right *Node
}

// 栈的定义
type Stack struct {
	list *list.List
}

// 利用list 模拟栈
func NewStack() *Stack {
	list := list.New()
	return &Stack{list}
}

func (stack *Stack) Push(value interface{}) {
	stack.list.PushBack(value)
}

func (stack *Stack) Pop() interface{} {
	if e := stack.list.Back(); e != nil { // 判断最后一个值存不存在
		stack.list.Remove(e)
		return e.Value
	}

	return nil
}

func (stack *Stack) Len() int {
	return stack.list.Len()
}

func (stack *Stack) Empty() bool {
	return stack.Len() == 0
}

// 前序遍历
func (root *Node) PreTravesal() {
	if root == nil {
		return
	}

	s := NewStack()
	s.Push(root)

	for !s.Empty() {
		cur := s.Pop().(*Node)
		fmt.Println(cur.Val)

		if cur.Right != nil {
			s.Push(cur.Right)
		}
		if cur.Left != nil {
			s.Push(cur.Left)
		}
	}
}

// 中序遍历
func (root *Node) InTravesal() {
	if root == nil {
		return
	}

	s := NewStack()
	cur := root
	for {
		for cur != nil {
			s.Push(cur)
			cur = cur.Left
		}

		if s.Empty() {
			break
		}

		cur = s.Pop().(*Node)
		fmt.Println(cur.Val)
		cur = cur.Right
	}
}

// 后序遍历
func (root *Node) PostTravesal() {
	if root == nil {
		return
	}

	s := NewStack()
	out := NewStack()
	s.Push(root)

	for !s.Empty() {
		cur := s.Pop().(*Node)
		out.Push(cur)

		if cur.Left != nil {
			s.Push(cur.Left)
		}

		if cur.Right != nil {
			s.Push(cur.Right)
		}
	}

	for !out.Empty() {
		cur := out.Pop().(*Node)
		fmt.Println(cur.Val)
	}
}

type Queue interface {
	Offer(e interface{})
	Poll() interface{}
	Clear() bool
	Size() int
	IsEmpty() bool
}

type LinkedList struct {
	elements []interface{}
}

func NewQueue() *LinkedList {
	return &LinkedList{}
}

// 实现队列
func (queue *LinkedList) Offer(e interface{}) {
	queue.elements = append(queue.elements, e)
}

func (queue *LinkedList) Poll() interface{} {
	if queue.IsEmpty() {
		fmt.Println("Poll error : queue is Empty")
		return nil
	}

	firstElement := queue.elements[0]
	queue.elements = queue.elements[1:]
	return firstElement
}

func (queue *LinkedList) Size() int {
	return len(queue.elements)
}

func (queue *LinkedList) IsEmpty() bool {
	return len(queue.elements) == 0
}

func (queue *LinkedList) Clear() bool {
	if queue.IsEmpty() {
		fmt.Println("queue is Empty!")
		return false
	}
	for i := 0; i < queue.Size(); i++ {
		queue.elements[i] = nil
	}
	queue.elements = nil
	return true
}

// 层序遍历
func (root *Node) LevelTravesal() {
	if root == nil {
		return
	}

	linkedList := NewQueue()
	linkedList.Offer(root)

	for !linkedList.IsEmpty() {
		cur := linkedList.Poll().(*Node)
		fmt.Println(cur.Val)

		if cur.Left != nil {
			linkedList.Offer(cur.Left)
		}

		if cur.Right != nil {
			linkedList.Offer(cur.Right)
		}
	}
}
