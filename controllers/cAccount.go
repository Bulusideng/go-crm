package controllers

import (
	"fmt"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type AccountController struct {
	beego.Controller
}

func (this *AccountController) Get() {
	curUser := GetCurAcct(this.Ctx)
	if curUser.IsGuest() {
		this.Redirect("/login", 302)
		return
	}
	if !curUser.IsManager() {
		this.Redirect("/status?msg=No access!!!", 302)
		return
	}

	this.Data["IsUser"] = true
	this.Data["CurUser"] = curUser

	reporters, err := curUser.GetReporters()
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Reporters"] = reporters
	}
	this.TplName = "accounts.html"
}

func (this *AccountController) View() {
	this.TplName = "account_view.html"
	this.Data["ViewAcct"] = true
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
}

func (this *AccountController) Register() {
	this.TplName = "account_register.html"
	this.Data["RegAcct"] = true
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	filter := map[string]string{"title": "Manager"}
	mgrs, _ = models.GetAcctounts(filter)
}

func (this *AccountController) Post() {
	acct := &models.Account{}
	err := this.ParseForm(acct)
	if err != nil {
		beego.Error(err)
		return
	}

	op := this.Input().Get("op") //Operations: register, change_contact, change_pwd, change_status
	fmt.Printf("%s User: %+v\n", op, *acct)

	curUsr := GetCurAcct(this.Ctx)
	err = nil
	switch op {
	case "register": //No accout validation required
		RePwd := this.GetString("RePwd")
		if RePwd == acct.Pwd {
			if err = acct.Register(); err == nil {
				this.Redirect("/status?msg=Register success, please wait manager approve", 302)
				return
			}
		} else {
			this.Redirect("/status?msg=Register failed, password mismatch", 302)
			return
		}
	case "change_contact":
		if curUsr.IsActive() {
			if err = models.UpdateContact(acct.Uname, acct.Email, acct.Mobile); err == nil {
				this.Redirect("/status?msg=Update contact success", 302)
				return
			}
		}
	case "change_pwd":
		curPwd := this.GetString("CurPwd", "")
		if curUsr.IsActive() {
			if curUsr.Uname == acct.Uname {
				if curPwd == curUsr.Pwd {
					if err = models.UpdatePwd(acct.Uname, acct.Pwd); err == nil {
						this.Redirect("/status?msg=Update passwrod success", 302)
						return
					}
				} else {
					this.Redirect("/status?msg=Update password failed, invalid password", 302)
					return
				}
			} else {
				this.Redirect("/status?msg=Update password failed, account mismatch", 302)
				return
			}
		}
	case "manage_acct":
		if curUsr.IsActive() {
			if curUsr.IsManagerOf(acct) {
				models.ManageAccount(acct.Uname, acct.Status, acct.Title, acct.Manager)
			} else {
				beego.Error("没有权限！！！")
				this.Redirect("/status?msg=Update account failed, permission denied", 302)
				return
			}
		}
	default:
		beego.Error("无效操作！！！")
		this.Redirect("/status?msg=Invalid operation", 302)
		return
	}
	if err != nil {
		beego.Error(err)
		this.Redirect("/status?msg="+err.Error(), 302)
	} else if curUsr.IsGuest() {
		this.Redirect("/login", 302)
		return
	} else if curUsr.Locked() {
		this.Redirect("/user/findpwd", 302)
		return
	} else if curUsr.Disabled() {
		beego.Error("账户禁用！！！")
		this.Redirect("/status?msg=Account disabled", 302)
		return
	}
	this.Redirect("/status?msg=unknown error", 302)
}

/*
func (this *AccountController) Update() {
	uname := this.GetString("uname", "")
	if uname == "" {
		return
	}
	curUsr := GetCurAcct(this.Ctx)
	if curUsr.IsGuest() {
		this.Redirect("/login", 302)
		return
	} else if curUsr.Title != "Admin" && curUsr.Uname != uname { //Only Admin self can update user
		fmt.Printf("User %+v want update %s\n", *curUsr, uname)
		return
	}
	usr, err := models.GetAccount(uname)
	if err != nil {
		beego.Error(err)
		this.Redirect("/", 302)
		return
	}

	this.TplName = "user_update.html"
	this.Data["CurUser"] = curUsr
	if curUsr.Title != "Admin" {
		this.Data["OP"] = "UPDATE"
	} else {
		this.Data["OP"] = "ADMIN" //Admin change other user info
	}

	this.Data["User"] = usr
}
*/
