package clientpool

import (
	"fmt"
	config "config"
	"sync"
	"sync/atomic"
	"time"

	iam "proto/iam"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/keepalive"
)

var (
	ErrConnShutdown = errors.New("grpc conn shutdown")

	defaultClientPoolCap    = 5
	defaultKeepAlive        = 30 * time.Second
	defaultKeepAliveTimeout = 10 * time.Second
	DefaultClientTimeout    = 10 * time.Second

	IAMIdentityClientPool *ClientPool
)

type ClientPool struct {
	option   *ClientOption
	capacity int64
	next     int64
	addr     string
	sync.Mutex
	conns []*grpc.ClientConn
}

type ClientOption struct {
	KeepAlive        time.Duration
	KeepAliveTimeout time.Duration
	ClientPoolSize   int
}

// InitClientPool 初始化客户端连接池
func InitClientPool() {
	option := NewDefaultClientOption()
	addr := fmt.Sprintf("%s:%s", config.DataSync.GetString("iam_identity.host"), config.DataSync.GetString("iam_identity.port"))
	IAMIdentityClientPool = NewClientPool(addr, option).init()
}

// NewDefaultClientOption NewDefaultClientOption
func NewDefaultClientOption() *ClientOption {
	return &ClientOption{
		KeepAlive:        defaultKeepAlive,
		KeepAliveTimeout: defaultKeepAliveTimeout,
	}
}

// NewClientPool 新建客户端连接池
func NewClientPool(addr string, option *ClientOption) *ClientPool {
	if option.ClientPoolSize <= 0 {
		option.ClientPoolSize = defaultClientPoolCap
	}

	return &ClientPool{
		addr:     addr,
		conns:    make([]*grpc.ClientConn, option.ClientPoolSize),
		capacity: int64(option.ClientPoolSize),
		option:   option,
	}
}

func (cp *ClientPool) getConn() (*grpc.ClientConn, error) {
	var (
		idx  int64
		next int64
		err  error
	)
	next = atomic.AddInt64(&cp.next, 1)
	idx = next % cp.capacity
	conn := cp.conns[idx]
	if conn != nil && cp.checkState(conn) == nil {
		return conn, nil
	}
	if conn != nil {
		conn.Close()
	}

	// double check, already inited
	conn = cp.conns[idx]
	if conn != nil && cp.checkState(conn) == nil {
		return conn, nil
	}

	conn, err = cp.connect()
	if err != nil {
		return nil, err
	}
	cp.conns[idx] = conn
	return conn, nil
}

func (cp *ClientPool) checkState(conn *grpc.ClientConn) error {
	state := conn.GetState()
	switch state {
	case connectivity.TransientFailure, connectivity.Shutdown:
		return ErrConnShutdown
	}
	return nil
}

func (cp *ClientPool) init() *ClientPool {
	for idx, _ := range cp.conns {
		Zap.Info("ClientPool init", zap.Int("Number", idx))
		conn, _ := cp.connect()
		cp.conns[idx] = conn
	}
	return cp
}

func (cp *ClientPool) connect() (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(cp.addr,
		grpc.WithInsecure(),
		grpc.WithBlock(),
		grpc.WithKeepaliveParams(keepalive.ClientParameters{
			Time:    cp.option.KeepAlive,
			Timeout: cp.option.KeepAliveTimeout},
		),
	)
	if err != nil {
		Zap.Error("client connect Error", zap.String("Error", err.Error()))
		return nil, err
	}

	return conn, nil
}

// Close 关闭连接
func (cp *ClientPool) Close() {
	cp.Lock()
	defer cp.Unlock()

	for _, conn := range cp.conns {
		if conn == nil {
			continue
		}

		conn.Close()
	}
}

// Get iam-identity client
func GetIAMIdentityClient() (iam.IAMIdentityClient, error) {
	conn, err := IAMIdentityClientPool.getConn()
	if err != nil {
		return nil, err
	}
	return iam.NewIAMIdentityClient(conn), nil
}
