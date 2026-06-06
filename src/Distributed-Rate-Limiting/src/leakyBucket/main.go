package main

import (
	"fmt"
	"time"
)

type LeakyBucket struct {
	rate       float64 // 漏桶速率，单位请求数/秒
	capacity   int     // 漏桶容量，最多可存储请求数
	water      int     // 当前水量，表示当前漏桶中的请求数
	lastLeakMs int64   // 上次漏水的时间戳，单位秒
}

func NewLeakyBucket(rate float64, capacity int) *LeakyBucket {
	return &LeakyBucket{
		rate:       rate,
		capacity:   capacity,
		water:      0,
		lastLeakMs: time.Now().Unix(),
	}
}

func (lb *LeakyBucket) Allow() bool {
	now := time.Now().Unix()

	elapsed := now - lb.lastLeakMs // 通过时间间隔来计算这段时间已经漏了多少水
	leakAmount := int(float64(elapsed) / 1000 * lb.rate)

	if leakAmount > 0 {
		if leakAmount > lb.water {
			lb.water = 0
		} else {
			lb.water -= leakAmount
		}
	}

	// 判断当前水量是否超过容量
	if lb.water > lb.capacity {
		lb.water-- // 如果超过容量，减去刚刚增加的水量
		return false
	}

	// 没有限流，增加水量，设置通过时间，也就是最近漏水时间
	lb.water++

	lb.lastLeakMs = now
	return true
}

func main() {
	// 创建一个漏桶，速率为每秒处理3个请求，容量为4个请求
	leakyBucket := NewLeakyBucket(3, 4)

	// 模拟请求
	for i := 1; i <= 100; i++ {
		now := time.Now().Format("15:04:05")
		if leakyBucket.Allow() {
			fmt.Printf(now+"  第 %d 个请求通过\n", i)
		} else {
			fmt.Printf(now+"  第 %d 个请求被限流\n", i)
		}
		time.Sleep(200 * time.Millisecond) // 模拟请求间隔
	}
}
