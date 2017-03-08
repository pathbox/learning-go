package main

import (
	_ "beeproject/routers"
	"github.com/astaxie/beego"
)

func main() {
	beego.BConfig.Listen.ServerTimeOut=60 //请求超时
	beego.BConfig.Listen.Graceful=false  // 是否开启热升级
	beego.BConfig.Listen.HTTPPort = 8081  // 这里会覆盖在app.conf中的配置
	beego.BConfig.Listen.EnableFcgi = false
	beego.BConfig.Log.AccessLogs = true
	beego.EnableXSRF = true
	beego.XSRFKEY = "61oETzKXQAGaYdkL5gEmGeJJFuYh7EQnp2XdTP1o"
	beego.XSRFExpire = 3600  //过期时间，默认1小时
	beego.Run()
}
