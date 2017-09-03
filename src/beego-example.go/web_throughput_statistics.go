package main

import (
	"encoding/json"
	"fmt"
	"log"
	// "strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
)

var globalSessions *session.Manager

type Index struct {
	beego.Controller
}

func init() {
	config := fmt.Sprintf(`{"cookieName":"gosessionid","gclifetime":%d,"enableSetCookie":true}`, 3600*24)
	conf := new(session.ManagerConfig)
	if err := json.Unmarshal([]byte(config), conf); err != nil {
		log.Fatal("json decode error", err)
	}

	globalSessions, _ := session.NewManager("memory", conf)
	go globalSessions.GC()
}

func main() {
	beego.BConfig.Listen.HTTPPort = 3000
	// beego.SetLogger("file", `{filename: "logs/stdout.log"}`)
	beego.BConfig.AppName = "巧克力糖"
	beego.BConfig.WebConfig.Session.SessionGCMaxLifetime = 3600 * 24
	beego.BConfig.WebConfig.Session.SessionCookieLifeTime = 3600 * 24

	beego.Router("/*", &Index{}, "*:Count")
	go beego.Run()

	for {
		time.Sleep(10 * time.Second)
	}
}

func (this *Index) Count() {
	path_url := this.Ctx.Request.URL.String()
	fmt.Println("get url:", path_url)
	if path_url == "/favicon.ico" { //忽略此路由地址请求
		this.Ctx.WriteString("")
		this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")
		return
	}

	fmt.Printf("===%v===\n", this.Ctx.Request)

	Client_Host := this.Ctx.Request.Host
	Client_Method := this.Ctx.Request.Method
	Client_User_Agent := this.Ctx.Request.Header.Get("User-Agent")
	Client_IP := this.Ctx.Request.Header.Get("Remote_addr")  //客户端IP
	Client_Referer := this.Ctx.Request.Header.Get("Referer") //来源
	if len(Client_IP) <= 7 {
		Client_IP = this.Ctx.Request.RemoteAddr //获取客户端IP
	}
	if strings.Contains(Client_IP, ":") {
		ip_boolA, ip_dataA := For_IP(string(Client_IP)) //获取IP
		if ip_boolA {
			Client_IP = ip_dataA
		}
	}
	// this.Ctx.ResponseWriter.Header().Set("Content-Type", "text/html")

	this.Ctx.WriteString(fmt.Sprintf("=====客户端IP:%v======</br>\n", Client_IP))
	this.Ctx.WriteString(fmt.Sprintf("=====访问域名:%v======</br>\n", Client_Host))
	this.Ctx.WriteString(fmt.Sprintf("=====请求路径:%v======</br>\n", path_url))
	this.Ctx.WriteString(fmt.Sprintf("=====来源来路:%v======</br>\n", Client_Referer))
	this.Ctx.WriteString(fmt.Sprintf("=====请求方式:%v======</br>\n", Client_Method))
	this.Ctx.WriteString(fmt.Sprintf("=====请求头:%v======</br>\n", Client_User_Agent))
	this.Ctx.WriteString(fmt.Sprintf("=====访问次数:%v======</br>\n", this.Cookie_session()))
	return

}

func For_IP(valuex string) (bool, string) {
	data_list := strings.Split(valuex, ":")
	if len(data_list) >= 2 {
		return true, data_list[0]
	}
	return false, ""
}

func (this *Index) Cookie_session() int { //id统计  PV  这样统计只能针对单个浏览器有效
	pv := 0
	=====================
	Cookie 统计法
	cook := this.Ctx.GetCookie("countnum") //获取Cookie
	if cook == "" {
		this.Ctx.SetCookie("countnum", "1", "/")
		pv = 1
	} else {
		xx, err := strconv.Atoi(cook)
		if err == nil {
			pv = xx + 1
			this.Ctx.SetCookie("countnum", strconv.Itoa(pv), "/")
		} else {
			pv = 0
		}
	}
	// return pv
	//=====================
	//session 统计法
	sess, _ := globalSessions.SessionStart(this.Ctx.ResponseWriter, this.Ctx.Request)
	ct := sess.Get("countnum")
	if ct == nil {
		sess.Set("countnum", 1)
		pv = 1
	} else {
		pv = ct.(int) + 1
		sess.Set("countnum", pv)
	}
	return pv
}
