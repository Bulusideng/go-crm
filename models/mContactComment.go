package models

import (
	"fmt"

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
		str += fmt.Sprintf(" {%s: %s -> %s}; ", c.Item, c.Last, c.Current)
	}
	return str
}

type Comment struct {
	Id          int
	Contract_id int64
	Title       string //The author's title
	Uname       string //The author
	Cname       string //The author
	Date        string
	Changes     string
	Content     string
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
		return err
	}
	return nil
}

func GetComments(contract_id string) ([]*Comment, error) {
	o := orm.NewOrm()
	comments := make([]*Comment, 0)
	qs := o.QueryTable("Comment")
	_, err := qs.Limit(-1).Filter("contract_id", contract_id).All(&comments)
	return comments, err
}
