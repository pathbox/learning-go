// http://www.jb51.net/article/89321.htm

package main

import (
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"io"
	"log"
	"net"
	"strings"
)

func main() {
	ln, err := net.Listen("tcp", "9090") // 底层实际就是 tcp Listen
	if err != nil {
		log.Panic(err)
	}

	for {
		conn, err := ln.Accept() // Accept() 操作
		if err != nil {
			log.Println("Accept err:", err)
		}
		for {
			handleConnection(conn)
		}
	}
}

func handleConnection(conn net.Conn) {
	content := make([]byte, 1024) // 1024 size 缓存
	_, err := conn.Read(content)
	log.Println(string(content))

	if err != nil {
		log.Println(err)
	}

	isHttp := false
	if string(content[0:3]) == "GET" {
		isHttp = true
	}
	log.Println("isHttp:", isHttp)
	if isHttp {
		headers := parseHandshake(string(content))
		log.Println("headers", headers)
		secWebsocketKey := headers["Sec-WebSocket-Key"]
		guid := "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
		// 计算Sec-WebSocket-Accept
		h := sha1.New()
		log.Println("accept raw:", secWebsocketKey+guid)
		io.WriteString(h, secWebsocketKey+guid)
		accept := make([]byte, 28)
		base64.StdEncoding.Encode(accept, h.Sum(nil)) // sha1后，进行base64，然后存到accept
		log.Println(string(accept))
		response := "HTTP/1.1 101 Switching Protocols\r\n"
		response = response + "Sec-WebSocket-Accept: " + string(accept) + "\r\n"
		response = response + "Connection: Upgrade\r\n"
		response = response + "Upgrade: websocket\r\n\r\n"
		log.Println("response:", response)
		if length, err := conn.Write([]byte(response)); err != nil {
			log.Println(err)
		} else {
			log.Println("send len:", length)
		}
		wsSocket := NewWsSocket(conn)
		for {
			data, err := wsSocket.ReadIframe()
			if err != nil {
				log.Println("readIframe err:", err)
			}
			log.Println("read data:", string(data))
			err = wsSocket.SendIframe([]byte("good"))
			if err != nil {
				log.Println("sendIframe err:", err)
			}
			log.Println("send data")
		}
	} else {
		log.Println(string(content))
	}
}

type WsSocket struct {
	MaskingKey []byte
	Conn       net.Conn
}

func NewWsSocket(conn net.Conn) *WsSocket {
	return &WsSocket{Conn: conn}
}

func (ws *WsSocket) SendIframe(data []byte) error {
	if len(data) > 125 {
		return errors.New("send iframe data error")
	}
	length := len(data)
	maskedData := make([]byte, length)
	for i := 0; i < length; i++ { // 一个字节一个字节的处理 data
		if ws.MaskingKey != nil {
			maskedData[i] = data[i] ^ ws.MaskingKey[i%4]
		} else {
			maskedData[i] = data[i]
		}
	}

	ws.Conn.Write([]byte{0x81})
	var payLenByte byte
	if ws.MaskingKey != nil && len(ws.MaskingKey) != 4 {
		payLenByte = byte(0x80) | byte(length)
		ws.Conn.Write([]byte{payLenByte})
		ws.Conn.Write(ws.MaskingKey)
	} else {
		payLenByte = byte(0x00) | byte(length)
		ws.Conn.Write([]byte{payLenByte})
	}
	ws.Conn.Write(data)
	return nil
}

func (ws *WsSocket) ReadIframe() (data []byte, err error) {
	err = nil
	//第一个字节：FIN + RSV1-3 + OPCODE

	opcodeByte := make([]byte, 1)
	ws.Conn.Read(opcodeByte)
	FIN := opcodeByte[0] >> 7
	RSV1 := opcodeByte[0] >> 6 & 1
	RSV2 := opcodeByte[0] >> 5 & 1
	RSV3 := opcodeByte[0] >> 4 & 1
	OPCODE := opcodeByte[0] & 15
	log.Println(RSV1, RSV2, RSV3, OPCODE)
	payloadLenByte := make([]byte, 1)
	ws.Conn.Read(payloadLenByte)
	payloadLen := int(payloadLenByte[0] & 0x7F)
	mask := payloadLenByte[0] >> 7
	if payloadLen == 127 {
		extendedByte := make([]byte, 8)
		ws.Conn.Read(extendedByte)
	}
	maskingByte := make([]byte, 4)
	if mask == 1 {
		ws.Conn.Read(maskingByte)
		ws.MaskingKey = maskingByte
	}
	payloadDataByte := make([]byte, payloadLen)
	ws.Conn.Read(payloadDataByte) // 将数据读取存到payloadDataByte
	log.Println("data:", payloadDataByte)
	dataByte := make([]byte, payloadLen)
	for i := 0; i < payloadLen; i++ {
		if mask == 1 {
			dataByte[i] = payloadDataByte[i] ^ maskingByte[i%4]
		} else {
			dataByte[i] = payloadDataByte[i]
		}
	}
	if FIN == 1 {
		data = dataByte
		return
	}
	nextData, err := ws.ReadIframe()
	if err != nil {
		return
	}
	data = append(data, nextData...)
	return
}

func parseHandshake(content string) map[string]string {
	headers := make(map[string]string, 10)
	lines := strings.Split(content, "\r\n")
	for _, line := range lines {
		if len(line) > 0 {
			words := strings.Split(line, ":") // key:value
			if len(words) == 2 {
				headers[strings.Trim(words[0], " ")] = strings.Trim(words[1], " ")
			}
		}
	}
}

/*
客户端发送消息：

GET /chat HTTP/1.1
  Host: server.example.com
  Upgrade: websocket
  Connection: Upgrade
  Sec-WebSocket-Key: dGhlIHNhbXBsZSBub25jZQ==
  Origin: http://example.com
  Sec-WebSocket-Version: 13
服务端返回消息：

HTTP/1.1 101 Switching Protocols
Upgrade: websocket
Connection: Upgrade
Sec-WebSocket-Accept: s3pPLMBiTxaQ9kYGzzhZRbK+xOo=
*/
