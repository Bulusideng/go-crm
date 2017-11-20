package controllers

import (
	"time"

	"github.com/Bulusideng/go-crm/models"
)

type AttachmentController struct {
	FileController
}

func (this *AttachmentController) Delete() {
	curUser := GetCurAcct(this.Ctx)
	cid := this.GetString("cid", "")
	id, err := this.GetInt("id", -1)
	fname := ""
	if err == nil {
		fname, err = models.DelAttachment(id)
	}
	if err != nil {
		this.RedirectTo("/status", "删除附件失败:"+err.Error(), "/attachment/view?cid="+cid, 302)
		return
	}
	cmt := &models.Comment{
		Contract_id: cid,
		Title:       curUser.Title,
		Uname:       curUser.Uname,
		Cname:       curUser.Cname,
		Date:        time.Now().Format(time.RFC1123),
		Changes:     "删除附件:" + fname,
	}
	err = models.AddComment(cmt)
	this.RedirectTo("/status", "删除附件成功:", "/attachment/view?cid="+cid, 302)
}
func (this *AttachmentController) View() {
	cid := this.GetString("cid", "-1")
	this.Data["Attachments"], _ = models.GetAttachments(cid)
	this.TplName = "attachments.html"
	this.Data["CurUser"] = GetCurAcct(this.Ctx)
	this.Data["cid"] = cid
}
