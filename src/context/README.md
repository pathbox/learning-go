WithCancel 对应的是 cancelCtx ,其中，返回一个 cancelCtx ，同时返回一个 CancelFunc，CancelFunc 是 context 包中定义的一个函数类型：type CancelFunc func()。调用这个 CancelFunc 时，关闭对应的c.done，也就是让他的后代goroutine退出

WithDeadline 和 WithTimeout 对应的是 timerCtx ，WithDeadline 和 WithTimeout 是相似的，WithDeadline 是设置具体的 deadline 时间，到达 deadline 的时候，后代 goroutine 退出，而 WithTimeout 简单粗暴，直接 return WithDeadline(parent, time.Now().Add(timeout))

WithValue 对应 valueCtx ，WithValue 是在 Context 中设置一个 map，拿到这个 Context 以及它的后代的 goroutine 都可以拿到 map 里的值

context的创建
所有的context的父对象，也叫根对象，是一个空的context，它不能被取消，它没有值，从不会被取消，也没有超时时间，它常常作为处理request的顶层context存在，然后通过WithCancel、WithTimeout函数来创建子对象来获得cancel、timeout的能力

当顶层的request请求函数结束后，我们就可以cancel掉某个context，从而通知别的routine结束

WithValue方法可以把键值对加入context中，让不同的routine获取

// 在 handle 环境中使用
```go

func handleSearch(w http.ResponseWriter, req *http.Request) {
    // ctx is the Context for this handler. Calling cancel closes the
    // ctx.Done channel, which is the cancellation signal for requests
    // started by this handler.
    var (
        ctx    context.Context
        cancel context.CancelFunc
    )
    // 获取参数 ...
    timeout, err := time.ParseDuration(req.FormValue("timeout"))
    if err == nil {
        // The request has a timeout, so create a context that is
        // canceled automatically when the timeout expires.
        // 获取成功, 则按照参数设置超时时间
        ctx, cancel = context.WithTimeout(context.Background(), timeout)
    } else {
        // 获取失败, 则在该函数结束时结束 ...
        ctx, cancel = context.WithCancel(context.Background())
    }
    // ----------------
    // 这样随着cancel的执行,所有的线程都随之结束了 ...
    go A(ctx) +1
    go B(ctx) +2
    go C(ctx) +3
    // ----------------
    defer cancel() // Cancel ctx as soon as handleSearch returns.
}
// 监听 ctx.Done() 结束 ...
func A(ctx context.Context) int {
    // ... TODO
    select {
    case <-ctx.Done():
            return -1
    default:
        // 没有结束 ... 执行 ...
    }
}
```

