// w.WorkerChannel <- w.Channel
// worker := <-WorkerChannel

package main

import (
	"fmt"
)

var WorkerChannel = make(chan chan Work)

type Work struct {
	ID   int
	Name string
}

type Worker struct {
	ID            int
	WorkerChannel chan chan Work
	Channel       chan Work
	End           chan struct{}
}

func main() {
	input := make(chan Work)
	worker := Worker{
		ID:            1,
		Channel:       make(chan Work), // receive the job
		WorkerChannel: WorkerChannel,   // belongs to the gloabl WorkerChannel
		End:           make(chan struct{}),
	}

	go func() { // worker process action 的`守护进程`
		for {
			fmt.Println("w.Channel: ", worker.Channel)
			worker.WorkerChannel <- worker.Channel
			select {
			case job := <-worker.Channel:
				fmt.Println("From Channel: ", job)
				fmt.Printf("Process job, job ID:%d, job Name: %s\n", job.ID, job.Name)
			}
		}
	}()

	go start(input)
	go func() { // This is the collector, producer create work, work can be the func()
		i := 0
		for {
			i++
			// time.Sleep(2 * time.Second)
			w := Work{ID: i, Name: "This is a work"}
			input <- w
		}
	}()

	for {
	}
}

func start(input chan Work) {
	for {
		select {
		case job := <-input:
			fmt.Println("job: ", job)
			worker := <-WorkerChannel // wait for available Worker channel
			fmt.Println("The worker from WorkerChannel: ", worker)
			worker <- job // dispatch work(job) to worker
			fmt.Println("after job: ", worker)
		default:
			// time.Sleep(time.Second * 3)
			fmt.Println("default")
		}
	}
}

/*
w.Channel:  0xc000066120
The worker from WorkerChannel:  0xc000066120
after job:  0xc000066120
job:  {1282 This is a work}
From Channel:  {1281 This is a work}
Process job, job ID:1281, job Name: This is a work
w.Channel:  0xc000066120
The worker from WorkerChannel:  0xc000066120
after job:  0xc000066120
job:  {1283 This is a work}
From Channel:  {1282 This is a work}
Process job, job ID:1282, job Name: This is a work
w.Channel:  0xc000066120

这里有四个channel

WorkerChannel chan chan Work
Channel       chan Work
End           chan struct{}

input         chan Work

End 显然是用于结束worker作用的

每个Worker有一个Channel，保存自身对应的Work，Worker就是从这里取出Work来消费处理的，所以外部需要把Work传到这个Channel中

WorkerChannel 是一个全局的channel，这个chan中的每个元素值就是 Channel这个channel
也就是，所有Worker对应的Channel会传入WorkerChannel中，WorkerChannel再提供给外部，接收Work.每次Worker的循环操作都需要worker.WorkerChannel <- worker.Channel， 因为Channel从WorkerChannel中取走之后，就没有再放回去了

input channel 的作用就是接收外部产生的work

逻辑阐述:

Collector 和input绑定， 外部产生的work会先传到 Collector这的input channel中

初始化Worker，每个Worker 已经Channel:make(chan Work)， 所以Worker对应的Channel：w.Channel不是nil，而是一个可以使用的channel

将Worker和全局的WorkerChannel绑定一下

在 worker process action 的`守护进程`中 的这一步： 	worker.WorkerChannel <- worker.Channel 是Worker把自身对应的Channel传入到全局的WorkerChannel保存，然后在
worker := <-WorkerChannel // wait for available Worker.Channel
worker <- job 这样 外部的work、job就传到了每个worker.Channel中，worker再从Channel中取出work进行消费处理

所以，WorkerChannel就相当于一个Channel池，存储所有worker的Channel

work的传递流程

new work => input => w.Channel From WorkerChannel => process work

优化：


WorkerChannel = make(chan chan Work, WorkerCount)
这样WorkerChannel带有缓冲，减少阻塞等待情况

可以不使用channel做pool，使用全局map作为pool
*/
