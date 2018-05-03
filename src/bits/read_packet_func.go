import (
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