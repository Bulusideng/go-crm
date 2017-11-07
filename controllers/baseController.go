package controllers

import (
	"log"
	"net/url"

	"github.com/astaxie/beego"
)

type baseController struct {
	beego.Controller
}

func (this *baseController) RedirectTo(to, msg, redirect string, code int) {
	u, err := url.Parse(to)
	if err != nil {
		log.Fatal(err)
	}

	q := u.Query()

	q.Set("Message", msg)
	q.Set("RedirectTo", redirect)
	u.RawQuery = q.Encode()
	str := to + "?" + u.RawQuery
	log.Printf("Redirect to:[%s]\n", str)
	this.Redirect(str, code)
	return
}
