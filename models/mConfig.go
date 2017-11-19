package models

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type Config struct {
	Id                    int
	Contract_view_control bool
}

func GetConfig(c *Config) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		beego.Debug("Add user failed: ", *c, "Error: ", err.Error())
		return err
	}
	beego.Debug("Add user success: ", *c)
	return nil
}

func UpdateConfig(c *Config) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
	if err != nil {
		beego.Debug("Add user failed: ", *c, "Error: ", err.Error())
		return err
	}
	beego.Debug("Add user success: ", *c)
	return nil
}
