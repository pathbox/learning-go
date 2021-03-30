package main
import (
    "fmt"
    "sync"
    "time"
)

// BucketLimiter 定义漏桶算法struct
type BucketLimiter struct {
    Lck      *sync.Mutex // 锁
    Rate     float64     //最大速率限制
    Balance  float64     //漏桶的余量
    Cap      float64     //漏桶的最大容量限制
    LastTime time.Time   //上次检查的时间
}

 // NewBucketLimiter 初始化BucketLimiter
func NewBucketLimiter(rate int, cap int) *BucketLimiter {
    return &BucketLimiter{
        Lck:      new(sync.Mutex),
        Rate:     float64(rate),
        Balance:  float64(cap),
        Cap:      float64(cap),
        LastTime: time.Now(),
    }
}

// leakyBucket 漏桶算法实现
func (r *BucketLimiter) LeakyBucket() bool {
    ok := false
    r.Lck.Lock()
    defer r.Lck.Unlock()
    now := time.Now()
    dur := now.Sub(r.LastTime).Seconds() //当前时间与上一次检查时间差
    r.LastTime = now
    water := dur * r.Rate //计算这段时间内漏桶流出水的流量water
    r.Balance += water    //漏桶流出water容量的水，自然漏桶的余量多出water
    if r.Balance > r.Cap {
        r.Balance = r.Cap
    }
    if r.Balance >= 1 { //漏桶余量足够容下当前的请求
        r.Balance -= 1
        ok = true
    }
    return ok
}
func main() {

    // 初始化 限制每秒2个请求 漏洞容量为5
    r := NewBucketLimiter(2, 5)
    for i := 0; i < 20; i++ {
        ok := r.LeakyBucket()
        if ok {
            fmt.Println("pass ", i)
        } else {
            fmt.Println("limit ", i)
        }
        time.Sleep(100 * time.Millisecond)
    }
}