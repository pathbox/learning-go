package main

import (
	"container/heap"
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Time    time.Time
	Comment string
}

type Tasks []Task

func (self Tasks) Len() int { return len(self) }

func (self Tasks) Less(i, j int) bool {
	return self[i].Time.Sub(self[j].Time) < 0
}

func (self Tasks) Swap(i, j int) {
	self[i], self[j] = self[j], self[i]
}

func (self *Tasks) Push(x interface{}) {
	*self = append(*self, x.(Task))
}

func (self *Tasks) Pop() interface{} {
	old := *self
	n := len(old)
	task := old[n-1]
	*self = old[0 : n-1]
	return task
}

type TaskQueue struct {
	lock     *sync.RWMutex
	cond     *sync.Cond
	tasks    Tasks
	timeouts Tasks
	push     chan Task
}

func NewTaskQueue() *TaskQueue {
	self := &TaskQueue{
		tasks: Tasks{},
		lock:  &sync.RWMutex{},
		push:  make(chan Task),
	}
	self.cond = sync.NewCond(self.lock)
	heap.Init(&self.tasks)

	go func() {
		wait := self.adjust()

		for {
			select {
			case task := <-self.push:
				self.lock.Lock()
				heap.Push(&self.tasks, task)
				wait = self.adjust()
				self.lock.Unlock()
				self.cond.Signal()

			case <-time.After(wait):
				self.lock.Lock()
				wait = self.adjust()
				self.lock.Unlock()
				self.cond.Signal()
			}
		}
	}()

	return self
}

func (self *TaskQueue) adjust() time.Duration {
	for {
		if len(self.tasks) == 0 {
			return time.Hour
		}
		task := self.tasks[0]
		diff := task.Time.Sub(time.Now())
		if diff > 0 {
			return diff
		}
		heap.Pop(&self.tasks)
		self.timeouts = append(self.timeouts, task)
	}
}

func (self *TaskQueue) Push(task Task) {
	self.push <- task
}

func (self *TaskQueue) Pop() (task Task) {
	self.cond.L.Lock()
	for {
		if len(self.timeouts) > 0 {
			task = self.timeouts[0]
			self.timeouts = self.timeouts[1:]
			break
		}
		if len(self.tasks) > 0 && time.Now().Sub(self.tasks[0].Time) >= 0 {
			task = heap.Pop(&self.tasks).(Task)
			break
		}
		self.cond.Wait()
	}
	self.cond.L.Unlock()
	return
}

func main() {
	queue := NewTaskQueue()

	queue.Push(Task{Time: time.Now().Add(-time.Second * 3), Comment: "-3"})
	queue.Push(Task{Time: time.Now().Add(time.Second * 3), Comment: "3"})
	queue.Push(Task{Time: time.Now().Add(time.Second * 2), Comment: "2"})
	queue.Push(Task{Time: time.Now().Add(time.Second), Comment: "1"})

	fmt.Println("start pop")

	go queue.Push(Task{Time: time.Now().Add(time.Millisecond * 500), Comment: "500ms"})

	for {
		task := queue.Pop()
		fmt.Println(task.Comment)
	}
}
