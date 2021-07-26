package main
import (
    "fmt"
    "math"
    "sync"
    "time"
)

// TokenBucket 定义令牌桶结构
type TokenBucket struct {
    LastTime time.Time // 当前请求时间
    Capacity float64   // 桶的容量（存放令牌的最大量）
    Rate     float64   // 令牌放入速度
    Tokens   float64   // 当前令牌总量
    Lck      *sync.Mutex
}

 // NewTokenBucket 初始化TokenBucket
func NewTokenBucket(rate int, cap int) *TokenBucket {
    return &TokenBucket{
        LastTime: time.Now(),
        Capacity: float64(cap),
        Rate:     float64(rate),
        Tokens:   float64(cap),
        Lck:      new(sync.Mutex),
    }
}

// getToken 判断是否获取令牌（若能获取，则处理请求）
func (r *TokenBucket) GetToken() bool {
    now := time.Now()
    r.Lck.Lock()
    defer r.Lck.Unlock()
    // 先添加令牌
    tokens := math.Min(r.Capacity, r.Tokens+now.Sub(r.LastTime).Seconds()*r.Rate)
    r.Tokens = tokens
    if tokens < 1 {
        // 若桶中一个令牌都没有了，则拒绝
        return false
    } else {
        // 桶中还有令牌，领取令牌
        r.Tokens -= 1
        r.LastTime = now
        return true
    }
}
 func main() {

    // 初始化 限制每秒2个请求 令牌桶容量为5
    r := NewTokenBucket(2, 5)
    for i := 0; i < 20; i++ {
        ok := r.GetToken()
        if ok {
            fmt.Println("pass ", i)
        } else {
            fmt.Println("limit ", i)
        }
        time.Sleep(100 * time.Millisecond)
    }
}