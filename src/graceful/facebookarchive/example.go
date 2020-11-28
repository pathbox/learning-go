package main

import (
	"net/http"
	"time"

	"github.com/facebookgo/grace/gracehttp"
)

func main() {
	gracehttp.Serve(
		&http.Server{Addr: ":5001", Handler: newGraceHandler()},
		&http.Server{Addr: ":5002", Handler: newGraceHandler()},
	)
}

func newGraceHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/sleep", func(w http.ResponseWriter, r *http.Request) {
		duration, err := time.ParseDuration(r.FormValue("duration"))
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}
		time.Sleep(duration)
		w.Write([]byte("Hello World"))
	})
	return mux
}

// 旧API不会断掉，会执行原来的逻辑，pid会变化
// curl "http://127.0.0.1:5001/sleep?duration=60s" &

/*
1.
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)
启动服务后，看到服务的PID是62844

2. 执行请求，这个请求会持续60秒，这个请求连接的
curl "http://127.0.0.1:5001/sleep?duration=60s" &

netstat -anlt| grep 5001
tcp4       0      0  127.0.0.1.5001         127.0.0.1.61724        ESTABLISHED
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    4u  IPv6 0x4ce79a57f7870183      0t0  TCP localhost:commplex-link->localhost:61748 (ESTABLISHED) 可以看到这个是curl和服务建立的连接
example 62844 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)

3.
kill -USR2 62844

再执行：
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 62844 pathbox    4u  IPv6 0x4ce79a57f7870183      0t0  TCP localhost:commplex-link->localhost:61748 (ESTABLISHED)
curl    65640 pathbox    3u  IPv4 0x4ce79a5804a1230b      0t0  TCP localhost:61748->localhost:commplex-link (ESTABLISHED)
example 65870 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)
example 65870 pathbox    4u  IPv6 0x4ce79a57ef2f3183      0t0  TCP localhost:commplex-link->localhost:61779 (ESTABLISHED)
curl    67548 pathbox    3u  IPv4 0x4ce79a57f93cb153      0t0  TCP localhost:61779->localhost:commplex-link (ESTABLISHED)
kill -USR2已经执行了，看到新的LISTEN已经创建，PID是65870，但是curl 的连接还在一直还在执行，并没有因为kill操作而断了。
这时候再执行一个新的curl请求，新的请求：localhost:61779 连接是和新的服务进程PID65870 建立了连接，请的请求请求到了重启后的服务进程上，不会请求到旧的服务进程上。旧的服务进程等原有的请求执行完了，就释放了

4. 最后留下最新的服务进程
lsof -i:5001
COMMAND   PID    USER   FD   TYPE             DEVICE SIZE/OFF NODE NAME
example 65870 pathbox    8u  IPv6 0x4ce79a57ef2f44a3      0t0  TCP *:commplex-link (LISTEN)

总结： grace重启后，旧的服务进程和新的服务进程同时存在，旧的服务进程在其原有的连接请求执行完后释放，新的请求会请求到新的服务进程进行处理
*/
