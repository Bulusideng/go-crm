package controllers

import (
	"github.com/astaxie/beego"
)

type StatusController struct {
	beego.Controller
}

func (this *StatusController) Get() {

	this.TplName = "status.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["Message"] = this.GetString("msg")
}
