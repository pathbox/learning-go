package simpletimewheel

import (
	"sync"
	"time"
)

type slot struct {
	id int
	elements map[interface{}]interface{} // 底层是一个map, 在论文中时间轮的slot是一个链表结构
}

type element struct {
	v             interface{}
	remainingTime int
}

func newSlot(id int) *slot {
	s := &slot{id: id}
	s.elements = make(map[interface{}]interface{})
	return s
}

func (s *slot) add(c interface{}) {
	s.elements[c] = c
}

func (s *slot) remove(c interface{}) {
	delete(s.elements, c)
}

type handler func(interface{})

func TimeWheel struct {
	tickDuration time.Duration
	ticksPerWheel int
	currentTickIndex int
	ticker *time.Ticker
	onTick handler
	wheel []*slot
	indicator map[interface{}]*slot
	sync.RWMutex

	taskChan chan interface{}
	quitChan chan interface{}
}

func New(tickDuration time.Duration, ticksPerWheel int, f handler) *TimeWheel {
	if tickDuration < 1 || ticksPerWheel < 1 || nil == f {
		return nil
	}

	ticksPerWheel++
	t := &TimeWheel{
		tickDuration:     tickDuration,
		ticksPerWheel:    ticksPerWheel, // 时间轮的刻度
		onTick:           f,// 回调函数
		currentTickIndex: 0,
		taskChan:         make(chan interface{}),
		quitChan:         make(chan interface{}),
	}
	t.indicator = make(map[interface{}]*slot, 0)

	t.wheel = make([]*slot, ticksPerWheel)
	for i := 0; i < ticksPerWheel; i++ {
		t.wheel[i] = newSlot(i)
	}

	return t
}


func (t *TimeWheel) Start() {
	t.ticker = time.NewTicker(t.tickDuration) // 定时器
	go t.run()
}

func (t *TimeWheel) Add(c interface{}) {
	t.taskChan <- &element{c,t.ticksPerWheel-1} // t.ticksPerWheel-1 相当于是尾部
}

func (t *TimeWheel) AddWithRemainingTime(c interface{}, remainingTime int) {
	t.taskChan <- &element{c,remainingTime}
}

func (t *TimeWheel) Remove(c interface{}) {
	if v, ok := t.indicator[c]; ok {
		v.remove(c)
	}
}

func (t *TimeWheel) getCurrentTickIndex() int {
	t.RLock()
	defer t.RUnlock()
	return t.currentTickIndex
}

func (t *TimeWheel) Stop() {
	close(t.quitChan)
}

func (t *TimeWheel) run() {
	for {
		select {
		case <-t.quitChan:
			t.ticker.Stop()
			break
		case <- t.ticker.C: // ⏲  每个tick去尝试取出slot的元素进行函数执行
			if t.ticksPerWheel == t.currentTickIndex { // 循环的方式，另一种是取余的方式
				t.currentTickIndex = 0
			}

			slot := t.wheel[t.currentTickIndex] // 取出当前轮询到的slot
			for _, v := range slot.elements { // 遍历slot中的所有元素, 取出该元素，
				slot.remove(v)
				delete(t.indicator, v)
				t.onTick(v) // 将元素参数代入，执行回调函数
			}

			t.currentTickIndex++ // 索引继续向前
		case v := <-t.taskChan:
			element, ok := v.(*element)
			if !ok {
				return
			}
			t.Remove(element.v)
			// 将新的任务需要的参数加入到 wheel
			elementIdx := t.getCurrentTickIndex()+element.remainingTime
			if elementIdx > t.ticksPerWheel {
				elementIdx = elementIdx - t.ticksPerWheel
			}
			elementIdx--
			slot := t.wheel[elementIdx]
			slot.add(element.v)
			t.indicator[element.v] = slot
		}
	}
}