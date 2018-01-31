// 一个用户对应一个连接
import (
	"net"
	"time"
)

type User struct {
	uid               int64
	conn              *MsgConn
	BKicked           bool // 被另外登陆的一方踢下线
	BHeartBeatTimeout bool // 心跳超时
}

type MsgConn struct {
	conn       net.Conn
	lastTick   time.Time // 上次接收到包时间
	remoteAddr string    // 为每个连接创建一个唯一标识符
	user       *User     // MsgConn与User一一映射
}

func ListenAndServe(network, address string) {
	tcpAddr, err := net.ResolveTCPAddr(network, address)
	if err != nil {
		logger.Fatalf(nil, "ResolveTcpAddr err:%v", err)
	}
	listener, err = net.ListenTCP(network, tcpAddr)
	if err != nil {
		logger.Fatalf(nil, "ListenTCP err:%v", err)
	}
	go accept()
}

func accept() {
	for {
		conn, err := listener.AcceptTCP()
		if err == nil {

			// 包计数，用来限制频率

			//anti-attack， 黑白名单
          ...

            // 新建一个连接
			imconn := NewMsgConn(conn)

			// run
			imconn.Run()
		}
	}
}

func (conn *MsgConn) Run() {
	conn.onConnect()

	go func() {
		tickerRecv := time.NewTicker(time.Second * time.Duration(rateStatInterval))
		for {
			select {
			case <-conn.stopChan:
				ickerRecv.Stop()
				return
			case <-tickerRecv.C:
				conn.packetsRecv = 0
			default:
			   // 在 conn.parseAndHandlePdu 里面通过Golang本身的io库里面提供的方法读取数据，如io.ReadFull
				conn_closed := conn.parseAndHandlePdu()
				if conn_closed {
					tickerRecv.Stop()
					return
				}
			}
		}
	}()
}

// 将 user 和 conn 一一对应起来
func (conn *MsgConn) onConnect() *User {
	user := &User{conn: conn, durationLevel: 0, startTime: time.Now(), ackWaitMsgIdSet: make(map[int64]struct{})}
	conn.user = user
	return user
}
