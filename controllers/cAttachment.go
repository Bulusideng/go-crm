package controllers

import (
	"github.com/Bulusideng/go-crm/models"
)

type AttachmentController struct {
	FileController
}

func (this *AttachmentController) Delete() {
	cid := this.GetString("cid", "")
	id, err := this.GetInt("id", -1)
	if err == nil {
		err = models.DelAttachment(id)
	}
	if err != nil {
		this.RedirectTo("/status", "删除附件失败:"+err.Error(), "/attachment/view?cid="+cid, 302)
		return
	}
	this.RedirectTo("/status", "删除附件成功:", "/attachment/view?cid="+cid, 302)
}
func (this *AttachmentController) View() {
	cid := this.GetString("cid", "-1")
	this.Data["Attachments"], _ = models.GetAttachments(cid)
	this.TplName = "attachments.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["cid"] = cid
}
