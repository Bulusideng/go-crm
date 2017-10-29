package controllers

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/Bulusideng/go-crm/models"
	"github.com/astaxie/beego"
)

type ContractController struct {
	beego.Controller
}

func (this *ContractController) Get() {
	this.Data["IsContract"] = true

	contracts, err := models.GetAllContracts()
	if err != nil {
		beego.Error(err)
	} else {
		this.Data["Contracts"] = contracts
		this.Data["Selectors"] = GetSelectors(contracts, map[string]string{})
	}

	this.TplName = "contracts.html"
	this.Data["CurUser"] = GetUserInfo(this.Ctx)
}

func (this *ContractController) Query() { //Filter contract
	this.Data["IsContract"] = true

	filters := map[string]string{}

	filter := this.GetString("filter")
	if filter != "" {
		beego.Warning("Filter: ", this.GetString("filter"))
		kv := strings.Split(filter, ":")
		if len(kv) == 2 {
			filters[kv[0]] = kv[1]
		}
	} else {
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
	}
	this.TplName = "contracts.html"
	this.Data["CurUser"] = GetUserInfo(this.Ctx)
}

func GetSelectors(contracts []*models.Contract, filters map[string]string) *models.ContractSelector {
	vals := models.NewContractSelector()
	for _, c := range contracts {
		vals.Contract_id.Values[c.Contract_id] = true
		vals.Client_name.Values[c.Client_name] = true
		vals.Consulter.Values[c.Consulter] = true
		vals.Secretary.Values[c.Secretary] = true
		vals.Country.Values[c.Country] = true
		vals.Project_type.Values[c.Project_type] = true
	}
	if key, ok := filters["Contract_id"]; ok {
		vals.Contract_id.Filter = key
	}
	if key, ok := filters["Client_name"]; ok {
		vals.Client_name.Filter = key
	}
	if key, ok := filters["Consulter"]; ok {
		vals.Consulter.Filter = key
	}
	if key, ok := filters["Secretary"]; ok {
		vals.Secretary.Filter = key
	}
	if key, ok := filters["Country"]; ok {
		vals.Country.Filter = key
	}
	if key, ok := filters["Project_type"]; ok {
		vals.Project_type.Filter = key
	}
	return vals
}

func (this *ContractController) Post() {
	if IsGuest(this.Ctx) {
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
	if op == "ADD" {
		err = models.AddContract(c)
	} else if op == "UPDATE" {
		err = models.UpdateContract(c)
	} else {
		beego.Error("Invalid op:%s", op)
	}

	if err != nil {
		beego.Error(err)
	}
	this.Redirect("/contract", 302)
}

func (this *ContractController) Add() {
	if IsGuest(this.Ctx) {
		this.Redirect("/login", 302)
		return
	}

	this.TplName = "contract_add.html"
	this.Data["IsAddContract"] = true
	this.Data["CurUser"] = GetUserInfo(this.Ctx)

}

func (this *ContractController) Update() {
	if IsGuest(this.Ctx) {
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
	this.Data["Contract"] = contract
	this.Data["Cid"] = cid

}

func (this *ContractController) View() {
	this.TplName = "contract_view.html"
	cid := this.Ctx.Input.Params()["0"]
	contract, err := models.GetContract(cid)
	if err != nil {
		beego.Error(err)
		this.Redirect("/contract", 302)
		return
	}
	this.Data["Contract"] = contract
	this.Data["Tid"] = this.Ctx.Input.Params()["0"]
	this.Data["Account"] = GetUserInfo(this.Ctx)
}
