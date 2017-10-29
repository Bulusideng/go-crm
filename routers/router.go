package routers

import (
	"github.com/Bulusideng/go-crm/controllers"
	"github.com/astaxie/beego"
)

func Register() {
	beego.Router("/", &controllers.MainController{})
	if false {
		beego.Router("contract/*", &controllers.ContractController{})
		beego.Router("contract/add", &controllers.ContractController{}, "get:Add")
		beego.Router("contract/update", &controllers.ContractController{}, "get:Update")
		beego.Router("contract/query", &controllers.ContractController{}, "post:Query")

	} else {
		beego.Router("contract/", &controllers.ContractController{})
		beego.AutoRouter(&controllers.ContractController{})
	}

	beego.Router("user/", &controllers.UserController{})
	beego.AutoRouter(&controllers.UserController{})

	beego.Router("/login", &controllers.LoginController{})

}
