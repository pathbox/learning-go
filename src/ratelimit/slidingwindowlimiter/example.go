package main

func main() {

}

// SlidingWindowLimiter 滑动窗口计数器限流
type SlidingWindowLimiter struct {
    Interval    int64       // 总计数时间
    LastTime    time.Time   // 上一个窗口时间
    Lck         *sync.Mutex // 锁
    WinCount    []int64     // 窗口中请求当前数量
    TicketSize  int64       // 窗口最大容量
    TicketCount int64       // 窗口个数
    CurIndex    int64       // 目前使用的窗口下标
}

 // NewSlidingWindowLimiter 初始化滑动窗口计数器限流
func NewSlidingWindowLimiter(interval int64, ticketCount int64, ticketSize int64) *SlidingWindowLimiter {
    return &SlidingWindowLimiter{
        Interval:    interval,
        LastTime:    time.Now(),
        TicketSize:  ticketSize,
        TicketCount: ticketCount,
        WinCount:    make([]int64, ticketSize, ticketSize),
        CurIndex:    0,
        Lck:         new(sync.Mutex),
    }
}

// slidingCounterLimit 滑动窗口计数器限流实现
func (r *SlidingWindowLimiter) slidingCounterLimit() bool {
    r.Lck.Lock()
    defer r.Lck.Unlock()
    eachTicketTime := r.Interval / r.TicketCount
    now := time.Now()
    // 如果间隔时间超过一个窗口的时间 当前窗口置0 指向下一个窗口
    if now.Unix()-r.LastTime.Unix() > eachTicketTime {
        r.WinCount[r.CurIndex] = 0
        r.CurIndex = (r.CurIndex + 1) % r.TicketCount
        r.LastTime = now
    }
    fmt.Println("当前窗口:", r.CurIndex)
    // 当前窗口未满则正常计数
    if r.WinCount[r.CurIndex] < r.TicketSize {
        r.WinCount[r.CurIndex]++
        return true
    }
    return false
}
 func main() {
    // 定义1秒10个时间窗口 每个窗口大小为1  即1秒10个请求
    r := NewSlidingWindowLimiter(1, 10, 1)
    for i := 0; i < 20; i++ {
        ok := r.slidingCounterLimit()
        if ok {
            fmt.Println("pass ", i)
        } else {
            fmt.Println("limit ", i)
        }
        time.Sleep(100 * time.Millisecond)
    }
}