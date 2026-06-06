package main

import (
	"fmt"
	"sync"
	"time"
)

// TokenBucket 表示一个令牌桶。
type TokenBucket struct {
	rate       float64    // 令牌添加到桶中的速率。
	capacity   float64    // 桶的最大容量。
	tokens     float64    // 当前桶中的令牌数量。
	lastUpdate time.Time  // 上次更新令牌数量的时间。
	mu         sync.Mutex // 互斥锁，确保线程安全。
}

// NewTokenBucket 创建一个新的令牌桶，给定令牌添加速率和桶的容量。
func NewTokenBucket(rate float64, capacity float64) *TokenBucket {
	return &TokenBucket{
		rate:       rate,
		capacity:   capacity,
		tokens:     capacity, // 初始时，桶是满的。
		lastUpdate: time.Now(),
	}
}

// Allow 检查是否可以从桶中移除一个令牌。如果可以，它移除一个令牌并返回 true。
// 如果不可以，它返回 false。
func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()

	// 根据经过的时间计算要添加的令牌数量。
	now := time.Now()
	elapsed := now.Sub(tb.lastUpdate).Seconds()
	tokensToAdd := elapsed * tb.rate

	tb.tokens += tokensToAdd
	if tb.tokens > tb.capacity {
		tb.tokens = tb.capacity // 确保令牌数量不超过桶的容量。
	}

	if tb.tokens >= 1.0 {
		tb.tokens--
		tb.lastUpdate = now
		return true
	}
	// 如果桶中没有足够的令牌，请求被拒绝。  令牌桶和漏桶的思维相反，令牌桶是桶里没令牌了，空了限流，漏桶是桶里的水满了限流。一个是生产的速率是固定的，没有可消费的了，一个消费的速率是固定，生产过快了溢出，容量满了
	return false
}

func main() {
	tokenBucket := NewTokenBucket(2.0, 3.0)

	for i := 1; i <= 100; i++ {
		now := time.Now().Format("15:04:05")
		if tokenBucket.Allow() {
			fmt.Printf(now+"  第 %d 个请求通过\n", i)
		} else { // 如果不能移除一个令牌，请求被拒绝。
			fmt.Printf(now+"  第 %d 个请求被限流\n", i)
		}
		time.Sleep(200 * time.Millisecond)
	}
}

/*
16:30:03  第 1 个请求通过
16:30:04  第 2 个请求通过
16:30:04  第 3 个请求通过
16:30:04  第 4 个请求通过
16:30:04  第 5 个请求被限流
16:30:04  第 6 个请求通过
16:30:05  第 7 个请求被限流
16:30:05  第 8 个请求通过
16:30:05  第 9 个请求通过
16:30:05  第 10 个请求被限流
16:30:05  第 11 个请求通过
16:30:06  第 12 个请求被限流
16:30:06  第 13 个请求通过
16:30:06  第 14 个请求被限流
16:30:06  第 15 个请求通过
16:30:06  第 16 个请求通过
16:30:07  第 17 个请求被限流
16:30:07  第 18 个请求通过
16:30:07  第 19 个请求被限流
16:30:07  第 20 个请求通过
16:30:07  第 21 个请求被限流
16:30:08  第 22 个请求通过
16:30:08  第 23 个请求通过
16:30:08  第 24 个请求被限流
16:30:08  第 25 个请求通过
16:30:08  第 26 个请求被限流
16:30:09  第 27 个请求通过
16:30:09  第 28 个请求被限流
16:30:09  第 29 个请求通过
16:30:09  第 30 个请求通过
16:30:09  第 31 个请求被限流
16:30:10  第 32 个请求通过
16:30:10  第 33 个请求被限流
16:30:10  第 34 个请求通过
16:30:10  第 35 个请求被限流
16:30:10  第 36 个请求通过
16:30:11  第 37 个请求通过
16:30:11  第 38 个请求被限流
16:30:11  第 39 个请求通过
16:30:11  第 40 个请求被限流
16:30:11  第 41 个请求通过
16:30:12  第 42 个请求被限流
16:30:12  第 43 个请求通过
16:30:12  第 44 个请求通过
16:30:12  第 45 个请求被限流
16:30:12  第 46 个请求通过
16:30:13  第 47 个请求被限流
16:30:13  第 48 个请求通过
16:30:13  第 49 个请求通过
16:30:13  第 50 个请求被限流
16:30:13  第 51 个请求通过
16:30:14  第 52 个请求被限流
*/
