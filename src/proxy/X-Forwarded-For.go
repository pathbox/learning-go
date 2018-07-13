/*X-Forwarded-For: client, proxy1, proxy2

最终客户端或者服务器端接受的请求，X-Forwarded-For 是没有最邻近节点的 ip 地址的，而这个地址可以通过 remote address 获取
每个节点（不管是客户端、代理服务器、真实服务器）都可以随便更改 X-Forwarded-For 的值，因此这个字段只能作为参考
*/

package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"
)

type Pxy struct{}

func (p *Pxy) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	fmt.Printf("Received request %s %s %s\n", req.Method, req.Host, req.RemoteAddr)

	// step 1
	outReq, _ := http.NewRequest("GET", "http://www.baidu.com", nil)

	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := outReq.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		outReq.Header.Set("X-Forwarded-For", clientIP)
		log.Println("X-Forwarded-For", outReq.Header)
	}
	// step 2

	res, err := http.DefaultClient.Do(outReq)
	if err != nil {
		log.Println(err)
		rw.WriteHeader(http.StatusBadGateway)
		return
	}
	for key, value := range res.Header {
		for _, v := range value {
			rw.Header().Add(key, v)
		}
	}
	// log.Println("response:", res, "body:", res.Body)
	b, _ := ioutil.ReadAll(res.Body) // res.Body just be read one time
	log.Println("Size: ", len(b))
	// step 3

	sio := strings.NewReader(string(b))

	rw.WriteHeader(res.StatusCode)
	io.Copy(rw, sio)
	defer res.Body.Close()

}

func main() {
	log.Println("Serve on :9090")
	http.Handle("/", &Pxy{}) // Handler 是一个接口，接口方法是 ServeHTTP(rw http.ResponseWriter, req *http.Request)，而Pxy 实现了这个方法
	http.ListenAndServe(":9090", nil)
}
