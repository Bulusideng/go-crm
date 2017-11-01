package models

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
)

//合同号 序号	客户姓名 国家 项目 咨询 文案 转案日期 目前状态

type Contract struct {
	Contract_id    string `orm:"pk"`
	Seq            string
	Client_name    string
	Client_tel     string
	Country        string
	Project_type   string
	Consulter      string
	Consulter_name string
	Secretary      string
	Secretary_name string
	Create_date    string
	Create_by      string
	Zhuan_an_date  string
	Current_state  string
	//Histories      []*History
}

type ContractSelector struct {
	Contract_id   FieldSelector
	Seq           FieldSelector
	Client_name   FieldSelector
	Client_tel    FieldSelector
	Country       FieldSelector
	Project_type  FieldSelector
	Consulter     FieldSelector
	Secretary     FieldSelector
	Create_date   FieldSelector
	Create_by     FieldSelector
	Zhuan_an_date FieldSelector
	Current_state FieldSelector
}
type FieldSelector struct {
	CurSelected string          //Filter value of field
	List        map[string]bool //All values of the field
}

func NewContractSelectors() *ContractSelector {
	return &ContractSelector{
		Contract_id:   FieldSelector{"ALL", map[string]bool{}},
		Seq:           FieldSelector{"ALL", map[string]bool{}},
		Client_name:   FieldSelector{"ALL", map[string]bool{}},
		Client_tel:    FieldSelector{"ALL", map[string]bool{}},
		Country:       FieldSelector{"ALL", map[string]bool{}},
		Project_type:  FieldSelector{"ALL", map[string]bool{}},
		Consulter:     FieldSelector{"ALL", map[string]bool{}},
		Secretary:     FieldSelector{"ALL", map[string]bool{}},
		Create_date:   FieldSelector{"ALL", map[string]bool{}},
		Create_by:     FieldSelector{"ALL", map[string]bool{}},
		Zhuan_an_date: FieldSelector{"ALL", map[string]bool{}},
		Current_state: FieldSelector{"ALL", map[string]bool{}},
	}
}

func AddContract(c *Contract) error {
	acct, err := GetAccount(c.Consulter)
	if err == nil {
		c.Consulter_name = acct.Cname
	}
	acct, err = GetAccount(c.Secretary)
	if err == nil {
		c.Secretary_name = acct.Cname
	}
	o := orm.NewOrm()
	_, err = o.Insert(c)
	if err != nil {
		fmt.Printf("Add Contract failed:%s %+v\n", err.Error(), *c)
		return err
	} else {
		fmt.Printf("Add Contract success %+v\n", *c)
	}
	return nil
}

func DelContract(c *Contract) error {
	o := orm.NewOrm()
	_, err := o.Delete(c)
	return err
}

func GetContract(contractId string) (c *Contract, err error) {
	o := orm.NewOrm()
	c = &Contract{
		Contract_id: contractId,
	}
	err = o.Read(c)
	if err != nil {
		fmt.Printf("GetContract[%s] failed %+v\n", contractId, err.Error())
	}
	return c, err
}

func UpdateContract(c *Contract) (*ChangeSlice, error) {
	acct, err := GetAccount(c.Consulter)
	if err == nil {
		c.Consulter_name = acct.Cname
	}
	acct, err = GetAccount(c.Secretary)
	if err == nil {
		c.Secretary_name = acct.Cname
	}
	o := orm.NewOrm()
	old := *c
	err = o.Read(&old)
	changes := ChangeSlice{}
	if err == nil {
		if !reflect.DeepEqual(*c, old) { //compae and update reflect.DeepEqual()
			if c.Client_name != old.Client_name {
				changes = append(changes, Change{Item: "客户姓名", Last: old.Client_name, Current: c.Client_name})
			}
			if c.Client_tel != old.Client_tel {
				changes = append(changes, Change{Item: "客户电话", Last: old.Client_tel, Current: c.Client_tel})
			}
			if c.Country != old.Country {
				changes = append(changes, Change{Item: "国家", Last: old.Country, Current: c.Country})
			}
			if c.Consulter != old.Consulter {
				changes = append(changes, Change{Item: "咨询", Last: old.Consulter_name + "[" + old.Consulter + "]", Current: c.Consulter_name + "[" + c.Consulter + "]"})
			}
			if c.Secretary != old.Secretary {
				changes = append(changes, Change{Item: "文案", Last: old.Secretary_name + "[" + old.Secretary + "]", Current: c.Secretary_name + "[" + c.Secretary + "]"})
			}
			if c.Zhuan_an_date != old.Zhuan_an_date {
				changes = append(changes, Change{Item: "转案日期", Last: old.Zhuan_an_date, Current: c.Zhuan_an_date})
			}
			if c.Current_state != old.Current_state {
				changes = append(changes, Change{Item: "状态", Last: old.Current_state, Current: c.Current_state})
			}
			_, err = o.Update(c)
		}
	}
	if err != nil {
		return nil, err
	}
	return &changes, nil
}

func GetContracts(filters map[string]string) ([]*Contract, error) {
	o := orm.NewOrm()

	qs := o.QueryTable("Contract")
	qs.Limit(-1)
	for k, v := range filters {
		qs = qs.Filter(strings.ToLower(k), v)
	}
	contracts := []*Contract{}
	_, err := qs.OrderBy("contract_id").All(&contracts)
	return contracts, err
}

func GetAllContracts() ([]*Contract, error) {
	o := orm.NewOrm()

	contracts := make([]*Contract, 0)

	qs := o.QueryTable("Contract")
	_, err := qs.Limit(-1).OrderBy("contract_id").All(&contracts)
	return contracts, err
}

func NewContract() *Contract {
	return &Contract{
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		"N/A",
		//nil,
	}
}
