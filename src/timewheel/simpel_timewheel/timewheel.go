package simpletimewheel

import (
	"sync"
	"time"
)

type slot struct {
	id int
	elements map[interface{}]interface{} // 底层是一个map
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
		ticksPerWheel:    ticksPerWheel,
		onTick:           f,
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
	t.ticker = time.NewTicker(t.tickDuration)
	go t.run()
}

func (t *TimeWheel) Add(c interface{}) {
	t.taskChan <- &element{c,t.ticksPerWheel-1}
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
		case <- t.ticker.C:
			if t.ticksPerWheel == t.currentTickIndex {
				t.currentTickIndex = 0
			}

			slot := t.wheel[t.currentTickIndex]
			for _, v := range slot.elements {
				slot.remove(v)
				delete(t.indicator, v)
				t.onTickv
			}

			t.currentTickIndex++
		case v := <-t.taskChan:
			element, ok := v.(*element)
			if !ok {
				return
			}
			t.Remove(element.v)

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