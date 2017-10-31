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
	Client_Tel     string
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

type Selector struct {
	CurSelected string
	List        map[string]bool
}

type ContractSelector struct {
	Contract_id   Selector
	Seq           Selector
	Client_name   Selector
	Country       Selector
	Project_type  Selector
	Consulter     Selector
	Secretary     Selector
	Zhuan_an_date Selector
	Current_state Selector
}

func NewContractSelectors() *ContractSelector {
	return &ContractSelector{
		Contract_id:   Selector{"ALL", map[string]bool{}},
		Seq:           Selector{"ALL", map[string]bool{}},
		Client_name:   Selector{"ALL", map[string]bool{}},
		Country:       Selector{"ALL", map[string]bool{}},
		Project_type:  Selector{"ALL", map[string]bool{}},
		Consulter:     Selector{"ALL", map[string]bool{}},
		Secretary:     Selector{"ALL", map[string]bool{}},
		Zhuan_an_date: Selector{"ALL", map[string]bool{}},
		Current_state: Selector{"ALL", map[string]bool{}},
	}
}

func AddContract(c *Contract) error {
	o := orm.NewOrm()
	_, err := o.Insert(c)
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
	return c, err
}

func UpdateContract(c *Contract) error {
	o := orm.NewOrm()
	old := *c
	err := o.Read(&old)
	if err == nil {
		if reflect.DeepEqual(*c, old) { //compae and update reflect.DeepEqual()
			_, err = o.Update(c)
		}
	}
	return err
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
