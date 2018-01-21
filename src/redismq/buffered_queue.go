package redismq

import (
	"fmt"
	"time"
)

// BufferedQueue provides an queue with buffered writes for increased performance.
// Only one buffered queue (per name) can be started at a time.
// Before terminating the queue should be flushed using FlushBuffer() to avoid package loss
type BufferedQueue struct {
	*Queue
	BufferSize   int
	Buffer       chan *Package // 使用chan 作为缓冲存储结构
	nextWrite    int64
	flushStatus  chan (chan bool)
	flushCommand chan bool
}

// CreateBufferedQueue returns BufferedQueue.
// To start writing the buffer to redis use Start().
// Optimal BufferSize seems to be around 200.
// Works like SelectBufferedQueue for existing queues

func CreateBufferedQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string, bufferSize int) *BufferedQueue {
	q := CreateQueue(redisHost, redisPort, redisPassword, redisDB, name)
	return &BufferedQueue{
		Queue:        q,
		BufferSize:   bufferSize,
		Buffer:       make(chan *Package, bufferSize*2),
		flushStatus:  make(chan chan bool, 1),
		flushCommand: make(chan bool, bufferSize*2),
	}
}

// SelectBufferedQueue returns a BufferedQueue if a queue with the name exists
func SelectBufferedQueue(redisHost, redisPort, redisPassword string, redisDB int64, name string, bufferSize int) (queue *BufferedQueue, err error) {
	q, err := SelectQueue(redisHost, redisPort, redisPassword, redisDB, name)
	if err != nil {
		return nil, err
	}

	return &BufferedQueue{
		Queue:        q,
		BufferSize:   bufferSize,
		Buffer:       make(chan *Package, bufferSize*2),
		flushStatus:  make(chan chan bool, 1),
		flushCommand: make(chan bool, bufferSize*2),
	}, nil
}

// Start dispatches the background writer that flushes the buffer.
// If there is already a BufferedQueue running it will return an error.
func (queue *BufferedQueue) Start() error {
	queue.redisClient.SAdd(masterQueueKey(), queue.Name)
	val := queue.redisClient.Get(queueHeartbeatKey(queue.Name)).Val()
	if val == "ping" {
		return fmt.Errorf("buffered queue with this name is already started")
	}
	queue.startHeartbeat()
	queue.startWritingBufferToRedis()
	queue.startPacemaker()
	return nil
}

// FlushBuffer tells the background writer to flush the buffer to redis
func (queue *BufferedQueue) FlushBuffer() {
	flushing := make(chan bool, 1)
	queue.flushStatus <- flushing
	queue.flushCommand <- true
	<-flushing
	return
}

func (queue *BufferedQueue) startWritingBufferToRedis() {
	go func() {
		queue.nextWrite = time.Now().Unix()
		for {
			if len(queue.Buffer) >= queue.BufferSize || time.Now().Unix() >= queue.nextWrite {
				size := len(queue.Buffer)
				a := []string{}
				for i := 0; i < size; i++ {
					p := <-queue.Buffer // 将buffer chan 中的每个package 读取出来，写到a中，再一次性的进行LPush到队列
					a = append(a, p.getString())
				}
				queue.redisClient.LPush(queueInputKey(queue.Name), a...)
				queue.incrRate(queueInputRateKey(queue.Name), int64(size))
				for i := 0; i < len(queue.flushStatus); i++ {
					c := <-queue.flushStatus
					c <- true
				}
				queue.nextWrite = time.Now().Unix() + 1
			}
			<-queue.flushCommand // 阻塞在这里,直到进行 FlushBuffer()或者
		}
	}()
}

func (queue *BufferedQueue) startPacemaker() {
	go func() {
		for {
			queue.flushCommand <- true // 每个 10微秒进行一次 flush命令，但是会被阻塞，直到有<-queue.flushCommand 命令，然后sleep 10Millisecond，继续下一次
			time.Sleep(10 * time.Millisecond)
		}
	}()
}

func (queue *BufferedQueue) startHeartbeat() {
	firstWrite := make(chan bool, 1) // chan 有size时，不会阻塞往里写，直到size满了，当为空的时候，read仍然会阻塞
	go func() {
		firstRun := true
		for {
			queue.redisClient.Set(queueHeartbeatKey(queue.Name), "ping", time.Second)
			if firstRun {
				firstWrite <- true
				firstRun = false
			}
			time.Sleep(500 * time.Millisecond)
		}
	}()
	<-firstWrite //没有明白这个的作用
	return
}
