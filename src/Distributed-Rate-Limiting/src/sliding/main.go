package main

import (
	"fmt"
	"sync"
	"time"
)

type SlidingWindowLimiter struct {
	windowSize   time.Duration //窗口大小
	maxRequests  int           // 最大请求数
	requests     []time.Time   //窗口内的请求时间
	requestsLock sync.Mutex    //请求锁
}

func NewSlidingWindowLimiter(windowSize time.Duration, maxRequests int) *SlidingWindowLimiter {
	return &SlidingWindowLimiter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		requests:    make([]time.Time, 0),
	}
}

func (limiter *SlidingWindowLimiter) AllowRequest() bool {
	limiter.requestsLock.Lock()
	defer limiter.requestsLock.Unlock()

	// 移除过期的请求

	currentTime := time.Now()
	// 当前接口的请求时间和 request数组中的时间进行求差比较，如果时间相距是在窗口范围外，则说明request中的时间过期，需要去除，不断的去处过期的时间，加入新的当前新请求的时间，然后在阈值范围内，则通过，超过阈值，则要限流
	for len(limiter.requests) > 0 && currentTime.Sub(limiter.requests[0]) > limiter.windowSize {
		limiter.requests = limiter.requests[1:] // 把过期的第一个时间数据移除，不断循环
	}

	// 检查请求数是否超过阈值,超过阈值表示限量，没超过阈值，则不限流，做“加”请求数的处理
	if len(limiter.requests) >= limiter.maxRequests {
		return false
	}

	limiter.requests = append(limiter.requests, currentTime)
	return true

}

func main() {
	limiter := NewSlidingWindowLimiter(500*time.Millisecond, 2) // 每秒最多允许4个请求
	for i := 0; i < 100; i++ {
		now := time.Now().Format("15:04:05")
		if limiter.AllowRequest() {
			fmt.Println(now + " 请求通过")
		} else {
			fmt.Println(now + " 请求被限流")
		}
		time.Sleep(100 * time.Millisecond)
	}
}
