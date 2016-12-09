package main

import (
	"container/heap"
	"fmt"
	"time"
)

type Task struct {
	Time    time.Time
	Comment string
}

type TaskQueue []Task

func (self TaskQueue) Len() int {
	return len(self)
}

func (self TaskQueue) Less(i, j int) bool {
	return self[i].Time.Sub(self[j].Time) < 0
}

func (self TaskQueue) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self *TaskQueue) Push(x interface{}) {
	*self = append(*self, x.(Task))
}

func (self *TaskQueue) Pop() interface{} {
	old := *self
	n := len(old)
	task := old[n-1]
	*self = old[0 : n-1]
	return task
}

func main() {
	queue := TaskQueue{}
	heap.Init(&queue)

	heap.Push(&queue, Task{Time: time.Now().Add(time.Second * 3), Comment: "3"})
	heap.Push(&queue, Task{Time: time.Now().Add(time.Second * 2), Comment: "2"})
	heap.Push(&queue, Task{Time: time.Now().Add(time.Second), Comment: "1"})

	for queue.Len() > 0 {
		task := heap.Pop(&queue).(Task)
		diff := task.Time.Sub(time.Now())
		if diff > 0 {
			time.Sleep(diff)
		}
		fmt.Println(task.Comment)
	}
}
