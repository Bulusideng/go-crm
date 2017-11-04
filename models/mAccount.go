package models

import (
	"errors"
	"fmt"
	"time"

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
		return nil, errors.New("用户名不存在: " + uname)
	}

	return c, nil
}

func GetValidAcct(uname, pwd string) (c *Account, err error) {
	c, err = GetAccount(uname)
	if err == nil {
		if c.Disabled() {
			return nil, errors.New("账户已禁用")
		} else if c.Locked() {
			return nil, errors.New("账户已锁定")
		} else if c.Pwd != pwd {
			return c, errors.New("密码错误")
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
		if mgr, err = GetAccount(nAcct.Manager); err == nil { //Validate the manager field
			if mgr.IsAdmin() || mgr.IsManager() {
				acct.Cname = nAcct.Cname
				acct.Title = nAcct.Title
				acct.Manager = nAcct.Manager
				fmt.Printf("manage acct status:%s\n", acct.Status)
				switch nAcct.Status {
				case "Active":
					acct.Enable()
				case "Disabled":
					acct.Disable()
				case "Locked":
					acct.Lock()
				}
				acct.Status = nAcct.Status
				_, err = o.Update(acct, "Cname", "Title", "Manager", "Status", "ErrCnt")
				fmt.Printf("manage acct success %+v...........\n", *acct)
			} else {
				fmt.Printf("manage acct error %+v...........\n", err)
				err = errors.New(nAcct.Manager + " is not manager")
			}
		} else {
			fmt.Printf("manage acct error %+v...........\n", err)
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

var maxErrCnt = 3

func UpdateErrCnt(uname string, delta int) (int, error) {
	remaining := 3
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
				remaining = 0
			} else {
				remaining = 3 - usr.ErrCnt
			}
			fmt.Printf("UpdateErrCnt %s delta:%d, newcnt:%d, locked:%t\n", uname, delta, usr.ErrCnt, usr.Locked())
			_, err = o.Update(usr, "err_cnt", "status")
		} else {
			return 0, errors.New("Invalid user status")
		}
	}
	return remaining, err
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
