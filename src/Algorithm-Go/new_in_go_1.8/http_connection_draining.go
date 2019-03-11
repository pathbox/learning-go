/*
这个issue请求实现http.Server的连接耗尽(draining)的功能。现在可以调用srv.Close可以立即停止http.Server,
也可以调用srv.Shutdown(ctx)等待已有的连接处理完毕(耗尽，draining, github.com/tylerb/graceful 的用户应
该熟悉这个特性)。
下面这个例子中，服务器当收到SIGINT信号后(^C)会优雅地关闭。
一旦收到SIGINT信号，服务器会立即停止接受新的连接，srv.ListenAndServe()会返回http.ErrServerClosed。
srv.Shutdown会一直阻塞，直到所有未完成的request都被处理完以及它们底层的连接被关闭
 */

package main

import (
	"context"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {
	// subscribe to SIGINT signals
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		io.WriteString(w, "Finished")
	}))
	srv := &http.Server{Addr: ":9090", Handler: mux}
	go func() {
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
		}
	}()

	<-stopChan // wait for SIGINT
	log.Println("Shutting down server...")
	ctx, _ := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	log.Println("Server gracefully stopped")
}
