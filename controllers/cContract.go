package controllers

import (
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type ContractController struct {
	baseController
}

func (this *ContractController) valid(curUser *models.Account) bool {
	if !curUser.IsWorker() {
		if curUser.IsAdmin() {
			this.Redirect("/account", 302)
			return false
		} else {
			this.Redirect("/login", 302)
			return false
		}
	}
	return true
}

func (this *ContractController) Get() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}
	contracts, err := models.GetAllContracts()
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Contracts"] = contracts
		this.Data["Selectors"] = GetSelectors(contracts, map[string]string{})
	}
	this.TplName = "contracts.html"
	this.Data["MgrClient"] = true
	this.Data["CurUser"] = curUser
}

func (this *ContractController) Query() { //Filter contract
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}
	filters := map[string]string{} //fieldname: value

	filter := this.GetString("filter")
	if filter != "" { //Single filter
		beego.Warning("Filter: ", this.GetString("filter"))
		kv := strings.Split(filter, ":")
		if len(kv) == 2 {
			filters[kv[0]] = kv[1]
		}
	} else { //Filter in form
		cfilter := &models.Contract{}
		this.ParseForm(cfilter)
		object := reflect.ValueOf(cfilter)
		myref := object.Elem()
		typeOfType := myref.Type()
		for i := 0; i < myref.NumField(); i++ {
			field := myref.Field(i)
			val := field.Interface().(string)
			if val != "" && val != "ALL" { //Get valid filters
				filters[typeOfType.Field(i).Name] = val
				fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
			}
		}
	}

	contracts, err := models.GetContracts(filters)
	if err != nil {
		beego.Error(err)
	} else {
		this.TplName = "contracts.html"
		this.Data["MgrClient"] = true
		this.Data["CurUser"] = curUser
		this.Data["Contracts"] = contracts
		this.Data["Filters"] = filters
		contracts, _ := models.GetAllContracts()
		this.Data["Selectors"] = GetSelectors(contracts, filters)
		//this.Data["Cnames"] = models.GetCnames()
	}
}

func (this *ContractController) Post() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}

	c := models.NewContract()
	err := this.ParseForm(c)
	if err != nil {
		beego.Error(err)
		return
	}

	op := this.Input().Get("OP")
	beego.Warning("op:%s: %+v", op, *c)
	cons, _ := models.GetAccount(c.Consulter)
	c.Consulter_name = cons.Cname
	if op == "ADD" {
		//c.Create_by = usr.Title + " " + usr.Cname + "[" + usr.Uname + "]"
		c.Create_by = curUser.Uname
		c.Create_date = time.Now().Format(time.RFC822)
		err = models.AddContract(c)
	} else if op == "UPDATE" {
		var changes *models.ChangeSlice
		changes, err = models.UpdateContract(c)
		if err == nil {
			txt := this.GetString("NewComment", "")
			if len(txt) > 0 || len(*changes) > 0 {
				cmt := &models.Comment{
					Contract_id: c.Contract_id,
					Title:       curUser.Title,
					Uname:       curUser.Uname,
					Cname:       curUser.Cname,
					Date:        time.Now().Format(time.RFC1123),
					Changes:     changes.String(),
					Content:     txt,
				}
				err = models.AddComment(cmt)
			}
			if err == nil {
				fmt.Printf("URL:%s, URI:%s, rq:%s, rp:%s\n", this.Ctx.Request.URL.String(), this.Ctx.Request.URL.RequestURI(), this.Ctx.Request.URL.RawQuery, this.Ctx.Request.URL.RawPath)
				this.RedirectTo("/status", "Update success!", this.Ctx.Request.URL.String()+"/update?cid="+c.Contract_id, 302)
			}
		}

	} else {
		beego.Error("Invalid op:%s", op)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/contract", 302)
}

func (this *ContractController) Add() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}

	this.TplName = "contract_add.html"
	this.Data["AddClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["Team"], _ = models.GetNonAdmins() //Consulter and Secretary are limited to this set
}

func (this *ContractController) Delete() {

	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}

	if !curUser.IsManager() {
		this.RedirectTo("/status", "Only manager can delete!!!", "/contract", 302)
		return
	}
	cid := this.GetString("cid", "")
	if len(cid) > 0 {
		err := models.DelContract(cid)
		if err == nil {
			this.RedirectTo("/status", "删除成功!!!", "/contract", 302)
			return
		} else {
			this.RedirectTo("/status", "Delete failed: "+err.Error(), "/contract", 302)
		}
	} else {
		this.RedirectTo("/status", "No contract id provide!!", "/contract", 302)
	}
}

func (this *ContractController) Update() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}

	cid := this.GetString("cid", "-1")
	contract, err := models.GetContract(cid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/contract", 302)
		return
	}
	this.TplName = "contract_update.html"
	this.Data["MgrClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["Contract"] = contract
	this.Data["Comments"], _ = models.GetComments(cid)
	this.Data["Team"], _ = models.GetNonAdmins() //Consulter and Secretary are limited to this set
}

func (this *ContractController) View() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}
	//cid := this.Ctx.Input.Params()["0"]
	cid := this.GetString("cid", "-1")
	contract, err := models.GetContract(cid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/contract", 302)
		return
	}
	this.TplName = "contract_view.html"
	this.Data["MgrClient"] = true
	this.Data["Contract"] = contract
	this.Data["Comments"], _ = models.GetComments(cid)
	this.Data["CurUser"] = curUser
}

func GetSelectors(contracts []*models.Contract, filters map[string]string) *models.ContractSelector {
	selectors := models.NewContractSelectors()
	for _, c := range contracts {
		selectors.Contract_id.List[c.Contract_id] = true
		selectors.Client_name.List[c.Client_name] = true
		selectors.Client_tel.List[c.Client_tel] = true
		selectors.Consulter.List[c.Consulter] = true
		selectors.Secretary.List[c.Secretary] = true
		selectors.Country.List[c.Country] = true
		selectors.Project_type.List[c.Project_type] = true
		selectors.Zhuan_an_date.List[c.Zhuan_an_date] = true
		selectors.Create_date.List[c.Create_date] = true
		selectors.Create_by.List[c.Create_by] = true
		selectors.Current_state.List[c.Current_state] = true

	}
	if key, ok := filters["Contract_id"]; ok {
		selectors.Contract_id.CurSelected = key
	}
	if key, ok := filters["Client_name"]; ok {
		selectors.Client_name.CurSelected = key
	}
	if key, ok := filters["Client_tel"]; ok {
		selectors.Client_tel.CurSelected = key
	}
	if key, ok := filters["Consulter"]; ok {
		selectors.Consulter.CurSelected = key
	}
	if key, ok := filters["Secretary"]; ok {
		selectors.Secretary.CurSelected = key
	}
	if key, ok := filters["Country"]; ok {
		selectors.Country.CurSelected = key
	}
	if key, ok := filters["Project_type"]; ok {
		selectors.Project_type.CurSelected = key
	}
	if key, ok := filters["Zhuan_an_date"]; ok {
		selectors.Zhuan_an_date.CurSelected = key
	}
	if key, ok := filters["Create_date"]; ok {
		selectors.Create_date.CurSelected = key
	}
	if key, ok := filters["Create_by"]; ok {
		selectors.Create_by.CurSelected = key
	}
	if key, ok := filters["Current_state"]; ok {
		selectors.Current_state.CurSelected = key
	}
	return selectors
}
