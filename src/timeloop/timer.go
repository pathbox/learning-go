package timer

import (
	"container/heap"
	"log"
	"runtime/debug"
	"sync"
	"time"
)

const MIN_TIMER = 20 * time.Millisecond

var (
	logger    *log.Logger
	IsRunning bool = true
)

/*
	map + 最小堆，最小堆用来排序，使用map 可以用来逻辑del，通过将elem 时间戳设置成当前时间
	回调逻辑因为可能阻塞主任务，采用异步方式入队列，队列满则写入磁盘，或者直接drop
*/

type TimerHander interface {
	AddFuncWithId(d time.Duration, taskid string, callBack func())
	StartTimerLoop()
	GetLen()
	ExitLoop()
	RemoveTask(taskid string)
	StartAsyncWorker(workNum int) //起一组协程池消费回调
}

type TimerEntry struct {
	runTime  time.Time    // 到期时间
	callback CallbackFunc // 回调方法
	taskId   string       // 业务ID
}

// time heap
type TimerHeapHandler struct {
	timers            []*TimerEntry
	locker            *sync.Mutex
	TaskTimerEntryMap map[string]*TimerEntry
	TaskQueue         chan func()
}

func New(workerNum int, qsize int) *TimerHeapHandler {
	var ttimers TimerHeapHandler
	heap.Init(&ttimers)

	h := TimerHeapHandler{
		locker:            new(sync.Mutex),
		TaskTimerEntryMap: make(map[string]*TimerEntry),
		TaskQueue:         make(chan func(), qsize),
	}
	h.startAsyncWorker(workerNum)
	return &h
}

// 实现heap
func (h *TimerHeapHandler) Len() int {
	return len(h.timers)
}

func (h *TimerHeapHandler) Less(i, j int) bool {
	t1, t2 := h.timers[i].runTime, h.timers[j].runTime
	if t1.Before(t2) {
		return true
	}
	return false
}

func (h *TimerHeapHandler) Swap(i, j int) {
	var tmp *TimerEntry
	tmp = h.timers[i]
	h.timers[i] = h.timers[j]
	h.timers[j] = tmp
}

func (h *TimerHeapHandler) Push(x interface{}) {
	h.timers = append(h.timers, x.(*TimerEntry))
}

func (h *TimerHeapHandler) Pop() (ret interface{}) {
	l := len(h.timers)
	h.timers, ret = h.timers[:l-1], h.timers[l-1]
	return
}

func (h *TimerHeapHandler) GetLength() map[string]int {
	return map[string]int{
		"heap_len": h.Len(),
		"id_map":   len(h.TaskTimerEntryMap),
	}
}

func (h *TimerHeapHandler) RemoveById(id string) {
	h.locker.Lock()
	if entry, ok := h.TaskTimerEntryMap[id]; ok {
		entry.runTime = time.Now() //将time 设置成now， 排序立刻会被放到最前或者最后
	}
	h.locker.Unlock()
}

func (h *TimerHeapHandler) AddFuncWithId(d time.Duration, taskId string, callback CallbackFunc) *TimerEntry {
	return h.addCallback(d,
		callback,
		taskId)
}

func (h *TimerHeapHandler) AddCronFuncWithId(d time.Duration, taskId string, callback CallbackFunc) *TimerEntry {
	return h.addCallback(d,
		callback,
		taskId)
}

func (h *TimerHeapHandler) addCallback(d time.Duration, callback CallbackFunc, taskId string) *TimerEntry {
	if d < MIN_TIMER {
		d = MIN_TIMER
	}

	t := &TimerEntry{
		runTime:  time.Now().Add(d),
		taskId:   taskId,
		callback: callback,
	}

	h.locker.Lock()
	heap.Push(h, t)
	h.TaskTimerEntryMap[taskId] = t
	h.locker.Unlock()
	return t
}

//every loop get all task expired
func (h *TimerHeapHandler) EventLoop() {
	now := time.Now()
	h.locker.Lock()
	for { // for循环不断调用
		if h.Len() <= 0 {
			break
		}

		// 使用小顶堆,先判断堆顶的时间是否到期了，到期了才会从堆顶拿出改节点数据
		nextRunTime := h.timers[0].runTime
		if nextRunTime.After(now) {
			// not due time
			break
		}

		t := heap.Pop(h).(*TimerEntry)
		if t.taskId != "" {
			delete(h.TaskTimerEntryMap, t.taskId)
		}

		callback := t.callback
		if callback == nil {
			continue
		}
		// h.runCallback(callback) don't use goroutine
		//use goroutine if full write to file or drop !!!
		select {
		case h.TaskQueue <- callback: // 将注册登录callback 传给TackQueue
		default:
			log.Printf("taskId %v droped", t.taskId)
		}
	}
	h.locker.Unlock()
}

func (h *TimerHeapHandler) StartTimerLoop(tickInterval time.Duration) {
	go func() {
		for IsRunning {
			time.Sleep(tickInterval)
			h.EventLoop()
		}
	}()
}

func (h *TimerHeapHandler) runCallback(callback CallbackFunc) {
	defer func() {
		err := recover()
		if err != nil {
			if logger != nil {
				logger.Printf("callback %v paniced: %v\n", callback, err)
			}
			debug.PrintStack()
		}
	}()

	if callback == nil {
		return
	}
	callback()
}

func (h *TimerHeapHandler) Exit() {
	IsRunning = false
}

func (h *TimerHeapHandler) startAsyncWorker(num int) {
	for index := 0; index < num; index++ {
		go h.asyncWorker()
	}
}

func (h *TimerHeapHandler) asyncWorker() {
	for IsRunning {
		select {
		case call := <-h.TaskQueue:
			h.runCallback(call)
		}
	}
}

type CallbackFunc func()
