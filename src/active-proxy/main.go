package main

import (
	"flag"
	"fmt"
)

func main() {
	address, etcdURLS := parseFlags()
	p := NewProxy(address)
	w := NewWatcher(etcdURLS)

	w.StartApplication(p)
	p.Start()
}

func parseFlags() (string, string) {
	port := flag.String("port", "8080", "Port where the proxy listens")
	host := flag.String("host", "localhost", "Host where the proxy is binded")
	etcdURLS := flag.String("etcdURLS", "http://127.0.0.1:4001", "URL(s) where the etcd daemon listens to separated by commas")

	flag.Parse()

	return fmt.Sprintf("%s:%s", *host, *port), *etcdURLS
}
