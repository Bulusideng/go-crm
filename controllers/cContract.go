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
	beego.Controller
}

func (this *ContractController) Get() {
	this.Data["MgrClient"] = true

	contracts, err := models.GetAllContracts()
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Contracts"] = contracts
		this.Data["Selectors"] = GetSelectors(contracts, map[string]string{})
	}

	this.TplName = "contracts.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
}

func (this *ContractController) Query() { //Filter contract
	this.Data["MgrClient"] = true

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
			if val != "" && val != "ALL" {
				filters[typeOfType.Field(i).Name] = val
				fmt.Printf("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
			}
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
		//this.Data["Cnames"] = models.GetCnames()
	}
	this.TplName = "contracts.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
}

func GetSelectors(contracts []*models.Contract, filters map[string]string) *models.ContractSelector {
	selectors := models.NewContractSelectors()
	for _, c := range contracts {
		selectors.Contract_id.List[c.Contract_id] = true
		selectors.Client_name.List[c.Client_name] = true
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

func (this *ContractController) Post() {
	usr := GetCurAcct(this.Ctx)
	if usr.IsGuest() {
		this.Redirect("/login", 302)
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
		c.Create_by = usr.Cname
		c.Create_date = time.Now().Format(time.RFC1123)
		err = models.AddContract(c)
	} else if op == "UPDATE" {
		var changes *models.ChangeSlice
		changes, err = models.UpdateContract(c)
		if err == nil {
			txt := this.GetString("NewComment", "")
			if len(txt) > 0 || len(*changes) > 0 {
				cmt := &models.Comment{
					Contract_id: c.Contract_id,
					Title:       usr.Title,
					Uname:       usr.Uname,
					Cname:       usr.Cname,
					Date:        time.Now().Format(time.RFC1123),
					Changes:     changes.String(),
					Content:     txt,
				}
				err = models.AddComment(cmt)
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
	if GetCurAcct(this.Ctx).IsGuest() {
		this.Redirect("/login", 302)
		return
	}

	this.TplName = "contract_add.html"
	this.Data["AddClient"] = true
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["Team"], _ = models.GetNonAdmins() //Consulter and Secretary are limited to this set
}

func (this *ContractController) Update() {
	if GetCurAcct(this.Ctx).IsGuest() {
		this.Redirect("/login", 302)
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
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["Contract"] = contract
	this.Data["Comments"], _ = models.GetComments(cid)
	this.Data["Team"], _ = models.GetNonAdmins() //Consulter and Secretary are limited to this set
}

func (this *ContractController) View() {
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
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
}
