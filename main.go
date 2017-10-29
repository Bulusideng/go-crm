package main

import (
	"github.com/Bulusideng/go-crm/models"
	"github.com/Bulusideng/go-crm/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func main() {
	models.Register()
	orm.Debug = false
	orm.RunSyncdb("default", false, true)
	models.RegisterAdmin()
	routers.Register()
	//models.RegCase()
	beego.Run()
}
