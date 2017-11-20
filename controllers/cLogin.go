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

func (this *LoginController) Get() {
	isExit := this.Input().Get("exit") == "true"
	failed := this.Input().Get("failed") != ""
	this.Data["chances"] = this.Input().Get("chances")

	if isExit || failed {
		this.Ctx.SetCookie("uname", "", -1, "/")
		this.Ctx.SetCookie("title", "", -1, "/")
		beego.Debug("Clear cookie...")
	}
	if isExit {
		this.Redirect("/", 302)
	} else {
		if failed {
			this.Data["Failed"] = true
		}
		this.TplName = "login.html"
		this.Data["PwdRstSpt"] = models.GetConfig().Pwdreset
	}
}

func (c *LoginController) Post() {
	uname := c.Input().Get("uname")
	pwd := c.Input().Get("pwd")

	//beego.Debug("Login ", uname, ":", pwd)

	autoLogin := c.Input().Get("autoLogin") == "on"

	usr, err := models.GetValidAcct(uname, pwd)
	if err == nil {
		maxAge := 1
		if autoLogin {
			maxAge = 1<<31 - 1
		}
		maxAge = 1<<31 - 1

		c.Ctx.SetCookie("uname", uname, maxAge, "/")
		c.Ctx.SetCookie("pwd", usr.Pwd, maxAge, "/")
		c.Ctx.SetCookie("title", usr.Title, maxAge, "/")
		beego.Debug("Login success, uname:", uname)
		models.UpdateErrCnt(uname, -1000) //Clear error cnt

		c.Ctx.ResponseWriter.ResponseWriter.Header().Set("Accept-Charset", "GBK,utf-8;q=0.7,*;q=0.3")
		c.Ctx.ResponseWriter.ResponseWriter.Header().Set("Accept-Language", "zh-CN,zh;q=0.8")

		if usr.IsAdmin() && false {
			c.Redirect("/account", 301)
			return
		} else {
			c.Redirect("/contract", 301)
			return
		}

	} else {
		c.Ctx.SetCookie("uname", "", -1, "/")
		c.Ctx.SetCookie("title", "", -1, "/")
		fmt.Printf("Clear cookie\n")
		if usr != nil {
			remainChance := 0
			remainChance, err = models.UpdateErrCnt(uname, 1)
			if remainChance > 0 {
				c.RedirectTo("/status", "登录失败，还有"+strconv.Itoa(remainChance)+"次机会!", "/login", 301)
				return
			} else {
				c.RedirectTo("/status", "登录失败，账户锁定，请联系你的经理: "+usr.Manager, "", 302)
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
	epwd := ck.Value
	usr, err := models.GetValidAcctEnc(uname, epwd)
	if err == nil {
		ctx.SetCookie("title", usr.Title) //Update title in case changed
		return usr
	}
	beego.Warning("Invalid pwd ", uname, ": ", epwd)
	return models.Guest()
}
