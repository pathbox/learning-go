package main

func main() {

}

func (mq *MyQueue) Pop() interface{} {
	old := *mq
	n := len(old)
	item := old[n-1] // last item
	item.index = -1
	*mq = old[0 : n-1] // 将mq 赋值为 除去最后一个item的队列
	return item        // 返回最后一个item
}

func (mq *MyQueue) Push(x interface{}) {
	n := len(*mq)
	item := x.(*Item)
	item.index = n
	*mq = append(*mq, item)
}

type MyQueue []*Item

type Item struct {
	Name  string
	Value interface{}
	index int
}
