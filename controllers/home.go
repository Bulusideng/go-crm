package controllers

import (
	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	baseController
}

func (this *MainController) Get() {
	beego.Warn("Remote addr:", this.Ctx.Request.RemoteAddr)
	this.Redirect("/login", 302)
	return
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	curUser := GetCurAcct(this.Ctx)
	this.Data["CurUser"] = curUser
	contracts, err := models.GetContracts(curUser, nil)
	if err != nil {
		beego.Error(err)
		return
	} else {
		this.Data["Contracts"] = contracts
	}
}
