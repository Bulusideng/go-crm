package controllers

import (
	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type ConfigController struct {
	baseController
}

func (this *ConfigController) Get() {
	curUser := GetCurAcct(this.Ctx)
	if !curUser.IsAdmin() {
		this.Redirect("/contract", 302)
		return
	}

	this.Data["SysConfig"] = true
	this.TplName = "config.html"
	this.Data["CurUser"] = curUser
	this.Data["Config"] = models.GetConfig()
}

func (this *ConfigController) Post() {
	for k, v := range this.Ctx.Request.Form {
		beego.Debug("Param ", k, ":", v)
	}

	curUser := GetCurAcct(this.Ctx)
	if !curUser.IsAdmin() {
		this.Redirect("/contract", 302)
		return
	}
	cfg := models.Config{}
	err := this.ParseForm(&cfg)
	if err == nil {
		err = models.UpdateConfig(&cfg)
	}
	if err == nil {
		this.RedirectTo("status", "设置成功！", "/config", 302)
	} else {
		this.RedirectTo("status", err.Error(), "/config", 302)
	}
}
