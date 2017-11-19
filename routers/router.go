package routers

import (
	"os"

	"github.com/Bulusideng/go-crm/controllers"
	"github.com/astaxie/beego"
)

func Register() {
	beego.Router("/", &controllers.MainController{})

	beego.Router("contract/", &controllers.ContractController{})
	beego.AutoRouter(&controllers.ContractController{})

	beego.Router("account/", &controllers.AccountController{})
	beego.AutoRouter(&controllers.AccountController{})

	beego.Router("/config", &controllers.ConfigController{})
	beego.Router("/status", &controllers.StatusController{})

	beego.Router("/login", &controllers.LoginController{})

	os.Mkdir("files", os.ModePerm)
	beego.Router("/files/*", &controllers.FileController{})

	beego.Router("/attachment/*", &controllers.AttachmentController{})
	beego.AutoRouter(&controllers.AttachmentController{})

}
