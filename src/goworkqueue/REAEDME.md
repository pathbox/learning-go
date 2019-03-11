### Understand Worker Queue Job


```go
type Queue struct {
	Jobs    chan interface{} // 生产者角色job 需要处理的job数量
	done    chan bool
	workers chan chan int // 消费者角色worker  处理job的多线程能力， worker的目的是多线程处理 job
	mux     sync.Mutex
}
```

- Queue 是一个队列work 队列
  - Jobs 接收每个job需要的参数入队列, 是Queue的入口, Queue和生产者的连接层
  - workers 工人，职蚁。 从Jobs不断取出job进行处理和消费，workers的数量,表示Queue起了多少线程(goroutinue)在并发处理Jobs队列,表示Queue的并发处理能力。每个worker应该是在不同的线程(goroutinue)中执行。worker 应该是最终的`消费者`

- 一个Queue就是消费者和生产者的模型
