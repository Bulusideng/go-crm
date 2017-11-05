package controllers

import (
	"io"
	"net/url"
	"os"

	"github.com/astaxie/beego"
)

type FileController struct {
	beego.Controller
}

func (c *FileController) Get() {
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
