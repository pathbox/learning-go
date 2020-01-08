package ratelimiter

import (
	"sync"
	"time"
)

const (
	lkrateMargin               = 0.01
	lkmaxWait    time.Duration = 0x7fffffffffffffff
)

//漏桶算法
type LeakyBucket struct {
	sync.Mutex
	//LeakyBucket创建时间
	startTime time.Time
	//桶的容量
	capacity int64
	//flow的速度(一个tick内)
	decrement int64
	//flow的周期(nanosecond --> 1/emmissionInterval = tick数)
	emissionInterval time.Duration
	//与capacity的差值
	inBucket int64
	//已经历的tick数
	availTick int64
}

func NewLeakyBucket(interval time.Duration, capacity int64) *LeakyBucket {
	return NewLeakyBucketWithDecrement(interval, capacity, int64(1))
}

func NewLeakyBucketWithRate(rate float32, capacity int64) *LeakyBucket {
	for decrement := int64(1); decrement < 1<<50; decrement = nextDecrement(decrement) {
		interval := 1e9 * float32(decrement) / rate
		if interval < 0 {
			continue
		}
		lb := NewLeakyBucketWithDecrement(time.Duration(interval), capacity, decrement)
		if diff := lb.Rate() - rate; diff <= lkrateMargin {
			return lb
		}
	}
	return nil
}

func NewLeakyBucketWithDecrement(interval time.Duration, capacity int64, decrement int64) *LeakyBucket {
	if interval <= 0 || capacity <= 0 || decrement <= 0 {
		return nil
	}
	return &LeakyBucket{
		startTime:        time.Now(),
		capacity:         capacity,
		decrement:        decrement,
		emissionInterval: interval,
	}
}

func nextDecrement(decrement int64) int64 {
	dcrm1 := decrement * 11 / 10
	if dcrm1 == decrement {
		dcrm1++
	}
	return int64(dcrm1)
}

func (lb *LeakyBucket) Wait(n int64) {
	if t := lb.Cal(n); t != 0 {
		time.Sleep(t)
	}
	return
}

func (lb *LeakyBucket) Cal(n int64) time.Duration {
	if n <= 0 {
		return 0
	}
	now := time.Now()
	currentTick := lb.adjust(now)
	if lb.inBucket+n <= lb.capacity {
		return 0
	}
	overflow := lb.inBucket + n - lb.capacity
	endTick := currentTick + overflow/lb.decrement
	endTime := now.Add(time.Duration(endTick) * lb.emissionInterval)
	waitTime := endTime.Sub(now)
	if waitTime >= lkmaxWait {
		return 0
	}
	lb.inBucket = lb.inBucket + n
	return waitTime
}

func (lb *LeakyBucket) adjust(now time.Time) int64 {
	currentTick := int64(now.Sub(lb.startTime) / lb.emissionInterval)
	if lb.inBucket <= 0 {
		return currentTick
	}
	lb.inBucket -= (currentTick - lb.availTick) * lb.decrement
	if lb.inBucket < 0 {
		lb.inBucket = 0
	}
	lb.availTick = currentTick
	return currentTick
}

func (lb *LeakyBucket) Rate() float32 {
	return 1e9 * float32(lb.decrement) / float32(lb.emissionInterval)
}

func (lb *LeakyBucket) Capacity() int64 {
	return lb.capacity
}
