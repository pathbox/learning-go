package main

import (
	"flag"

	"github.com/golang/glog"
)

func main() {
	flag.Parse()

	glog.Info("hello, glog")
	glog.Warning("Warning glog")
	glog.Error("error glog")

	glog.Infof("info %d", 1)
	glog.Warning("warning %d", 2)
	glog.Error("error %d", 3)

	glog.V(3).Infoln("info with v 3")
	glog.V(2).Infoln("info with v 2")
	glog.V(1).Infoln("info with v 1")
	glog.V(0).Infoln("info with v 0")

	// 退出时调用，确保日志写入文件中
	glog.Flush()
}
