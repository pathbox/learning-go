package timewheel

import (
	"math/rand"
	"time"
)

type ConcurrentTimer struct {
	interval     time.Duration
	loopInterval time.Duration
	slotNum      int
	slots        []*Timer
}

func NewConcurrentTimer(slotNum int, loopInterval time.Duration) (*ConcurrentTimer, error) {
	if slotNum <= 0 {
		return nil, ErrLtMinDelay
	}

	var n = &ConcurrentTimer{
		interval:     time.Second,
		slotNum:      slotNum,
		loopInterval: loopInterval,
		slots:        make([]*Timer, slotNum),
	}

	n.init()
	return n, nil
}

func (ct *ConcurrentTimer) init() {
	for idx, _ := range ct.slots {
		var tm = NewTimer()
		tm.loopInterval = ct.loopInterval
		ct.slots[idx] = tm
	}
}

func (ct *ConcurrentTimer) GetOneTimer() *Timer {
	pos := rand.Intn(ct.slotNum - 1)
	return ct.slots[pos]
}

func (ct *ConcurrentTimer) Start() {
	for _, tm := range ct.slots {
		tm.Start()
	}
}

func (ct *ConcurrentTimer) Stop() {
	for _, tm := range ct.slots {
		tm.Stop()
	}
}

type ConcurrentTimerEntry struct {
	timer *Timer
	C     chan time.Time
	event *Event
}

func (ce *ConcurrentTimerEntry) Stop() {
	ce.timer.Del(ce.event)
}

func (ce *ConcurrentTimerEntry) Reset(delay time.Duration) {
	ce.timer.Set(ce.event, delay)
}