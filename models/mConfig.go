package models

import (
	"errors"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Config struct {
	Id                    int
	Contract_view_control bool
	Pwdreset              bool //Support pwd support, if yes, we must have email acct setup
	Mailaddr              string
	Smtpaddr              string
	Mailpwd               string
	inited                bool
}

var (
	ID     = 9876543
	config = &Config{
		Id: ID,
		Contract_view_control: true,
		Mailaddr:              "dbiti@163.com",
		Smtpaddr:              "smtp.163.com",
		Mailpwd:               "oxiwangyi.",
	}
)

func GetConfig() *Config {
	if !config.inited {
		c, e := readConfig()
		if e == nil {
			config = c
			config.inited = true
		}
	}
	return config
}

func readConfig() (c *Config, err error) {
	c = &Config{
		Id: ID,
	}
	o := orm.NewOrm()
	err = o.Read(c)
	if err != nil {
		_, err = o.Insert(c)
	}
	beego.Warn("Read config:", c)
	return c, err
}

func UpdateConfig(c *Config) error {
	c.Id = ID

	idx := strings.Index(c.Mailaddr, "@")
	if idx <= 0 {
		if c.Pwdreset {
			return errors.New("非法邮件地址")
		} else {
			c.Mailaddr = ""
			c.Smtpaddr = ""
		}
	} else {
		c.Smtpaddr = "smtp." + c.Mailaddr[idx+1:]
	}

	config = c
	o := orm.NewOrm()

	old := &Config{
		Id: ID,
	}

	old, err := readConfig()
	if err != nil {
		_, err = o.Insert(c)
	} else {
		if !reflect.DeepEqual(c, old) {
			_, err = o.Update(c)
			if err != nil {
				beego.Error("Update config failed: ", *c, "Error: ", err.Error())
				return err
			}
		} else {
			return errors.New("配置没有改变")
		}
	}
	if err == nil {
		beego.Warn("Update config success:", c)
	} else {
		beego.Warn("Update config failed: ", err.Error())
	}

	return err
}
