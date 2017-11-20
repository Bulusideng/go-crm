package models

import (
	"os"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type ALink string

type Attachment struct {
	Id          int
	Author      string
	Time        string
	Contract_id string
	Link        string
	Name        string
	Descrip     string
}

func AddAttachment(a *Attachment) error {

	a.Time = time.Now().Format(time.RFC3339)
	o := orm.NewOrm()
	_, err := o.Insert(a)
	if err != nil {
		beego.Debug("Add Attachment failed: ", *a, "Error: ", err.Error())
		return err
	}
	beego.Debug("Add Attachment success: ", *a)
	return nil
}

func DelAttachment(id int) (string, error) {
	o := orm.NewOrm()

	a := &Attachment{
		Id: id,
	}
	err := o.Read(a)
	if err != nil {
		return "", err
	}
	err = os.Remove(a.Link[1:])
	if err == nil || os.IsNotExist(err) { //Success or file not exit, delete the record in db
		_, err = o.Delete(a)
	}
	return a.Name, err
}

func GetAttachments(cid string) ([]*Attachment, error) {
	o := orm.NewOrm()
	atts := []*Attachment{}
	_, err := o.QueryTable("attachment").Filter("contract_id", cid).All(&atts)
	return atts, err
}
