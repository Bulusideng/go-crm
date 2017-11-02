package controllers

type StatusController struct {
	baseController
}

func (this *StatusController) Get() {

	this.TplName = "status.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["Message"] = this.GetString("Message")
	this.Data["RedirectTo"] = this.GetString("RedirectTo")
}
