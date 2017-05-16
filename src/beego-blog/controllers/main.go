package controllers

import (
	"github.com/ulricqin/beego-blog/g"
	"github.com/ulricqin/beego-blog/models/blog"
	"github.com/ulricqin/beego-blog/models/catalog"
)

type MainController struct {
	BaseController
}

func (this *MainController) Get() {
	this.Data["Catalogs"] = catalog.All()
	this.Data["PageTitle"] = "首页"
	this.Layout = "layout/default.html"
	this.TplName = "index.html"
}

func (this *MainController) Read() {
	ident := this.GetString(":ident")
	b := blog.OneByIdent(ident)
}
