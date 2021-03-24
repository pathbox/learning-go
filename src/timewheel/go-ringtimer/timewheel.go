package timewheel

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

const (
	posWriteMode = iota
	posReadMode

	SecondInterval = time.Second
)

var (
	ErrLtMinDelay  = errors.New("lt delay min")
	ErrInvalidSlot = errors.New("invalid slot num")
)

type TimeWheel struct {
	counter int64
	interval time.Duration
	ticker *time.Ticker
	slots []*Timer
	timer      map[interface{}]int
	currentPos int
	slotNum    int
	started    bool

	concurrentTimerStarted int32
	concurrentTimer        *ConcurrentTimer

	ctx    context.Context
	cancel context.CancelFunc

	updateEventChannel chan *Event
}

func NewTimeWheel(interval time.Duration, slotNum int) (*TimeWheel, error) {
	if interval < SecondInterval {
		return nil, ErrLtMinDelay
	}

	if slotNum <= 0 {
		return nil, ErrInvalidSlot
	}

	ccTimer, err := NewConcurrentTimer(5, interval)
	if err != nil {
		return nil, err
	}

	var ctx, cancel = context.WithCancel(context.Background())
	tw := &TimeWheel{
		interval: interval,
		slots: make([]*Timer, slotNum),
		timer: make(map[interface{}]int),
		slotNum: slotNum,
		ctx: ctx,
		cancel: cancel,
		concurrentTimer: ccTimer,
	}

	tx.initSlots()
	return tw, nil
}

func (tw *TimeWheel) initSlots() {
	for i := 0; i < tw.slotNum; i++ {
		tw.slots[i] = NewTimer()
	}
}

func (tw *TImeWheel) Start() {
	if tw.started {
		return
	}

	tw.currentPos = tw.getInitPosition()
	tw.ticker = time.NewTicker(tw.interval / 2)
	tw.started = true
	go tw.start()
}

func (tw *TimeWheel) Stop() {
	tw.cancel()
}

func (tw *TimeWheel) ResetTimer(entry *TimerEntry, delay time.Duration) bool {
	if entry.event == nil {
		return false
	}

	entry.Reset(delay)
	return true
}

func (tw *TimeWheel) AddCronTimer(delay time.Duration, fn ExpireFunc) (*TimerEntry, error) {
	if atomic.CompareAndSwapInt32(&tw.concurrentTimerStarted, 0, 1) {
		tw.concurrentTimer.Start()
	}

	timer := tw.concurrentTimer.GetOneTimer()

	// new TimerEntry
	entry := new(TimerEntry)
	entry.init()

	// new event
	ev := timer.addAny(delay, fn, false, entry.C)

	// link
	entry.event = ev
	entry.timer = timer
	entry.tw = tw

	return entry, nil
}

func (tw *TimeWheel) AddTimer(delay time.Duration, fn ExpireFunc) (*TimerEntry, error) {
	if delay < time.Millisecond {
		return nil, ErrLtMinDelay
	}

	var (
		pos   = tw.getWritePosition(delay)
		timer = tw.slots[pos]
	)

	// new TimerEntry
	entry := new(TimerEntry)
	entry.init()

	// new event
	ev := timer.addAny(delay, fn, false, entry.C)
	ev.slotPos = pos

	// link
	entry.event = ev
	entry.timer = timer
	entry.tw = tw

	return entry, nil
}

func (tw *TimeWheel) RemoveTimer(entry *TimerEntry) {
	if entry.event == nil {
		return
	}

	entry.Stop()
}

func (tw *TimeWheel) Sleep(delay time.Duration) {
	var (
		pos   = tw.getWritePosition(delay)
		timer = tw.slots[pos]
	)

	timer.Sleep(delay)
}

func (tw *TimeWheel) After(delay time.Duration) <-chan time.Time {
	var (
		pos   = tw.getWritePosition(delay)
		timer = tw.slots[pos]
	)

	return timer.After(delay)
}

func (tw *TimeWheel) AfterFunc(delay time.Duration, fn ExpireFunc) (*TimerEntry, error) {
	return tw.AddTimer(delay, fn)
}

func (tw *TimeWheel) GetTimerCount() int64 {
	return atomic.LoadInt64(&tw.counter)
}

func (tw *TimeWheel) start() {
	tw.tickHandler()
	for {
		select {
		case <-tw.ticker.C:
			go tw.tickHandler()

		case <-tw.updateEventChannel:

		case <-tw.ctx.Done():
			return
		}
	}
}

func (tw *TimeWheel) tickHandler() {
	pos := tw.getInitPosition()
	tw.currentPos = pos
	timer := tw.slots[pos]
	timer.LoopOnce()

	// wheel full, reset init 0
	// if tw.currentPos == tw.slotNum-1 {
	// 	tw.currentPos = 0
	// } else {
	// 	tw.currentPos++
	// }
}

func (tw *TimeWheel) GetTimers() []*Timer {
	return tw.slots
}

type TimerStatsRes struct {
	SlotID int
	Len    int
}

func (tw *TimeWheel) GetEachTimerLength() []TimerStatsRes {
	var res = make([]TimerStatsRes, tw.slotNum)
	for idx, tm := range tw.slots {
		res[idx] = TimerStatsRes{
			SlotID: idx,
			Len:    tm.Len(),
		}
	}

	return res
}

func (tw *TimeWheel) GetTimersLength() int {
	var counter int
	for _, tm := range tw.slots {
		counter += tm.Len()
	}

	return counter
}

func (tw *TimeWheel) getPosition(d time.Duration) (pos int) {
	return tw.callGetPosition(d, posReadMode)
}

func (tw *TimeWheel) getWritePosition(d time.Duration) (pos int) {
	return tw.callGetPosition(d, posWriteMode)
}

func (tw *TimeWheel) callGetPosition(delay time.Duration, mode int) int {
	var (
		pos       int
		plus      int
		delayUnit int
	)

	if tw.interval >= time.Millisecond && tw.interval < time.Second {
		delayUnit = int(delay.Nanoseconds() / 1000 / 1000)
		plus = int(time.Now().Unix()) + delayUnit
	} else {
		// defualt second unit
		delayUnit = int(delay.Seconds())
		plus = int(time.Now().Unix()) + delayUnit
	}

	pos = plus % tw.slotNum

	if mode == posWriteMode && pos == tw.currentPos {
		pos++
	}

	return pos
}

func (tw *TimeWheel) getInitPosition() int {
	var pos int
	if tw.interval >= time.Millisecond && tw.interval < time.Second {
		pos = int(time.Now().Nanosecond()/1000/1000) % tw.slotNum
	} else {
		pos = int(time.Now().Unix()) % tw.slotNum
	}

	return pos
}

func (tw *TimeWheel) getTimerInSlot(delay time.Duration) *Timer {
	pos := tw.getPosition(delay)
	return tw.slots[pos]
}

func (tw *TimeWheel) incr() {
	atomic.AddInt64(&tw.counter, 1)
}

func (tw *TimeWheel) deincr() {
	atomic.AddInt64(&tw.counter, -1)
}

func (tw *TimeWheel) loadIncr() int64 {
	return atomic.LoadInt64(&tw.counter)
}

type TimerEntry struct {
	timer  *Timer
	tw     *TimeWheel
	event  *Event
	stoped bool

	C chan time.Time
}

func (te *TimerEntry) init() {
	te.C = make(chan time.Time)
}

func (te *TimerEntry) Stop() {
	// remove event
	te.timer.Del(te.event)
	te.stoped = true
}

func (te *TimerEntry) Reset(delay time.Duration) bool {
	if te.stoped {
		return false
	}

	// del
	te.timer.Del(te.event)

	// add
	if delay < time.Millisecond {
		return false
	}

	var (
		pos   = te.tw.getWritePosition(delay)
		timer = te.tw.slots[pos]
	)

	ev := timer.addAny(delay, te.event.fn, te.event.cron, te.C)
	ev.slotPos = pos
	ev.c = te.C

	te.timer = timer
	te.event = ev

	return true
}