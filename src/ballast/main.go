package main

import (
    "runtime"
    "math"
    "time"
)

func main() {
    ballast := make([]byte, 10*1024*1024*1024)

    <-time.After(time.Duration(math.MaxInt64))
    runtime.KeepAlive(ballast)
}

/* https://gocn.vip/topics/17480
上面的代码就初始化了一个 ballast，利用 runtime.KeepAlive 来保证 ballast 不会被 GC 给回收掉。

利用这个特性，就能保证 GC 在 10G 的一倍时才能被触发，这样就能够比较精准控制 GO GC 的触发时机。

这里你可能有一个疑问，这里初始化一个 10G 的数组，不就占用了 10 G 的物理内存呢？ 答案其实是不会的。

*/