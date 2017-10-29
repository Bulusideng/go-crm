package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

const (
	EAdmin = iota //0
	EManager
	EConsulter
	ESecretary
	EGuest
)

var Titles = []string{
	"Admin",
	"Manager",
	"Consulter",
	"Secretary",
	"Guest",
}

func NewGuest() *UserInfo {
	return &UserInfo{
		Uname: "Guest",
		Title: "Guest",
		Cname: "访客",
	}
}

type UserInfo struct {
	Uname  string `orm:"pk"`
	Title  string //Administrator, Manager, Consulter or Secretary
	Cname  string //Chinese name
	Email  string
	Mobile string
}

type User struct {
	UserInfo
	Pwd string
}

func AddUser(c *User) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		fmt.Printf("Add user %+v failed\n", *c)
		return err
	}
	fmt.Printf("Add user %+v success\n", *c)
	return nil
}

func DelUser(c *User) error {
	o := orm.NewOrm()
	_, err := o.Delete(c)
	return err
}

func GetUser(uname string) (c *User, err error) {
	o := orm.NewOrm()
	c = &User{
		UserInfo: UserInfo{Uname: uname},
	}
	err = o.Read(c)

	return c, err
}

func GetValidUser(uname, pwd string) (c *User, err error) {
	c, err = GetUser(uname)
	if err == nil {
		if c.Pwd != pwd {
			return nil, errors.New("Invalid password")
		}
	}
	return c, err
}

func UpdateUser(c *User) error {
	o := orm.NewOrm()
	old := *c
	err := o.Read(&old)
	if err == nil {
		if *c != old { //compae and update reflect.DeepEqual()
			_, err = o.Update(c)
		}
	}
	return err
}

func GetUsers(filters map[string]string) ([]*User, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("User")
	for k, v := range filters {
		qs = qs.Filter(strings.ToLower(k), v)
	}
	users := []*User{}
	_, err := qs.All(&users)
	return users, err
}

func GetAllUsers() ([]*User, error) {
	o := orm.NewOrm()

	users := make([]*User, 0)

	qs := o.QueryTable("User")
	qs.Limit(-1)
	_, err := qs.All(&users)
	if err != nil {
		fmt.Printf("GetAllUsers failed:%s\n", err.Error())
	}
	return users, err
}

func RegisterAdmin() {
	users, err := GetUsers(map[string]string{"title": "Admin"})
	for err != nil {
		time.Sleep(time.Second)
		users, err = GetAllUsers()
	}

	if err == nil && len(users) == 0 { ///No user yet, add the user in config to DB
		uname := beego.AppConfig.String("uname")
		pwd := beego.AppConfig.String("pwd")
		usr := &User{
			UserInfo: UserInfo{
				Uname:  uname,
				Title:  Titles[EAdmin],
				Mobile: "15010275683",
				Email:  "673764189@qq.com",
			},
			Pwd: pwd,
		}
		AddUser(usr)
	}

}
