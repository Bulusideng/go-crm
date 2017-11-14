package controllers

import (
	"fmt"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type AccountController struct {
	baseController
}

func (this *AccountController) Get() {
	curUser := GetCurAcct(this.Ctx)
	if curUser.IsGuest() {
		this.Redirect("/login", 302)
		return
	}
	if !(curUser.IsManager() || curUser.IsAdmin()) {
		this.RedirectTo("/status", "No access!!!", "/account", 302)
		return
	}

	this.Data["MgrAcct"] = true
	this.Data["CurUser"] = curUser
	var reporters []*models.Account
	var err error
	if curUser.IsAdmin() || true {
		reporters, err = models.GetAllAccts()
	} else if curUser.IsManager() {
		reporters, err = curUser.GetReporters()
	}

	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Reporters"] = reporters
	}
	this.TplName = "accounts.html"

}

func (this *AccountController) Register() {
	this.TplName = "account_register.html"
	this.Data["RegAcct"] = true
	this.Data["RICH"] = IsRichView()
	mgrs, _ := models.AcctByTitle("Manager")
	admins, _ := models.AcctByTitle("Admin")
	this.Data["Managers"] = append(mgrs, admins...)
}

func (this *AccountController) PwdForget() {
	this.TplName = "account_pwdforget.html"
	this.Data["RICH"] = IsRichView()
}

func (this *AccountController) PwdReset() {

	uname := this.GetString("Uname", "")
	randomstr := this.GetString("Random", "")
	if uname == "" || randomstr == "" {
		this.RedirectTo("/status", "非法链接!", "/login", 302)
		return
	}
	acct, err := models.GetAccount(uname)
	if err != nil {
		this.RedirectTo("/status", "用户不存在!", "/login", 302)
		return
	}
	if randomstr != acct.Random {
		this.RedirectTo("/status", "重置链接已失效!", "/login", 302)
		return
	}
	this.TplName = "account_pwdreset.html"
	this.Data["Uname"] = uname
	this.Data["RICH"] = IsRichView()
}

func (this *AccountController) View() {
	curUsr := GetCurAcct(this.Ctx)
	if curUsr.IsGuest() {
		this.Redirect("/login", 302)
		return
	}
	this.TplName = "account.html"
	this.Data["ViewAcct"] = true
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
}

func (this *AccountController) Manage() {
	uname := this.GetString("uname", "")
	if uname == "" {
		return
	}
	acct, err := models.GetAccount(uname)
	if err != nil {
		this.RedirectTo("/status", "Manage account error, invalid Uname: "+uname, "/account", 302)
		return
	}
	curUsr := GetCurAcct(this.Ctx)
	if curUsr.IsGuest() {
		this.Redirect("/login", 302)
		return
	} else if !curUsr.IsManagerOf(acct) {
		this.RedirectTo("/status", "Current user: "+curUsr.Uname+"is not manager of: "+uname, "/account", 302)
		beego.Warn("User ", *curUsr, " want to update ", uname)
		return
	}

	this.TplName = "account_manage.html"
	this.Data["IsUser"] = true
	this.Data["CurUser"] = curUsr
	this.Data["Acct"] = acct

	mgrs, _ := models.AcctByTitle("Manager")
	admins, _ := models.AcctByTitle("Admin")
	this.Data["Managers"] = append(mgrs, admins...)
}

func (this *AccountController) Post() {
	acct := &models.Account{}
	err := this.ParseForm(acct)
	if err != nil {
		beego.Error(err)
		return
	}
	curUsr := GetCurAcct(this.Ctx)
	op := this.Input().Get("op") //Operations: register, change_contact, change_pwd, change_status
	beego.Debug("User: ", curUsr.Uname, "[", op, "] account:", *acct)

	err = nil
	switch op {
	case "register": //No accout validation required
		RePwd := this.GetString("RePwd")
		//fmt.Printf("p:[%s], rep:[%s]\n", acct.Pwd, RePwd)
		if RePwd == acct.Pwd {
			if err = acct.Register(); err == nil {
				this.RedirectTo("/status", "Register success, please wait manager approve", "/account/register", 302)
				return
			} else {
				this.RedirectTo("/status", "Register failed: 用户名已存在", "/account/register", 302)
				return
			}
		} else {
			this.RedirectTo("/status", "Register failed, password mismatch", "/account/register", 302)
			return
		}
	case "forgetpwd":
		err := acct.ForgetPwd()
		if err != nil {
			this.RedirectTo("/status", "找回密码失败:"+err.Error(), "/login", 302)
			return
		}
		link := "/account/pwdreset?Uname=" + acct.Uname + "&Random=" + acct.Random
		//TODO send link through email
		if true {
			beego.Debug("Host: ", this.Ctx.Request.Host)
			this.RedirectTo("/status", "重置密码链接已发送至邮箱！", link, 302)
		} else {
			if SendEmail(acct.Email, link) {
				this.RedirectTo("/status", "重置密码链接已发送至邮箱！", "", 302)
			} else {
				this.RedirectTo("/status", "重置邮件发送失败！", "", 302)
			}
		}

		return

	case "resetpwd":
		uname := this.GetString("Uname", "")
		pwd := this.GetString("Pwd", "")
		rePwd := this.GetString("RePwd", "")
		if uname == "" {
			this.RedirectTo("/status", "非法用户名", "", 302)
			return
		}
		if pwd == "" || pwd != rePwd {
			this.RedirectTo("/status", "密码不一致", "", 302)
			return
		}
		if err = models.UpdatePwd(uname, pwd, true); err == nil {
			this.RedirectTo("/status", "Update password success", "/login", 302)
			return
		}
	case "change_contact":
		if curUsr.IsActive() {
			if err = models.UpdateContact(acct.Uname, acct.Email, acct.Mobile); err == nil {
				this.RedirectTo("/status", "Update contact success", "/account/view", 302)
				return
			}
		}
	case "change_pwd":
		curPwd := this.GetString("CurPwd", "")
		RePwd := this.GetString("RePwd", "")
		//fmt.Printf("p:[%s], rep:[%s]\n", acct.Pwd, RePwd)

		if !curUsr.ValidPwd(curPwd) {
			this.RedirectTo("/status", "Update password failed, Current password incorrect", "/account/view", 302)
			return
		}
		if curUsr.IsActive() {
			if curUsr.Uname == acct.Uname {
				if RePwd == acct.Pwd {
					if err = models.UpdatePwd(acct.Uname, acct.Pwd, false); err == nil {
						this.RedirectTo("/status", "Update passwrod success", "/account/view", 302)
						return
					}
				} else {
					this.RedirectTo("/status", "Update password failed, password mismatch", "/account/view", 302)
					return
				}
			} else {
				this.RedirectTo("/status", "Update password failed, account mismatch", "/account/view", 302)
				return
			}
		}
	case "manage_acct":
		CurPwd := this.GetString("CurPwd", "")
		if !curUsr.ValidPwd(CurPwd) {
			this.RedirectTo("/status", "Manage account failed, incorrect password!", "/account", 302)
			return
		}

		if curUsr.IsActive() {
			origAcct, _ := models.GetAccount(acct.Uname)
			if curUsr.IsManagerOf(origAcct) {
				if err = models.ManageAccount(acct); err == nil {
					this.RedirectTo("/status", "Manage account success", "/account", 302)
					ct, _ := models.GetAccount(acct.Uname)
					beego.Debug("After manage: ", *ct)
					return
				}
			} else {
				beego.Error("没有权限！！！")
				this.RedirectTo("/status", "Manage account failed, permission denied, "+curUsr.Uname+" is not manager of: "+acct.Uname, "/account", 302)
				return
			}
		} else {
			this.RedirectTo("/status", "Manage account failed, current account is inactive!", "/account", 302)
			return
		}
	default:
		beego.Error("无效操作！！！")
		this.RedirectTo("/status", "Invalid operation", "/account", 302)
		return
	}
	if err != nil {
		beego.Error(err)
		this.RedirectTo("/status", ""+err.Error(), "/account", 302)
	} else if curUsr.IsGuest() {
		this.Redirect("/login", 302)
		return
	} else if curUsr.Locked() {
		this.Redirect("/user/findpwd", 302)
		return
	} else if curUsr.Disabled() {
		beego.Error("账户禁用！！！")
		this.RedirectTo("/status", "Account disabled", "/account", 302)
		return
	}
	fmt.Printf("账户操作失败: %s\n", op)
	//fmt.Printf("CurUser: %+v\n", *curUsr)
	//fmt.Printf("Account: %+v\n", *acct)
	this.RedirectTo("/status", "unknown error", "/account", 302)
}
