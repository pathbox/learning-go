package socket

import (
	"bytes"
	"encoding/binary"
)

const bitlength = 10 //数据长度占用字节数

var heartBeatBytes = []byte{0x11, 0x22, 0x13, 0x24, 0x15, 0x26, 0x17, 0x28, 0x19, 0x11}

type Protocol struct {
	data       chan []byte    //解析成功的数据
	byteBuffer *bytes.Buffer  //数据存储中心
	dataLength int64          //当前消息数据长度
	heartbeat  []byte         //心跳包数据，如果设置而且接收到改数据会被忽略而不被输出
	Handler    PackageHandler //数据编码handler
}

// chanLength为解析成功数据channel缓冲长度
func NewProtocol(chanLength ...int) *Protocol {
	length := 100
	if chanLength != nil && len(chanLength) > 0 {
		length = chanLength[0]
	}
	return &Protocol{
		data:       make(chan []byte, length),
		byteBuffer: bytes.NewBufferString(""),
	}
}

// Packet 封包
func (p *Protocol) Packet(message []byte) []byte {
	if p.Handler != nil {
		message = p.Handler.Packet(message)
	}
	return append(IntToByte(int64(len(message))), message...)
}

func (p *Protocol) Read() <-chan []byte {
	return p.data
}

//设置心跳包数据内容，如果接收到的一条消息刚好于设置的心跳包
//内容一致,这条消息将会忽略不会进入读取成功的消息队列中
func (p *Protocol) SetHeartBeat(b []byte) {
	p.heartbeat = b
}

//解析成功的数据请用Read方法获取
func (p *Protocol) Unpack(buffer []byte) {
	p.byteBuffer.Write(buffer)
	for { //多条数据循环处理
		length := p.byteBuffer.Len()
		if length < bitlength { //前面8个字节是长度
			return
		}
		p.dataLength = ByteToInt(p.byteBuffer.Bytes()[0:bitlength])
		if int64(length) < p.dataLength+bitlength { //数据长度不够,等待下次读取数据
			return
		}
		data := make([]byte, p.dataLength+bitlength)
		p.byteBuffer.Read(data)
		msg := data[bitlength]
		if p.Handler != nil {
			msg = p.Handler.Unpackage(msg)
		}
		if p.heartbeat != nil && bytes.Equal(msg, p.heartbeat) {
			//对比接收到的内容如果和设置的内容一致忽略该条消息
			continue
		}
		p.data <- msg
	}
}

//重置
func (p *Protocol) Reset() {
	p.dataLength = 0
	p.byteBuffer.Reset() //清空重新开始
}

func IntToByte(length int64) []byte {
	ret := make([]byte, bitlength)
	binary.PutUvarint(ret, length)
	return ret
}

func ByteToInt(data []byte) int64 {
	x, _ := binary.Varint(data)
	return x
}
