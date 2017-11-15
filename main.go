package main

import (
	"flag"

	"github.com/Bulusideng/go-crm/controllers"
	"github.com/Bulusideng/go-crm/models"
	"github.com/Bulusideng/go-crm/routers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
)

func main() {
	rich_view := flag.Bool("rich_view", true, "using rich display")
	flag.Parse()

	models.Register()
	orm.Debug = false
	orm.RunSyncdb("default", false, true)

	controllers.SetViewType(*rich_view)
	models.Init()
	routers.Register()
	//logs.SetLogger(logs.AdapterConsole, `{"level":3}`)
	logs.SetLogger(logs.AdapterFile, `{"filename":"log.log","level":7}`)
	beego.Run()
}
