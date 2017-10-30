package controllers

import (
	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (this *MainController) Get() {
	this.Redirect("/contract", 302)
	return
	this.Data["IsHome"] = true
	this.TplName = "home.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	contracts, err := models.GetAllContracts()
	if err != nil {
		beego.Error(err)
		return
	} else {
		this.Data["Contracts"] = contracts
	}
}
