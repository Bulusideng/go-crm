package controllers

import (
	"fmt"
	"reflect"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) Get() {
	if IsGuest(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}
	this.Data["IsUser"] = true
	this.Data["CurUser"] = GetUserInfo(this.Ctx)

	users, err := models.GetAllUsers()
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Users"] = users
	}
	this.TplName = "users.html"
}

func (this *UserController) Query() { //Filter contract
	this.Data["IsContract"] = true
	this.Data["CurUser"] = GetUserInfo(this.Ctx)
	filters := map[string]string{}

	cfilter := &models.Contract{}
	this.ParseForm(cfilter)
	object := reflect.ValueOf(cfilter)
	myref := object.Elem()
	typeOfType := myref.Type()
	for i := 0; i < myref.NumField(); i++ {
		field := myref.Field(i)
		val := field.Interface().(string)
		if val != "" && val != "ALL" {
			filters[typeOfType.Field(i).Name] = val
			fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
		}
	}

	contracts, err := models.GetContracts(filters)
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Contracts"] = contracts
		this.Data["Filters"] = filters
		contracts, _ := models.GetAllContracts()
		this.Data["Selectors"] = GetSelectors(contracts, filters)
	}
	this.TplName = "contracts.html"
}

func (this *UserController) Post() {
	if IsGuest(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	usr := &models.User{}
	err := this.ParseForm(usr)
	if err != nil {
		beego.Error(err)
		return
	}

	curPwd := this.GetString("CurPwd", "")
	this.Ctx.SetCookie("pwd", curPwd)

	if IsGuest(this.Ctx) { //Incorrect pwd, force to login again
		this.Redirect("/login", 302)
		return
	}

	op := this.Input().Get("OP")
	fmt.Printf("%s User:%+v\n", op, *usr)

	if op == "ADD" {
		err = models.AddUser(usr)
	} else if op == "UPDATE" {
		err = models.UpdateUser(usr)
	} else if op == "ADMIN" {
		if IsAdmin(this.Ctx) {
			err = models.UpdateUser(usr)
		}
	} else {
		beego.Error("Invalid op:", op)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/user", 302)
}

func (this *UserController) Register() {
	this.TplName = "user_register.html"
	this.Data["CurUser"] = GetUserInfo(this.Ctx)
	this.Data["OP"] = "ADD"
	this.Data["IsAddUser"] = true

}

func (this *UserController) Update() {
	uname := this.GetString("uname", "")
	if uname == "" {
		return
	}
	curUsr := GetUserInfo(this.Ctx)
	if curUsr.Title == "Guest" {
		this.Redirect("/login", 302)
		return
	} else if curUsr.Title != "Admin" && curUsr.Uname != uname { //Only Admin self can update user
		fmt.Printf("User %+v want update %s\n", *curUsr, uname)
		return
	}
	usr, err := models.GetUser(uname)
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

	this.Data["User"] = &usr.UserInfo
}
