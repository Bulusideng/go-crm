package controllers

import (
	"github.com/astaxie/beego"
	"io"
	"net/url"
	"os"
)

type AttachmentController struct {
	beego.Controller
}

func (c *AttachmentController) Get() {
	filePath, err := url.QueryUnescape(c.Ctx.Request.RequestURI[1:])
	if err != nil {
		beego.Error(err)
		return
	}
	f, err := os.Open(filePath)
	if err != nil {
		c.Ctx.WriteString(err.Error())
	}
	defer f.Close()

	_, err = io.Copy(c.Ctx.ResponseWriter, f)

}
