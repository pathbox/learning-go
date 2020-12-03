package main

import (
	"fmt"
	"log"
	"net/http"
	"time"
)

func main() {
	http.HandleFunc("/", SayHelloContext) // 设置访问的路由

	log.Fatalln(http.ListenAndServe(":8080", nil))
}

// 这里我假定请求需要耗时2s，在请求2s后返回，我们期望监控goroutine在打印2次Current request is in progress后即停止。但运行发现，监控goroutine打印2次后，其仍不会结束，而会一直打印下去
func SayHello(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(&request)

	go func() {
		for range time.Tick(time.Second) {
			fmt.Println("Current request is in progress")
		}
	}()

	time.Sleep(2 * time.Second)
	writer.Write([]byte("Hi"))
}

func SayHelloContext(writer http.ResponseWriter, request *http.Request) {
	fmt.Println(&request)

	go func() {
		for range time.Tick(time.Second) {
			select {
			case <-request.Context().Done():
				fmt.Println("request is outgoing")
				return
			default:
				fmt.Println("Current request is in progress")
			}
		}
	}()

	time.Sleep(2 * time.Second)
	writer.Write([]byte("Hi"))
}

/*
https://mp.weixin.qq.com/s/JKMHUpwXzLoSzWt_ElptFg
 context包可以提供一个请求从API请求边界到各goroutine的请求域数据传递、取消信号及截至时间等能力
 Context 的主要作用就是在不同的 Goroutine 之间同步请求特定的数据、消信号以及处理请求的截止日期

type Context interface {
 Deadline() (deadline time.Time, ok bool)
 Done() <-chan struct{}
 Err() error
 Value(key interface{}) interface{}
}

func WithCancel(parent Context) (ctx Context, cancel CancelFunc)
func WithDeadline(parent Context, deadline time.Time) (Context, CancelFunc)
func WithTimeout(parent Context, timeout time.Duration) (Context, CancelFunc)
func WithValue(parent Context, key, val interface{}) Context

ithCancel函数，传递一个父Context作为参数，返回子Context，以及一个取消函数用来取消Context。WithDeadline函数，和WithCancel差不多，它会多传递一个截止时间参数，意味着到了这个时间点，会自动取消Context，当然我们也可以不等到这个时候，可以提前通过取消函数进行取消。

WithTimeout和WithDeadline基本上一样，这个表示是超时自动取消，是多少时间后自动取消Context的意思。

WithValue函数和取消Context无关，它是为了生成一个绑定了一个键值对数据的Context，这个绑定的数据可以通过Context.Value方法访问到

此函数接收 context 并返回派生 context，其中值 val 与 key 关联，并通过 context 树与 context 一起传递。这意味着一旦获得带有值的 context，从中派生的任何 context 都会获得此值。不建议使用 context 值传递关键参数，而是函数应接收签名中的那些值，使其显式化。

func main()  {
 ctx,cancel := context.WithCancel(context.Background())
 defer cancel()
 go Speak(ctx)
 time.Sleep(10*time.Second)
}

func Speak(ctx context.Context)  {
 for range time.Tick(time.Second){
  select {
  case <- ctx.Done():
   return
  default:
   fmt.Println("balabalabalabala")
  }
 }
}

func main()  {
 now := time.Now()
 later,_:=time.ParseDuration("10s")

 ctx,cancel := context.WithDeadline(context.Background(),now.Add(later))
 defer cancel()
 go Monitor(ctx)

 time.Sleep(20 * time.Second)

}

func Monitor(ctx context.Context)  {
 select {
 case <- ctx.Done():
  fmt.Println(ctx.Err())
 case <-time.After(20*time.Second):
  fmt.Println("stop monitor")
 }
}

不要把Context放在结构体中，要以参数的方式传递。

以Context作为参数的函数方法，应该把Context作为第一个参数，放在第一位。

给一个函数方法传递Context的时候，不要传递nil，如果不知道传递什么，就使用context.TODO

Context的Value相关方法应该传递必须的数据，不要什么数据都使用这个传递。context.Value 应该很少使用，它不应该被用来传递可选参数。这使得 API 隐式的并且可以引起错误。取而代之的是，这些值应该作为参数传递。

Context是线程安全的，可以放心的在多个goroutine中传递。同一个Context可以传给使用其的多个goroutine，且Context可被多个goroutine同时安全访问。

Context 结构没有取消方法，因为只有派生 context 的函数才应该取消 context
*/
