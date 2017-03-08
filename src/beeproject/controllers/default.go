package controllers

import (
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	c.Data["Website"] = "beego.me"
	c.Data["Email"] = "astaxie@gmail.com"
	c.TplName = "index.tpl"
}

//
//beego.AppConfig.String("mysqluser")
//beego.AppConfig.String("mysqlpass")
//beego.AppConfig.String("mysqlurls")
//beego.AppConfig.String("mysqldb")
