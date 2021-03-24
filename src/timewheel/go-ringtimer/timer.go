package timewheel

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

var (
	defaultAllocCap     = 1024
	defaultLoopInterval = 500 * time.Millisecond

	ErrStarted    = errors.New("timer has started")
	ErrNotStarted = errors.New("timer has not started")
	ErrStopped    = errors.New("timer has stopped")
)

const (
	// optimize GC when timer is very many
	onSyncPool = false
)

type Timer struct {
	mu         sync.RWMutex
	ctx        context.Context
	cancelFunc context.CancelFunc

	loopInterval time.Duration

	// avoid multi call
	inited  bool
	started int32
	stopped int32

	allocCap int      // the cap that reallocate more events
	free     *Event   // free events
	events   []*Event // min heap array
}

// New returns an timer instance with the default allocate capicity
func NewTimer() *Timer {
	t := &Timer{}
	t.init(defaultAllocCap)
	return t
}

// NewWithCap returns an timer instance with the given allocate capicity
func NewWithCap(cap int) *Timer {
	t := &Timer{}
	t.init(cap)
	return t
}

// Init inits the timer instance with the given allocate capicity
func (t *Timer) Init(cap int) {
	t.init(cap)
	return
}

func (t *Timer) init(cap int) {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.inited {
		return
	}

	t.ctx, t.cancelFunc = context.WithCancel(context.Background())
	t.allocCap = cap
	t.allocate()
	t.inited = true
}

// Len return the length of min heap array.
func (t *Timer) Len() int {
	return len(t.events)
}

// Events return the currently events in the timer.
func (t *Timer) Events() []*Event {
	t.mu.Lock()
	defer t.mu.Unlock()
	return t.events
}

// allocate is used to expand min-heap array
func (t *Timer) allocate() {
	events := make([]Event, t.allocCap)
	t.free = &events[0]
	for i := 0; i < t.allocCap; i++ {
		if i < t.allocCap-1 {
			events[i].next = &events[i+1]
		}
	}
}

// Add is used to add new event
func (t *Timer) Add(ttl time.Duration, fn ExpireFunc) *Event {
	return t.addAny(ttl, fn, false, nil)
}

func (t *Timer) AddWithChan(ttl time.Duration, fn ExpireFunc, notifyQueue chan time.Time) *Event {
	return t.addAny(ttl, fn, false, notifyQueue)
}

// AddCron is used to add new crontab event
func (t *Timer) AddCron(ttl time.Duration, fn ExpireFunc) *Event {
	return t.addAny(ttl, fn, true, nil)
}

func (t *Timer) addAny(ttl time.Duration, fn ExpireFunc, cron bool, notifyQueue chan time.Time) *Event {
	// get from sync.Pool cache
	var (
		event = t.get()
	)

	// init event field value
	event.ttl = ttl
	event.expire = time.Now().Add(ttl)
	event.fn = fn
	event.cron = cron
	if notifyQueue != nil {
		event.c = notifyQueue
	}

	// add event
	t.mu.Lock()
	t.add(event)
	t.mu.Unlock()

	return event
}

// like stdlib time.After()
func (t *Timer) After(ttl time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	t.mu.Lock()
	defer t.mu.Unlock()

	event := t.get()
	event.ttl = ttl
	event.expire = time.Now().Add(ttl)
	event.fn = func() {
		queue <- time.Now()
	}

	t.add(event)
	return queue
}

// like stdlib time.Sleep()
func (t *Timer) Sleep(ttl time.Duration) {
	q := t.After(ttl)
	select {
	case <-t.ctx.Done():
	case <-q:
	}
}

func (t *Timer) get() *Event {
	var event *Event
	if onSyncPool {
		// sync.pool
		event = globalEventPool.get()
		event.alone = true
	} else {
		event = new(Event)
	}

	return event
}

func (t *Timer) add(event *Event) {
	event.index = len(t.events)
	t.events = append(t.events, event)
	t.upEvent(event.index)
	return
}

func (t *Timer) upEvent(j int) {
	for {
		i := (j - 1) / 2
		if i == j || !t.events[j].Less(t.events[i]) {
			break
		}
		t.swapEvent(i, j)
		j = i
	}
}

func (t *Timer) swapEvent(i, j int) {
	t.events[i], t.events[j] = t.events[j], t.events[i]
	t.events[i].index = i
	t.events[j].index = j
}

// Del is used to remove event from timer.
// If event is nil, will retrun directly.
func (t *Timer) Del(event *Event) {
	if event == nil {
		return
	}

	t.mu.Lock()
	defer t.mu.Unlock()
	if t.del(event) || event.alone {
		t.put(event)
	}
}

func (t *Timer) put(event *Event) {
	if onSyncPool {
		globalEventPool.put(event)
	}
}

func (t *Timer) del(event *Event) bool {
	i := event.index
	last := len(t.events) - 1
	if i < 0 || i > last || t.events[i] != event {
		// invalid event or event has been removed
		return false
	}

	if i != last {
		t.swapEvent(i, last)
		t.downEvent(i)
		t.upEvent(i)
	}

	// remove the last event
	t.events[last].index = -1
	t.events = t.events[:last]
	return true
}

func (t *Timer) downEvent(i int) {
	n := len(t.events) - 1
	for {
		left := 2*i + 1
		if left >= n || left < 0 {
			// greather than max index or number overflow
			break
		}
		j := left
		if right := left + 1; right < n && t.events[right].Less(t.events[left]) {
			j = right
		}
		if t.events[i].Less(t.events[j]) {
			break
		}
		t.swapEvent(i, j)
		i = j
	}
}

// Start is used to start the timer.
func (t *Timer) Start() error {
	if atomic.CompareAndSwapInt32(&t.started, 0, 1) {
		go func() {
			for {
				t.loop()
				if t.loopInterval == 0 {
					time.Sleep(defaultLoopInterval)
					continue
				}

				time.Sleep(t.loopInterval)
			}
		}()
		return nil
	}
	return ErrStarted
}

// Stop is used to stop the timer.
func (t *Timer) Stop() error {
	if !atomic.CompareAndSwapInt32(&t.started, 1, 1) {
		return ErrNotStarted
	}
	if atomic.CompareAndSwapInt32(&t.stopped, 0, 1) {
		t.cancelCtx()
		return nil
	}
	return ErrStopped
}

// cancelCtx
func (t *Timer) cancelCtx() {
	if t.cancelFunc != nil {
		t.cancelFunc()
	}
}

// LoopOnce run cone
func (t *Timer) LoopOnce() {
	t.loop()
}

func (t *Timer) loop() {
	var (
		d     time.Duration
		fn    ExpireFunc
		event *Event
	)

	select {
	case <-t.ctx.Done():
		return
	default:
		for {
			t.mu.RLock()
			if len(t.events) == 0 {
				t.mu.RUnlock()
				break
			}

			event = t.events[0]
			if d = event.Delay(); d >= 0 {
				t.mu.RUnlock()
				break
			}
			t.mu.RUnlock()

			t.mu.Lock()
			fn = event.fn
			if event.cron {
				t.Set(event, event.ttl)
			} else {
				t.del(event)
			}
			t.mu.Unlock()

			if fn != nil {
				go fn()
			}

			if event.c != nil {
				event.sendNotify()
			}
		}
	}
}

// IsStopped is used to show the timer is whether stopped.
func (t *Timer) IsStopped() bool {
	return atomic.CompareAndSwapInt32(&t.stopped, 1, 1)
}

// Reset alias Set func
func (t *Timer) Reset(event *Event, ttl time.Duration) bool {
	return t.Set(event, ttl)
}

// Set
func (t *Timer) Set(event *Event, ttl time.Duration) bool {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.del(event)

	// events has been recyled
	if !event.alone {
		return false
	}

	event.ttl = ttl
	event.expire = time.Now().Add(ttl)
	t.add(event)
	return true
}

func stdlog(v ...interface{}) {
	fmt.Println(v...)
}