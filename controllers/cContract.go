package controllers

import (
	"errors"
	"path"
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

	/*
		if !curUser.IsWorker() {
			if curUser.IsAdmin() {
				this.Redirect("/account", 302)
				return false
			} else {
				this.Redirect("/login", 302)
				return false
			}
		}
	*/
	if curUser.IsGuest() {
		this.Redirect("/login", 302)
		return false
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
	this.Data["RICH"] = IsRichView()

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
		cfilter := models.NewAllFilter()
		this.ParseForm(cfilter)
		object := reflect.ValueOf(cfilter)
		myref := object.Elem()
		typeOfType := myref.Type()
		for i := 0; i < myref.NumField(); i++ {
			field := myref.Field(i)
			val := field.Interface().(string)
			if val != "ALL" { //Get valid filters
				filters[typeOfType.Field(i).Name] = val
				beego.BeeLogger.Debug("%d. %s %s = %v \n", i, typeOfType.Field(i).Name, field.Type(), field.Interface())
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

	op := this.Input().Get("op")
	beego.Debug("op: [", op, "] on: ", *c)

	for k, v := range this.Ctx.Request.Form {
		beego.Debug("Param ", k, ":", v)
	}

	contractURL := "/contract/view?cid=" + string(c.Contract_id)
	if op == "add" {
		//c.Create_by = usr.Title + " " + usr.Cname + "[" + usr.Uname + "]"
		c.Create_by = curUser.Uname
		c.Create_date = time.Now().Format(time.RFC822)
		c.Secretaries = strings.Join(this.Ctx.Request.Form["Secretaries"], "&")
		err = models.AddContract(c)
		if err == nil {
			this.RedirectTo("/status", "添加成功!", contractURL, 302)
		}
	} else if op == "update" {
		oldContractId := this.GetString("oldContractId", "")
		if oldContractId == "" {
			oldContractId = c.Contract_id
		}
		c.Consulters = strings.Join(this.Ctx.Request.Form["Consulters"], "&")
		c.Secretaries = strings.Join(this.Ctx.Request.Form["Secretaries"], "&")
		//fmt.Printf("New values: %+v\n", *c)
		var changes *models.ChangeSlice
		changes, err = models.UpdateContract(oldContractId, c)
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
				if err == nil {
					this.RedirectTo("/status", "更新成功!", contractURL, 302)
				} else {
					this.RedirectTo("/status", "更新失败: "+err.Error(), contractURL, 302)
				}
			} else {
				this.RedirectTo("/status", "没有改动!", contractURL, 302)
			}
		}

	} else if op == "backup" {
		pwd := this.GetString("pwd", "")
		if !curUser.ValidPwd(pwd) {
			this.RedirectTo("/status", "密码错误!", "/contract/backup", 302)
			return
		}
		if !(curUser.IsManager() || curUser.IsAdmin()) {
			this.RedirectTo("/status", "当前帐号没有备份权限!", "/contract", 302)
			return
		}

		cat := this.GetString("cat", "")
		fn := ""
		withComments := (this.GetString("withComments", "") == "on")
		fn, err = models.ExportContracts(cat, withComments)

		if err != nil {
			this.RedirectTo("/status", "备份失败:"+err.Error(), "/contract/backup", 302)
		} else {
			this.Redirect("/contract/backup?file="+fn, 302)
		}

		return
	} else {
		beego.Error("Invalid op:%s", op)
		err = errors.New("非法操作: " + op)
		this.RedirectTo("/status", "操作失败: "+err.Error(), "/contract", 302)
	}

	if err != nil {
		beego.Error(err)
		this.RedirectTo("/status", "操作失败: "+err.Error(), "/contract/backup", 302)
	}

}

func (this *ContractController) Backup() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}
	if !(curUser.IsManager() || curUser.IsAdmin()) {
		this.RedirectTo("/status", "没有权限！", "/account", 302)
		return
	}

	this.TplName = "contract_backup.html"
	this.Data["BakClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["RICH"] = IsRichView()
	fn := this.GetString("file", "")
	if fn != "" {
		this.Data["FILE"] = path.Join("/files", fn)
	}
}

func (this *ContractController) Add() {
	curUser := GetCurAcct(this.Ctx)
	if !this.valid(curUser) {
		return
	}
	if !(curUser.IsManager() || curUser.IsAdmin()) {
		this.RedirectTo("/status", "没有权限！", "/account", 302)
		return
	}

	this.TplName = "contract_add.html"
	this.Data["AddClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["RICH"] = IsRichView()
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

func (this *ContractController) View() {
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
	Perm := "Read"
	//Manager, full write permission, consulter or secretary, partial wirte permission
	if curUser.IsManager() || curUser.IsAdmin() {
		Perm = "Write"
	} else {
		if strings.Contains(contract.Consulters, curUser.Cname) || strings.Contains(contract.Secretaries, curUser.Cname) {
			Perm = "ParWrite" //Partial write
		}
	}

	this.TplName = "contract.html"
	this.Data["MgrClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["Contract"] = contract
	this.Data["Perm"] = Perm
	this.Data["Comments"], _ = models.GetComments(cid)
	this.Data["Team"], _ = models.GetNonAdmins() //Consulter and Secretary are limited to this set
}

func (this *ContractController) Viewo() {
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
	this.TplName = "contract.html"
	this.Data["MgrClient"] = true
	this.Data["CurUser"] = curUser
	this.Data["Contract"] = contract
	this.Data["Comments"], _ = models.GetComments(cid)

}

func GetSelectors(contracts []*models.Contract, filters map[string]string) *models.ContractSelector {
	selectors := models.NewContractSelectors()
	for _, c := range contracts {
		selectors.Contract_id.List[c.Contract_id] = true
		selectors.Client_name.List[c.Client_name] = true
		selectors.Client_tel.List[c.Client_tel] = true
		selectors.Consulters.List[c.Consulters] = true
		selectors.Secretaries.List[c.Secretaries] = true
		selectors.Country.List[c.Country] = true
		selectors.Project_type.List[c.Project_type] = true
		selectors.Contract_date.List[c.Contract_date] = true
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
	if key, ok := filters["Consulters"]; ok {
		selectors.Consulters.CurSelected = key
	}
	if key, ok := filters["Secretaries"]; ok {
		selectors.Secretaries.CurSelected = key
	}
	if key, ok := filters["Country"]; ok {
		selectors.Country.CurSelected = key
	}
	if key, ok := filters["Project_type"]; ok {
		selectors.Project_type.CurSelected = key
	}
	if key, ok := filters["Contract_date"]; ok {
		selectors.Contract_date.CurSelected = key
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
