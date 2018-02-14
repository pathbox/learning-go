package main

import (
	"github.com/coreos/go-etcd/etcd"
	"log"
	"strings"
)

type watcher struct {
	etcdLeader   string
	etcdMachines []string
	client       *etcd.Client
}

func NewWatcher(etcdURLS string) *watcher {
	w := &watcher{}
	w.Init(etcdURLS)
	return w
}

func (w *watcher) Init(etcdURLS string) {
	w.client = etcd.NewClient(string.Split(etcdURLS, ","))

	if len(w.etcdMachines) > 0 {
		w.client.SetCluster(w.etcdMachines)
	}
}

func (w *watcher) StartApplications(p *proxy) {
	go w.loadApplications(p)
	go w.watchApplications(p)
}

func (w *watcher) loadApplications(p *proxy) {
	values, err := w.client.Get("applications", true, true)

	if err == nil {
		for _, entry := range values.Node.Nodes {
			app := strings.Split(entry.Key, "/")[2]
			w.registerApp(app, p)
		}
	}
}

func (w *watcher) watchApplications(p *proxy) {
	appsChannel := make(chan *etcd.Response, 10)

	go w.client.Watch("applications", 0, true, appsChannel, nil)
	for entry := range appsChannel {
		app := strings.Split(entry.Node.Key, "/")[2]

		w.registerApp(app, p)
	}
}

func (w *watcher) registerApp(app string, p *proxy) {
	values, err := w.client.Get("applications/"+app, true, true)

	if err != nil {
		log.Printf("Error getting settings for: %s\nReason: %s", app, err.Error())
	} else {
		a := &application{Name: app}

		for _, value := range values.Node.Nodes {
			switch value.Key {
			case "/applications/" + app + "/port":
				a.Port = value.Value
			case "/applications/" + app + "/test":
				a.Test = value.Value
			}
		}

		p.Route(a)
	}
}
