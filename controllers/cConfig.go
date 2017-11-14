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

func SetViewType(rich bool) {
	RICH_VIEW = rich
}

var (
	RICH_VIEW = true
)

func IsRichView() bool {
	if false {
		rich, err := beego.AppConfig.Bool("rich_view")
		if err != nil {
			rich = false
		}
		beego.Warn("Rich display: ", rich)
		return rich
	} else {
		return RICH_VIEW
	}

}
