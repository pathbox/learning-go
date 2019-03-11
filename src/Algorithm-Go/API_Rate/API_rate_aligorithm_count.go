// 接口限流算法 计数算法, 在interval的时间内，计数加一，检查次数是否超过限制。时间过期了，重置时间和次数

package main

import (
	"fmt"
	"time"
)

func main() {
	req_count := 0
	limit := 1000
	interval := int64(60)
	t := time.Now().Unix()
	result := rate_count(req_count, limit, interval, t)
	fmt.Println("The result is:", result)
}

func rate_count(req_count int, limit int, interval int64, t int64) bool {
	now := time.Now().Unix()
	if now < t+interval {
		req_count++
		return req_count <= limit
	} else {
		t = now
		req_count = 1
		return true
	}
}
