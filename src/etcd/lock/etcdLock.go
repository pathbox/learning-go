package etcdservice

import (
	"log"

	clientv3 "github.com/coreos/etcd/clientv3"
	concurrency "github.com/coreos/etcd/clientv3/concurrency"
	"github.com/kinwyb/go/err1"
)

var (
	LockFail   = err1.NewError(-100, "锁失败")
	UnlockFail = err1.NewError(-101, "解锁失败")
)

type EtcdLockFactory struct {
	etcd *clientv3.Client
}

// 初始化ETCD连接
func NewEtcdLockFactory(endpoints []string) (*EtcdLockFactory, error) {
	ret := &EtcdLockFactory{}
	var err error
	ret.etcd, err = clientv3.New(clientv3.Config{Endpoints: endpoints})
}

//获取一个分布式锁对象。如果出错会painc需要处理panic情况
func (e *EtcdLockFactory) GetLock(path string) (*Lock, error) {
	etcdSession, err := concurrency.NewSession(e.etcd)
	if err != nil {
		return nil, err
	}

	return &Lock{
		mux:    concurrency.NewMutex(etcdSession, path),
		client: e.etcd,
	}, nil
}

//获取etcd连接对象
func (e *EtcdLockFactory) GetEtcdClient() *clientv3.Client {
	return e.etcd
}
func (e *EtcdLockFactory) GetEtcdSession() (*concurrency.Session, error) {
	return concurrency.NewSession(e.etcd)
}

func (e *EtcdLockFactory) Close() {
	if e.etcd != nil {
		e.etcd.Close()
		log.Println("etcd close")
	}
}

type Lock struct {
	mux    *concurrency.Mutex
	client *clientv3.Client
	err    error
	islock bool
}

func (l *Lock) Lock() {
	if l.mux != nil && l.client != nil {
		l.err = l.mux.Lock(l.client.Ctx())
		l.islock = l.err == nil
	}
}

func (l *Lock) Unlock() {
	if l.mux != nil && l.client != nil {
		l.mux.Unlock(l.client.Ctx())
		l.islock = false
	}
}

func (l *Lock) IsLock() bool {
	return l.islock
}

func (l *Lock) LockErr() error {
	return l.err
}

func (l *Lock) CallFunc(fun func() err1.Error) err1.Error {
	if fun == nil {
		return nil
	}
	l.Lock()
	if l.err != nil {
		return LockFail
	}
	err := fun()
	if err != nil {
		l.Unlock() // 报错也要解锁，要不会造成死锁
		return err
	}
	l.Unlock()
	return nil
}
