package timewheel

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var globalEventPool = newEventPool()

// ExpireFunc represents a function will be executed when a event is trigged.
type ExpireFunc func()

// An Event represents an elemenet of the events in the timer.
type Event struct {
	slotPos int            // mark timeWheel slot index
	index   int            // index in the min heap structure
	c       chan time.Time // like go time.NewTimer().C
	closed  int32

	ttl    time.Duration // wait delay time
	expire time.Time     // due timestamp
	fn     ExpireFunc    // callback function

	next    *Event
	cron    bool // repeat task
	cronNum int  // cron circle num
	alone   bool // indicates event is alone or in the free linked-list of timer
}

func (e *Event) init() {
	e.c = make(chan time.Time)
}

func (e *Event) close() {
	var ok = atomic.CompareAndSwapInt32(&e.closed, 0, 1)
	if ok {
		return
	}

	close(e.c)
}

func (e *Event) sendNotify() {
	if atomic.CompareAndSwapInt32(&e.closed, 1, 1) {
		// already closed
		return
	}

	select {
	case e.c <- time.Now():
	default:
		return
	}
}

// clear field
func (e *Event) clear() {
	e.index = 0
	e.slotPos = 0
	e.cron = false
	e.fn = nil
	e.alone = false
	e.c = nil
}

// Less is used to compare expiration with other events.
func (e *Event) Less(o *Event) bool {
	return e.expire.Before(o.expire)
}

// Delay is used to give the duration that event will expire.
func (e *Event) Delay() time.Duration {
	return e.expire.Sub(time.Now())
}

func (e *Event) String() string {
	return fmt.Sprintf("index %d ttl %v, expire at %v", e.index, e.ttl, e.expire)
}

func newEventPool() *eventPool {
	return &eventPool{}
}

type eventPool struct {
	p sync.Pool
}

func (ep *eventPool) get() *Event {
	if t, _ := ep.p.Get().(*Event); t != nil {
		t.clear()
		return t
	}

	return new(Event)
}

func (ep *eventPool) put(ev *Event) {
	ep.p.Put(ev)
}