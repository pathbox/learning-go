# 四种主流限流算法详解

> 限流（Rate Limiting）是分布式系统中保护服务、防止过载的重要手段。本文介绍四种最常见的限流算法：固定窗口、滑动窗口、漏桶、令牌桶。

---

## 目录

- [1. 固定窗口限流（Fixed Window Counter）](#1-固定窗口限流fixed-window-counter)
- [2. 滑动窗口限流（Sliding Window）](#2-滑动窗口限流sliding-window)
- [3. 漏桶限流（Leaky Bucket）](#3-漏桶限流leaky-bucket)
- [4. 令牌桶限流（Token Bucket）](#4-令牌桶限流token-bucket)
- [横向对比](#横向对比)
- [选型建议](#选型建议)

---

## 1. 固定窗口限流（Fixed Window Counter）

### 原理

把时间切成等长的窗口（比如每 1 秒一个窗口），每个窗口内维护一个计数器：
- 请求进来 `count++`
- 超过阈值就拒绝
- 窗口切换时计数器清零

```
时间轴: |--- 1s ---|--- 1s ---|--- 1s ---|
计数:    [count=80]  [count=20]  [count=...]
阈值:    100         100         100
```

### 伪代码

```go
type FixedWindow struct {
    count     int64
    windowEnd time.Time
    limit     int64
    duration  time.Duration
}

func (f *FixedWindow) Allow() bool {
    now := time.Now()
    if now.After(f.windowEnd) {
        f.count = 0
        f.windowEnd = now.Add(f.duration)
    }
    if f.count >= f.limit {
        return false
    }
    f.count++
    return true
}
```

Redis 实现：

```bash
INCR  rate_limit:user:123:1696500000
EXPIRE rate_limit:user:123:1696500000 1
```

### 优点

- **实现极简**：一个计数器 + 一个时间戳即可
- **性能高**：O(1) 时间复杂度，内存占用极小
- **易理解**：业务方一看就懂

### 缺点

- **临界突刺问题**（最致命）：窗口边界处可能出现 **2 倍流量**
  - 例如：限制 100 QPS，第 1 秒的最后 100ms 来了 100 个请求，第 2 秒的前 100ms 又来了 100 个，那么这 200ms 内实际通过了 200 个请求，瞬时 QPS 翻倍击穿下游
- 流量分布不均匀

### 适用场景

- 对精度要求不高的场景：API 日调用次数限制、用户每天发帖次数
- 简单的防刷场景
- 资源受限、追求简单的内部系统

---

## 2. 滑动窗口限流（Sliding Window）

### 原理

为了解决固定窗口的临界问题，把窗口"滑动"起来。两种常见实现：

#### (a) 滑动窗口日志（Sliding Window Log）

记录每个请求的精确时间戳（通常用 Redis ZSet 存），每次请求时：

1. 删除窗口外的旧记录（`ZREMRANGEBYSCORE`）
2. 统计窗口内的数量（`ZCARD`）
3. 未超阈值则插入当前时间戳

#### (b) 滑动窗口计数（Sliding Window Counter）

把大窗口切成 N 个小格子，滑动时按比例计算：

```
当前时间在第 2 个窗口的 30% 处
有效计数 = 第1窗口count × 70% + 第2窗口count
```

### 伪代码（Redis Lua 实现日志方式）

```lua
local key = KEYS[1]
local now = tonumber(ARGV[1])
local window = tonumber(ARGV[2])
local limit = tonumber(ARGV[3])

redis.call('ZREMRANGEBYSCORE', key, 0, now - window)
local count = redis.call('ZCARD', key)
if count >= limit then
    return 0
end
redis.call('ZADD', key, now, now)
redis.call('EXPIRE', key, window / 1000)
return 1
```

### 优点

- **平滑限流**：彻底解决临界突刺问题
- **精度可调**：格子越细，精度越高
- 流量曲线更平滑

### 缺点

- **内存开销大**（日志方式）：高 QPS 下每个请求都要存时间戳，10w QPS × 1s 窗口就是 10w 条记录
- **实现复杂**：尤其分布式场景下要保证原子性，通常需要 Lua 脚本
- 性能略低于固定窗口

### 适用场景

- 对限流精度要求高的场景：核心接口防刷、登录尝试次数限制
- 风控、反爬虫
- 用户级精细化限流（每个用户独立窗口）

---

## 3. 漏桶限流（Leaky Bucket）

### 原理

想象一个底部有小孔的水桶：

- 请求 = 注水（任意速率，桶满则拒绝/排队）
- 处理 = 漏水（**恒定速率**流出）

```
   请求流入（任意速率）
        ↓ ↓ ↓ ↓ ↓
       ┌─────────┐
       │ ▓▓▓▓▓▓▓ │  桶容量 = 100
       │ ▓▓▓▓▓▓▓ │
       └────┬────┘
            ↓ 恒定速率漏出 = 10/s
         下游服务
```

### 伪代码

```go
type LeakyBucket struct {
    capacity   int64     // 桶容量
    water      int64     // 当前水量
    leakRate   int64     // 漏出速率（个/秒）
    lastLeakAt time.Time
}

func (l *LeakyBucket) Allow() bool {
    now := time.Now()
    elapsed := now.Sub(l.lastLeakAt).Seconds()
    leaked := int64(elapsed * float64(l.leakRate))
    l.water = max(0, l.water - leaked)
    l.lastLeakAt = now

    if l.water >= l.capacity {
        return false
    }
    l.water++
    return true
}
```

### 优点

- **强制平滑输出**：无论入口流量多突兀，下游永远是恒定速率
- **保护下游**：对下游系统极其友好（数据库、第三方 API）
- 实现相对简单（一个队列 + 定时器）

### 缺点

- **无法应对突发流量**：即使系统空闲，突发请求也只能排队或被拒，**资源利用率低**
- **可能增加延迟**：请求排队等待漏出
- 不适合需要"快速响应"的场景

### 适用场景

- **下游是脆弱系统**：调用第三方 API（短信、支付）有严格 QPS 限制
- 需要强制平滑输出：消息推送、爬虫调度
- 流量整形（Traffic Shaping）：网络包发送、视频流推送

---

## 4. 令牌桶限流（Token Bucket）

### 原理

桶里以恒定速率投放令牌（满则丢弃多余令牌）：

- 请求来了先**取令牌**，有令牌就放行，没令牌就拒绝/等待
- 桶容量 = 允许的最大突发量

```
   令牌生成器（恒定速率 10/s）
            ↓
       ┌─────────┐
       │ ●●●●●●● │  桶容量 = 100
       │ ●●●●●●● │
       └────┬────┘
            ↑ 取令牌
        请求流入
```

### 伪代码

```go
type TokenBucket struct {
    capacity   int64     // 桶容量
    tokens     int64     // 当前令牌数
    rate       int64     // 令牌生成速率（个/秒）
    lastFillAt time.Time
}

func (t *TokenBucket) Allow() bool {
    now := time.Now()
    elapsed := now.Sub(t.lastFillAt).Seconds()
    add := int64(elapsed * float64(t.rate))
    t.tokens = min(t.capacity, t.tokens + add)
    t.lastFillAt = now

    if t.tokens <= 0 {
        return false
    }
    t.tokens--
    return true
}
```

### 优点

- **允许突发流量**：桶里攒了令牌时可以瞬间放行一批，**资源利用率高**
- **平均速率有保障**：长期看仍然受令牌生成速率限制
- **灵活**：通过调整桶容量和生成速率，可同时控制"平均速率"和"突发上限"
- 业界标杆：Google Guava `RateLimiter`、Nginx `limit_req`、Sentinel 都基于此

### 缺点

- 实现略复杂于固定窗口
- 突发流量可能瞬间打满下游（若下游不能承受突发，应改用漏桶）
- 分布式实现需要原子操作（Redis + Lua）

### 适用场景

- **API 网关限流**：允许一定突发，整体保持平均速率
- 用户行为限流：登录、支付、下单接口
- 大多数 Web 服务的通用限流方案
- 微服务网关、Service Mesh 限流

---

## 横向对比

| 维度           | 固定窗口        | 滑动窗口          | 漏桶              | 令牌桶            |
| -------------- | --------------- | ----------------- | ----------------- | ----------------- |
| **实现难度**   | ⭐ 极简         | ⭐⭐⭐ 较复杂     | ⭐⭐ 中等         | ⭐⭐ 中等         |
| **性能**       | 极高            | 中                | 高                | 高                |
| **内存占用**   | 极小            | 较大（日志方式）  | 小                | 小                |
| **临界突刺**   | ❌ 有           | ✅ 无             | ✅ 无             | ✅ 无             |
| **允许突发**   | 部分允许        | 不允许            | ❌ 不允许         | ✅ 允许           |
| **平滑输出**   | ❌              | ✅                | ✅✅ 强制平滑      | ⭕ 平均平滑       |
| **资源利用率** | 中              | 中                | 低                | 高                |
| **延迟**       | 低              | 低                | 可能高（排队）    | 低                |
| **典型实现**   | Redis INCR      | Redis ZSet + Lua  | 队列 + 定时器     | Guava/Nginx       |

---

## 选型建议

### 决策树

```
是否需要严格保护下游（如第三方 API、DB）？
├── 是 → 漏桶（强制恒定速率）
└── 否 → 是否允许突发流量？
        ├── 是 → 令牌桶（推荐，业界主流）
        └── 否 → 是否对精度要求高？
                ├── 是 → 滑动窗口
                └── 否 → 固定窗口（最简单）
```

### 常见场景对应

| 场景                          | 推荐算法     | 原因                                 |
| ----------------------------- | ------------ | ------------------------------------ |
| API 网关通用限流              | 令牌桶       | 兼顾平均速率与突发能力               |
| 调用第三方支付/短信 API       | 漏桶         | 第三方有严格 QPS 限制，必须平滑      |
| 用户每日发帖次数              | 固定窗口     | 精度要求低，按天计数即可             |
| 登录失败次数限制              | 滑动窗口     | 需精确防爆破，避免临界绕过           |
| 秒杀场景前端拦截              | 令牌桶       | 允许突发抢购，桶容量控制总放行量     |
| 网络包/视频流发送             | 漏桶         | 必须恒定速率，避免网络抖动           |
| 反爬虫 IP 限频                | 滑动窗口     | 需要精确统计单位时间内访问次数       |
| 微服务内部调用限流            | 令牌桶       | 高性能，支持突发，业界标准           |

---

## 进阶思考

1. **单机 vs 分布式**：单机用内存即可（如 `golang.org/x/time/rate` 就是令牌桶）；分布式需要 Redis + Lua 保证原子性
2. **限流维度**：可按 IP / 用户 ID / 接口 / 租户 多维度组合
3. **降级策略**：限流后是直接拒绝（fail-fast）、排队等待，还是降级返回兜底数据
4. **预热（Warm Up）**：冷启动时令牌桶逐步放量，避免突然全开
5. **自适应限流**：根据下游 RT、CPU 负载动态调整限流阈值（如 Sentinel 的 BBR 算法）

---

## 参考实现

- Go 标准库：`golang.org/x/time/rate`（令牌桶）
- Java：Google Guava `RateLimiter`（令牌桶）、Resilience4j、Sentinel
- Nginx：`limit_req_module`（漏桶）、`limit_conn_module`
- 分布式：Redis-Cell 模块、Sentinel-Cluster
