package controllers

import (
	"fmt"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	beego.Controller
}

func (c *LoginController) Get() {
	isExit := c.Input().Get("exit") == "true"
	failed := c.Input().Get("failed") != ""
	if isExit || failed {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("title", "", -1, "/")
		fmt.Printf("Clear cookie\n")
	}
	if isExit {
		c.Redirect("/", 302)
	} else {
		if failed {
			c.Data["Failed"] = true
		}
		c.TplName = "login.html"
	}

}

func (c *LoginController) Post() {
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")

	fmt.Printf("Login %s:%s\n", uname, pwd)
	autoLogin := c.Input().Get("autoLogin") == "on"

	usr, err := models.GetValidAcct(uname, pwd)
	if err == nil {
		maxAge := 10
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		maxAge = 1<<31 - 1
		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")
		c.Ctx.SetCookie("title", usr.Title, maxAge, "/")
		beego.Warning("Login success:", uname, ": ", pwd)
		c.Redirect("/", 301)
	} else {
		c.Redirect("/login?failed=true", 301)
	}
	return
}

func GetCurAcct(ctx *context.Context) *models.Account {
	for _, cook := range ctx.Request.Cookies() {
		fmt.Printf("ctx cook: %+v\n", *cook)
	}
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		beego.Warning("No uname")
		return models.Guest()
	}
	uname := ck.Value
	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		beego.Warning("No pwd")
		return models.Guest()
	}
	pwd := ck.Value
	usr, err := models.GetValidAcct(uname, pwd)
	if err == nil {
		ctx.SetCookie("title", usr.Title) //Update title in case changed
		return usr
	}
	beego.Warning("Invalid pwd ", uname, ": ", pwd)
	return models.Guest()
}

/*
func IsGuest(ctx *context.Context) bool {
	usr := GetCurUser(ctx)
	fmt.Printf("User info:%+v\n", usr)
	return usr.Uname == "Guest"
}
func IsAdmin(ctx *context.Context) bool {
	usr := GetCurUser(ctx)
	return usr.Title == "Admin"
}
func GetTitle(ctx *context.Context) string {
	usr := GetCurUser(ctx)
	fmt.Printf("User info:%+v\n", usr)
	return usr.Title
}
*/
