package models

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func AddAccount(c *Account) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		fmt.Printf("Add user %+v failed\n", *c)
		return err
	}
	fmt.Printf("Add user %+v success\n", *c)
	return nil
}

func DelAccount(c *Account) error {
	o := orm.NewOrm()
	_, err := o.Delete(c)
	return err
}

func GetAccount(uname string) (c *Account, err error) {
	o := orm.NewOrm()
	c = &Account{Uname: uname}
	err = o.Read(c)

	return c, err
}

func GetValidAcct(uname, pwd string) (c *Account, err error) {
	c, err = GetAccount(uname)
	if err == nil {
		if c.Pwd != pwd {
			return nil, errors.New("Invalid password")
		} else if !c.IsActive() {
			return nil, errors.New("Account is inactive")
		}
	}
	return c, err
}

//Only one's manager can do this!!!
func ManageAccount(uname, status, manager, title string) error {
	o := orm.NewOrm()
	acct := &Account{Uname: uname}
	err := o.Read(acct)
	if err == nil {
		if mgr, err := GetAccount(manager); err == nil && mgr.IsManagerOf(acct) {
			acct.Status = status
			acct.Manager = manager
			acct.Title = title
			_, err = o.Update(acct, "Status", "Manager", "Title")
		} else {
			err = errors.New("Invalid manager")
		}
	}
	if err != nil {
		fmt.Printf("manage acct failed[%s, %s, %s]...........%s\n", uname, status, manager, err.Error())
	} else {
		fmt.Printf("manage acct success[%s, %s, %s]...........%s\n", uname, status, manager)
	}
	return err
}

func UpdateContact(uname, email, mobile string) error {
	o := orm.NewOrm()
	usr := &Account{Uname: uname}
	err := o.Read(usr)
	if err == nil {
		if usr.IsActive() {
			usr.Email = email
			usr.Mobile = mobile
			_, err = o.Update(usr, "email", "mobile")
		} else {
			return errors.New("Invalid user status")
		}

	}
	return err
}

func UpdatePwd(uname, pwd string) error {
	o := orm.NewOrm()
	usr := &Account{Uname: uname}
	err := o.Read(usr)
	if err == nil {
		if usr.IsActive() {
			usr.Pwd = pwd
			_, err = o.Update(usr, "pwd")
		} else {
			return errors.New("Invalid user status")
		}
	}
	return err
}

func UpdateErrCnt(uname string, delta int) error {
	o := orm.NewOrm()
	usr := &Account{Uname: uname}
	err := o.Read(usr)
	if err == nil {
		if usr.IsActive() {
			usr.ErrCnt += delta
			if usr.ErrCnt < 0 {
				usr.ErrCnt = 0
			} else if usr.ErrCnt >= 3 {
				usr.Lock()
			}
			_, err = o.Update(usr, "err_cnt", "status")
		} else {
			return errors.New("Invalid user status")
		}
	}
	return err
}

func GetAcctounts(filters map[string]string) ([]*Account, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("Account")
	for k, v := range filters {
		qs = qs.Filter(strings.ToLower(k), v)
	}
	users := []*Account{}
	_, err := qs.All(&users)
	return users, err
}

func GetAllAccts() ([]*Account, error) {
	o := orm.NewOrm()

	users := make([]*Account, 0)

	qs := o.QueryTable("Account")
	qs.Limit(-1)
	_, err := qs.All(&users)
	if err != nil {
		fmt.Printf("GetAllUsers failed:%s\n", err.Error())
	}
	return users, err
}

func GetManagers() ([]*Account, error) {
	filter := map[string]string{
		"Title": "Admin",
	}
	admins, _ := GetAcctounts(filter)
	filter["Title"] = "Manager"
	mgrs, _ := GetAcctounts(filter)
	return append(mgrs, admins...), nil

}

func RegisterAdmin() {
	users, err := GetAcctounts(map[string]string{"title": "Admin"})
	for err != nil {
		time.Sleep(time.Second)
		users, err = GetAllAccts()
	}
	if err == nil && len(users) == 0 { ///No user yet, add the user in config to DB
		uname := beego.AppConfig.String("uname")
		pwd := beego.AppConfig.String("pwd")
		usr := &Account{
			Uname:  uname,
			Cname:  "管理员",
			Title:  Titles[EAdmin],
			Email:  "364718765@qq.com",
			Pwd:    pwd,
			Status: "Active",
			ErrCnt: 0,
		}
		AddAccount(usr)
		usr = &Account{
			Uname:  "dengx",
			Cname:  "管理员",
			Title:  Titles[EAdmin],
			Mobile: "15010275683",
			Email:  "364718765@qq.com",
			Pwd:    "dengx",
			Status: "Active",
			ErrCnt: 0,
		}
		AddAccount(usr)
	}

}
