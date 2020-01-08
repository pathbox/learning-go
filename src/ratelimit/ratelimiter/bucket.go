package ratelimiter

import (
	"sync"
	"time"
)

const (
	maxWait    time.Duration = 0x7fffffffffffffff
	rateMargin float32       = 0.01
)

// Token Bucket
type Bucket struct {
	sync.Mutex
	//Token Bucket创建时间
	startTime time.Time
	//容量(QPS/TPS值)
	capacity int64
	//填充周期
	fillInterval time.Duration
	//填充单位值
	increment int64
	//剩余容量
	avail int64
	//已走的Tick数
	availTick int64
}

//初始化Token Bucket
func NewBucket(interval time.Duration, capacity int64) *Bucket {
	return NewBucketWithIncrement(interval, capacity, int64(1))
}

//按rate初始化Token Bucket
func NewBucketWithRate(rate float32, capacity int64) *Bucket {
	//根据rate计算increment和interval
	for increment := int64(1); increment <= 1<<50; increment = nextIncrement(increment) {
		fillInterval := (1e9 * float32(increment) / rate)
		if fillInterval < 0 {
			continue
		}
		b := NewBucketWithIncrement(time.Duration(fillInterval), capacity, int64(increment))
		if margin := b.Rate() - rate; margin <= rateMargin {
			return b
		}
	}
	return nil
}

func NewBucketWithIncrement(interval time.Duration, capacity int64, increment int64) *Bucket {
	if interval <= 0 || capacity <= 0 || increment <= 0 {
		return nil
	}
	return &Bucket{
		startTime:    time.Now(),
		capacity:     capacity,
		increment:    increment,
		fillInterval: interval,
		avail:        capacity,
	}
}

func nextIncrement(increment int64) int64 {
	icrm1 := increment * 11 / 10
	if icrm1 == increment {
		icrm1++
	}
	return int64(icrm1)
}

//流量控制函数
func (tb *Bucket) Wait(n int64) {
	if t := tb.Cal(n); t != 0 {
		time.Sleep(t)
	}
	return
}

//计算是否需要等待
func (tb *Bucket) Cal(n int64) time.Duration {
	if n <= 0 {
		return 0
	}
	tb.Lock()
	defer tb.Unlock()
	now := time.Now()
	currentTick := tb.adjust(now)
	left := tb.avail - n
	if left >= 0 {
		tb.avail = left
		return 0
	}
	endTick := currentTick + -left/tb.increment
	endTime := tb.startTime.Add(time.Duration(endTick) * tb.fillInterval)
	waitTime := endTime.Sub(now)
	if waitTime > maxWait {
		return 0
	}
	tb.avail = left
	return waitTime
}

func (tb *Bucket) adjust(now time.Time) int64 {
	currentTick := int64(now.Sub(tb.startTime) / tb.fillInterval)
	if tb.avail >= tb.capacity {
		return currentTick
	}
	tb.avail += (currentTick - tb.availTick) * tb.increment
	if tb.avail > tb.capacity {
		tb.avail = tb.capacity
	}
	tb.availTick = currentTick
	return currentTick
}

func (tb *Bucket) Rate() float32 {
	r := float32(tb.increment) * 1e9 / float32(tb.fillInterval)
	return r
}

func (tb *Bucket) Capacity() int64 {
	return tb.capacity
}
