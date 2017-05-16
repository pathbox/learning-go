package main

import (
	"github.com/astaxie/beego"
)

type HomeController struct {
  beego.Controller
}

func (this *HomeController) Get(){
  this.Ctx.WriteString("hello world")
}

func main() {
  beego.Router("/", &HomeController{})

  beego.Run(":9000")
}
