package controllers

//"github.com/Bulusideng/go-crm/models"
import "github.com/astaxie/beego"

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
	this.TplName = "home.html"
	this.Data["CurUser"] = curUser
}

func IsRichView() bool {
	rich, err := beego.AppConfig.Bool("rich_view")
	if err != nil {
		rich = false
	}
	beego.Warn("Rich display: ", rich)
	return rich
}
