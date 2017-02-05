package request_center

import (
	"bytes"
	"io"
	"net"
	"runtime"
	"strings"
	"time"
)

type QRequest struct {
	conn net.Conn
}

func CreateRequest(typename string, host string) *QRequest {
	request := new(QRequest)
	runtime.SetFinalizer(request, func(q *QRequest) {
		q.Close()
	})
	conn, err := net.Dial(typename, host)
	if err == nil {
		request.conn = conn
		return request
	}
	return nil
}

func (q *QRequest) Close() {
	q.conn.SetDeadline(time.Now())
	q.conn.Close()
}

func (q *QRequest) SendData(heads []string, body string) bool {
	headsStr := strings.Join(heads, "\r\n") + body
	count, err := q.conn.Write([]byte(headsStr))
	if err != nil || count == 0 {
		return false
	}
	return true
}

func (q *QRequest) ReceiveData(timeout time.Duration) []byte {
	if timeout > 0 {
		q.conn.SetReadDeadline(time.Now().Add(time.Millisecond * timeout))
	}

	var buf bytes.Buffer
	buffer := make([]byte, 8192)
	for {
		sizenew, err := q.conn.Read(buffer)
		buf.Write(buffer[:sizenew])
		if err == io.EOF || sizenew < 8192 { // 读到 EOF 或 最后的字节数的bytes数据后，就停止循环
			break
		}
	}
	return buf.Bytes()
}
