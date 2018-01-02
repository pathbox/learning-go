package dworker

import (
	"fmt"
	"time"

	"github.com/coreos/etcd/client"
	"golang.org/x/net/context"
)

var kHeartBeatInterval = time.Second * 2
var kTTL = time.Second * 5
var kRoot = "service"

type Worker struct {
	kapi    client.KeysAPI
	key     string
	extInfo string
	active  bool
	stop    bool
}

func NewWorker(serviceName string, node string, extInfo string, endpoints []string) (*Worker, error) {
	cfg := client.Config{
		Endpoints:               endpoints,
		HeaderTimeoutPerRequest: time.Second * 2,
	}
	c, err := client.New(cfg)
	if err != nil {
		return nil, err
	}
	worker := &Worker{
		kapi:    client.NewKeysAPI(c),
		key:     fmt.Sprintf("%s/%s/%s", kRoot, serviceName, node),
		extInfo: extInfo,
		active:  false,
		stop:    false,
	}
	return worker, nil
}

func (w *Worker) Register() {
	w.heartbeat()
	go w.heartbeatPeriod() // 起一个goroutine 进行心跳
}

func (w *Worker) Unregister() {
	w.stop = true
	/// no need to wait result
}

func (w *Worker) IsActive() bool {
	return w.active
}

func (w *Worker) IsStop() bool {
	return w.stop
}

func (w *Worker) heartbeatPeriod() {
	for !w.stop {
		w.heartbeat()
		time.Sleep(kHeartBeatInterval) // 每个kHeartBeatInterval 进行判断
	}
}

func (w *Worker) heartbeat() error {
	_, err := w.kapi.Set(context.Background(), w.key, w.extInfo, &client.SetOptions{
		TTL: kTTL,
	})
	w.active = err == nil
	return err
}
