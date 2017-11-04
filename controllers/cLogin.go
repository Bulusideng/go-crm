package controllers

import (
	"fmt"
	"strconv"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

type LoginController struct {
	baseController
}

func (c *LoginController) Get() {
	isExit := c.Input().Get("exit") == "true"
	failed := c.Input().Get("failed") != ""
	c.Data["chances"] = c.Input().Get("chances")

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
		maxAge := 1
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		maxAge = 1<<31 - 1

		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", pwd, maxAge, "/")
		c.Ctx.SetCookie("title", usr.Title, maxAge, "/")
		beego.Warning("Login success:", uname, ": ", pwd)
		models.UpdateErrCnt(uname, -1000) //Clear error cnt

		if usr.IsAdmin() {
			c.Redirect("/account", 301)
			return
		} else {
			c.Redirect("/contract", 301)
			return
		}

	} else {
		if usr != nil {
			remainChance := 0
			remainChance, err = models.UpdateErrCnt(uname, 1)
			if remainChance > 0 {
				c.Redirect("/login?chances="+strconv.Itoa(remainChance), 301)
				return
			} else {
				c.RedirectTo("/status", "Login failed, account locked!", "/login", 302)
				return
			}
		}
	}

	c.RedirectTo("/status", "Login failed: "+err.Error(), "/login", 302)
	return
}

//Get account in cookie, check password and status is active, if failed return guest
func GetCurAcct(ctx *context.Context) *models.Account {
	ck, err := ctx.Request.Cookie("uname")
	if err != nil {
		beego.Warning("Get uname cookie error: ", err.Error())
		return models.Guest()
	}
	uname := ck.Value

	ck, err = ctx.Request.Cookie("pwd")
	if err != nil {
		beego.Warning("Get pwd cookie error: ", err.Error())
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
