package pingpong

import "github.com/leesper/tao"

const (
	PingPontMessage int32 = 1
)

type Message struct {
	Info string
}

func (pp Message) MessageNumber() int32 {
	return PingPontMessage
}

func (pp Message) Serialize() ([]byte, error) {
	return []byte(pp.Info), nil
}

func DeserializeMessage(data []byte) (mseeage tao.Message, err error) {
	if data == nil {
		return nil, tao.ErrNilData
	}
	info := string(data)
	msg := Message{
		Info: info,
	}
	return msg, nil
}

// func ProcessPingPongMessage(ctx tao.Context, conn tao.Connection) {
//   if serverConn, ok := conn.(*tao.ServerConnection); ok {
//     if serverConn.GetOwner() != nil {
//       connections := serverConn.GetOwner().GetAllConnections()
//       for v := range connections.IterValues() {
//         c := v.(tao.Connection)
//         c.Write(ctx.Message())
//       }
//     }
//   }
// }
