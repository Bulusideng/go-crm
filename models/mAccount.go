package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func AddAccount(c *Account) error {
	c.CreateDate = time.Now().Format(time.RFC3339)
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
	if err != nil {
		fmt.Printf("Get account failed: [%s]\n", uname)
	}

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
func ManageAccount(nAcct *Account) error {
	o := orm.NewOrm()
	acct := &Account{Uname: nAcct.Uname}
	err := o.Read(acct)
	if err == nil {
		var mgr *Account
		if mgr, err = GetAccount(nAcct.Manager); err == nil {
			if mgr.IsAdmin() || mgr.IsManager() {
				acct.Cname = nAcct.Cname
				acct.Title = nAcct.Title
				acct.Manager = nAcct.Manager
				acct.Status = nAcct.Status
				_, err = o.Update(acct, "Cname", "Title", "Manager", "Status")
				fmt.Printf("manage acct success %+v...........\n", *acct)
			} else {
				err = errors.New(nAcct.Manager + " is not manager")
			}
		} else {
			err = errors.New("Invalid manager: " + nAcct.Manager)
		}
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

func GetAllAccts() ([]*Account, error) {
	o := orm.NewOrm()

	users := make([]*Account, 0)

	qs := o.QueryTable("Account")
	qs.Limit(-1)
	_, err := qs.OrderBy("Title").All(&users)
	if err != nil {
		fmt.Printf("GetAllUsers failed:%s\n", err.Error())
	}
	return users, err
}

func GetAcctounts(filters map[string]string) ([]*Account, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("Account")
	for field, value := range filters {
		//qs = qs.Filter(strings.ToLower(field), value)
		qs = qs.Filter(field, value)
	}
	users := []*Account{}
	_, err := qs.OrderBy("Title").All(&users)
	return users, err
}

func AcctByTitle(title string) ([]*Account, error) {
	filter := map[string]string{
		"Title": title,
	}
	return GetAcctounts(filter)
}

func GetNonAdmins() ([]*Account, error) {
	o := orm.NewOrm()

	users := make([]*Account, 0)

	qs := o.QueryTable("Account")
	qs.Limit(-1)
	_, err := qs.Exclude("Title", "Admin").OrderBy("Title").All(&users) //All account exclude admin
	if err != nil {
		fmt.Printf("GetNonAdmins failed:%s\n", err.Error())
	}
	return users, err
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
