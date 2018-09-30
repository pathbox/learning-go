import (
	"bufio"
	"encoding/binary"
	"io"
	"net"
)

// 从conn 中读取数据
func ReadPacket(conn net.Conn) ([]byte, error) {
	var head [2]byte // 用于先行 读取packet的长度size

	if _, err := io.ReadFull(conn, head[:]); err != nil {
		return nil, err
	}

	size := binary.BigEndian.Uint16(head) // 将 head值 转为 uint16

	packet := make([]byte, size) // 知道长度了，定义size大小的[]byte 用于正式读取数据

	if _, err := io.ReadFull(conn, packet); err != nil { // 读取这次的整个packet
		return nil, err
	}

	return packet, err

}

// 使用bufio优化版

type PacketConn struct {
	net.Conn
	reader *bufio.Reader
}

func NewPacketConn(conn net.Conn) *PacketConn {
	return &PacketConn{conn, bufio.NewReader(conn)} // conn 被bufio包裹处理
}

func (pconn *PacketConn) ReadPacket() []byte {

	var head [2]byte

	if _, err := io.ReadFull(pconn.reader, head[:]); err != nil {
		return err
	}

	size := binary.BigEndian.Uint16(head) // 将byte 转为 uint16
	packet := make([]byte, size)

	if _, err := io.ReadFull(conn.reader, packet); err != nil {
		return err
	}

	return packet

}

func (conn *PacketConn) Read(p []byte) (int, error) {
	return conn.reader.Read(p)
}