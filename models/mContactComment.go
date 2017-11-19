package models

import (
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Change struct {
	Id      int
	Item    string
	Last    string
	Current string
}

type ChangeSlice []Change

func (this *ChangeSlice) String() string {
	str := ""
	for _, c := range *this {
		str += fmt.Sprintf(`{%s: %s -> %s}; `, c.Item, c.Last, c.Current)
	}
	return str
}

type Comment struct {
	Id          int
	Contract_id string
	Title       string //The author's title
	Uname       string //The author
	Cname       string //The author
	Date        string
	Changes     string
	Content     string
	Attach      string
}

func AddComment(c *Comment) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		beego.Error("Add Comment failed: ", err.Error(), *c)
		return err
	} else {
		beego.Debug("Add Comment success: ", *c)
	}
	return nil
}

func UpdateComment(c *Comment) error {
	o := orm.NewOrm()
	_, err := o.Update(c)
	if err != nil {
		beego.Error("Update Comment failed: ", err.Error(), *c)
		return err
	} else {
		beego.Debug("Update Comment success: ", *c)
	}
	return nil
}

func DelComment(id int) error {
	o := orm.NewOrm()
	if _, err := o.Delete(&Comment{Id: id}); err != nil {
		beego.Error("Delete Comment failed: ", id)
		return err
	}
	return nil
}

func GetComments(contract_id string) ([]*Comment, error) {
	o := orm.NewOrm()
	comments := make([]*Comment, 0)
	qs := o.QueryTable("Comment")
	if contract_id != "" {
		qs = qs.Limit(-1).Filter("contract_id", contract_id).OrderBy("-id")
	}
	_, err := qs.All(&comments)
	return comments, err
}
