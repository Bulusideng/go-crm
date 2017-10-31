package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Comment struct {
	Id          int
	Contract_id string
	User        string
	Date        string
	Text        string
}

func AddComment(c *Comment) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		fmt.Printf("Add Comment failed:%s %+v\n", err.Error(), *c)
		return err
	} else {
		fmt.Printf("Add Comment success %+v\n", *c)
	}
	return nil
}

func DelComment(id int) error {
	o := orm.NewOrm()
	if _, err := o.Delete(&Comment{Id: id}); err != nil {
		fmt.Printf("Delete Comment failed %d\n", id)
	}
	return err
}

func GetComments(contract_id string) ([]*Comment, error) {
	o := orm.NewOrm()
	comments := make([]*Comment, 0)
	qs := o.QueryTable("Comment")
	_, err := qs.Limit(-1).All(&comments)
	return comments, err
}
