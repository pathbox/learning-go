package main

import (
	"fmt"
	"sync"
	"time"
)

type FixedWindowLimiter struct {
	windowSize  time.Duration // 窗口大小
	maxRequests int           // 最大请求数
	requests    int           // 当前窗口请求数
	lastRest    int64         // 上次窗口重置时间(秒级时间戳)
	resetMutex  sync.Mutex    // 重置锁
}

func NewFixedWindowLimiter(windowSize time.Duration, maxRequests int) *FixedWindowLimiter {
	return &FixedWindowLimiter{
		windowSize:  windowSize,
		maxRequests: maxRequests,
		lastRest:    time.Now().Unix(),
	}
}

func (limiter *FixedWindowLimiter) AllowRequest() bool {
	limiter.resetMutex.Lock()
	defer limiter.resetMutex.Unlock()

	// 检查是否x重置窗口
	if time.Now().Unix()-limiter.lastRest >= int64(limiter.windowSize.Seconds()) {
		limiter.requests = 0
		limiter.lastRest = time.Now().Unix()
	}

	// 检查请求数是否超过阈值
	if limiter.requests >= limiter.maxRequests {
		return false
	}

	limiter.requests++
	return true

}

func main() {
	limiter := NewFixedWindowLimiter(1*time.Second, 3) // 每秒最多允许3个请求，算QPS的话 一般是1秒

	for i := 0; i < 100; i++ {
		now := time.Now().Format("15:04:05")
		if limiter.AllowRequest() {
			fmt.Println(now + ":pass")
		} else {
			fmt.Println(now + ":limited")
		}
		time.Sleep(100 * time.Millisecond)
	}
}

/*
优点：

灵活性：滑动窗口算法可以根据实际情况动态调整窗口的大小，以适应流量的变化。这种灵活性使得算法能够更好地应对突发流量和请求分布不均匀的情况。

实时性：由于滑动窗口算法在每个时间窗口结束时都会进行窗口滑动，它能够更及时地响应流量的变化，提供更实时的限流效果。

精度：相比于固定窗口算法，滑动窗口算法的颗粒度更小，可以提供更精确的限流控制。

缺点：

内存消耗：滑动窗口算法需要维护一个窗口内的请求时间列表，随着时间的推移，列表的长度会增长。这可能会导致较大的内存消耗，特别是在窗口大小较大或请求频率较高的情况下。

算法复杂性：相比于简单的固定窗口算法，滑动窗口算法的实现较为复杂。它需要处理窗口滑动、请求计数和过期请求的移除等逻辑，可能需要更多的代码和计算开销
*/
